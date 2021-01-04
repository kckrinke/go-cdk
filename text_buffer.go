// Copyright 2020 The CDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cdk

import (
	"strings"
	"sync"
)

var (
	TapSpace = "    "
)

type TextBuffer interface {
	Set(input string, style Style)
	Len() (count int)
	Draw(canvas *Canvas, singleLineMode bool, wordWrap WrapMode, justify Justification, align VerticalAlignment) EventFlag
}

type CTextBuffer struct {
	lines []*WordLine
	style Style

	sync.Mutex
}

func NewEmptyTextBuffer(style Style) *CTextBuffer {
	return &CTextBuffer{
		lines: make([]*WordLine, 0),
		style: style,
	}
}

func NewTextBuffer(input string, style Style) *CTextBuffer {
	tb := &CTextBuffer{
		lines: make([]*WordLine, 0),
		style: style,
	}
	tb.Set(input, style)
	return tb
}

func (b *CTextBuffer) Set(input string, style Style) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		b.lines = append(b.lines, NewWordLine(line, style))
	}
}

func (b *CTextBuffer) Len() (count int) {
	for _, line := range b.lines {
		count += line.CharacterCount()
	}
	return
}

func (b *CTextBuffer) LenWords() (wordCount int) {
	for _, line := range b.lines {
		wordCount += line.WordCount()
	}
	return
}

func (b *CTextBuffer) Draw(canvas *Canvas, singleLine bool, wordWrap WrapMode, justify Justification, align VerticalAlignment) EventFlag {
	if len(b.lines) == 0 {
		// non-operation
		return EVENT_PASS
	}
	if !singleLine && canvas.size.H == 1 {
		singleLine = true
	}
	maxChars := canvas.size.W
	if singleLine {
		var atCanvasLine int
		switch align {
		case ALIGN_MIDDLE:
			atCanvasLine = (canvas.size.H / 2) - (len(b.lines) / 2)
		case ALIGN_BOTTOM:
			atCanvasLine = canvas.size.H - len(b.lines)
		case ALIGN_TOP:
		default:
			atCanvasLine = 0
		}
		spaces := b.Len()
		if spaces > maxChars {
			b.truncateSingleLine(atCanvasLine, 0, canvas, wordWrap, maxChars)
			return EVENT_STOP
		}
		b.writeJustifiedLines(atCanvasLine, []*WordLine{b.lines[0]}, canvas, wordWrap, maxChars, justify)
		return EVENT_STOP
	}
	b.wrapMultiLine(canvas, wordWrap, justify, align)
	return EVENT_STOP
}

func (b *CTextBuffer) wrapMultiLine(canvas *Canvas, wordWrap WrapMode, justify Justification, align VerticalAlignment) {
	var atLine int
	switch align {
	case ALIGN_MIDDLE:
		atLine = (canvas.size.H / 2) - (len(b.lines) / 2)
	case ALIGN_BOTTOM:
		atLine = canvas.size.H - len(b.lines)
	case ALIGN_TOP:
	default:
		atLine = 0
	}
	maxChars := canvas.size.W
	sorted := b.wrapSort(maxChars, wordWrap)
	b.writeJustifiedLines(atLine, sorted, canvas, wordWrap, maxChars, justify)
}

func (b *CTextBuffer) wrapSortNop(maxChars, atLine int, sorted []*WordLine, line *WordLine) (out []*WordLine, lid int) {
	out, lid = sorted, atLine
	if lid >= len(out) {
		out = append(out, line)
	} else {
		out[lid] = line
	}
	lid++
	return
}
func (b *CTextBuffer) wrapSortWord(maxChars, atLine int, sorted []*WordLine, line *WordLine) (out []*WordLine, lid int) {
	out, lid = sorted, atLine
	if !line.HasSpace() {
		// truncate at maxChars, or do nothing?
	} else if len(out) < lid {
		for _, word := range line.words {
			if out[lid].CharacterCount()+word.Len() < maxChars {
				out[lid].words = append(out[lid].words, word)
			} else {
				lid++
				out = append(out, NewWordLine("", b.style))
				out[lid].words = append(out[lid].words, word)
			}
		}
	}
	return
}
func (b *CTextBuffer) wrapSortChar(maxChars, atLine int, sorted []*WordLine, line *WordLine) (out []*WordLine, lid int) {
	out, lid = sorted, atLine
	count := 0
	for _, word := range line.words {
		if count+word.Len() < maxChars {
			out[lid].words = append(out[lid].words, word)
		} else {
			delta := maxChars - count
			firstHalf := NewEmptyWordCell()
			secondHalf := NewEmptyWordCell()
			for i := 0; i < word.Len(); i++ {
				c := word.GetCharacter(i)
				if i <= delta {
					firstHalf.AppendRune(c.Value(), c.Style())
				} else {
					secondHalf.AppendRune(c.Value(), c.Style())
				}
			}
			out[lid].words = append(out[lid].words, firstHalf)
			out = append(out, NewEmptyWordLine())
			lid = len(out) - 1
			out[lid].words = append(out[lid].words, secondHalf)
		}
		count += word.Len()
	}
	return
}
func (b *CTextBuffer) wrapSortTruncate(maxChars, atLine int, sorted []*WordLine, line *WordLine) (out []*WordLine, lid int) {
	out, lid = sorted, atLine
	count := 0
	for _, word := range line.words {
		if count+word.Len() < maxChars {
			out[lid].words = append(out[lid].words, word)
		} else {
			delta := maxChars - count
			partial := NewEmptyWordCell()
			for i := 0; i < delta; i++ {
				c := word.GetCharacter(i)
				partial.AppendRune(c.Value(), c.Style())
			}
			out[lid].words = append(out[lid].words, partial)
			break
		}
		count += word.Len()
	}
	return
}
func (b *CTextBuffer) wrapSort(maxChars int, wordWrap WrapMode) []*WordLine {
	var sorted []*WordLine
	sorted = append(sorted, NewEmptyWordLine())
	lid := 0
	for _, line := range b.lines {
		if line.CharacterCount() < maxChars {
			// no need to wrap at all
			sorted, lid = b.wrapSortNop(maxChars, lid, sorted, line)
			continue
		}
		// must wrap or truncate
		switch wordWrap {
		case WRAP_WORD:
			sorted, lid = b.wrapSortWord(maxChars, lid, sorted, line)
		case WRAP_WORD_CHAR:
			// attempt wrap word
			if line.HasSpace() {
				sorted, lid = b.wrapSortWord(maxChars, lid, sorted, line)
				break
			}
			// no breaks in line, fallback to wrap_char
			fallthrough
		case WRAP_CHAR:
			sorted, lid = b.wrapSortChar(maxChars, lid, sorted, line)
		case WRAP_NONE:
			fallthrough
		default:
			// truncate the line
			sorted, lid = b.wrapSortTruncate(maxChars, lid, sorted, line)
		}
	}
	return sorted
}

func (b *CTextBuffer) writeJustifiedLines(fromCanvasLine int, lines []*WordLine, canvas *Canvas, wordWrap WrapMode, maxChars int, justify Justification) {
	cSize := canvas.GetSize()
	for y, line := range lines {
		if fromCanvasLine+y >= cSize.H {
			break
		}
		var justified *WordLine
		switch justify {
		case JUSTIFY_RIGHT:
			justified = line.MakeRight(maxChars, b.style)
		case JUSTIFY_CENTER:
			justified = line.MakeCenter(maxChars, b.style)
		case JUSTIFY_FILL:
			justified = line.MakeFill(maxChars)
		case JUSTIFY_LEFT:
			fallthrough
		default:
			justified = line.MakeLeft(maxChars)
		}
		for x := 0; x < justified.CharacterCount(); x++ {
			c := justified.GetCharacter(x)
			canvas.SetRune(x, fromCanvasLine+y, c.Value(), c.Style())
		}
	}
}


func (b *CTextBuffer) truncateSingleLine(atCanvasLine, forWordLine int, canvas *Canvas, wordWrap WrapMode, maxChars int) {
	if len(b.lines) < forWordLine || len(b.lines[forWordLine].words) == 0 {
		// non-operation
		return
	}
	var line *WordLine
	switch wordWrap {
	case WRAP_WORD:
		line, _ = b.lines[forWordLine].MakeTruncated(maxChars, true)
	case WRAP_WORD_CHAR:
		var didTruncate bool
		if line, didTruncate = b.lines[forWordLine].MakeTruncated(maxChars, true); didTruncate {
			break
		}
		fallthrough
	case WRAP_CHAR:
		fallthrough
	case WRAP_NONE:
		fallthrough
	default:
		line, _ = b.lines[forWordLine].MakeTruncated(maxChars, false)
	}
	if line != nil {
		x := 0
		for _, word := range line.words {
			for _, char := range word.characters {
				canvas.SetRune(x, atCanvasLine, char.Value(), char.Style())
				x += char.Width()
			}
		}
	}
	return
}
