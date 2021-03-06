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
	"sync"

	"github.com/kckrinke/go-cdk/utils"
)

var (
	TapSpace = "    "
)

type TextBuffer interface {
	Set(input string, style Style)
	SetInput(input WordLine)
	Style() Style
	SetStyle(style Style)
	CharacterCount() (cellCount int)
	WordCount() (wordCount int)
	PlainText(wordWrap WrapMode, ellipsize bool, justify Justification, maxChars int) (plain string)
	PlainTextInfo(wordWrap WrapMode, ellipsize bool, justify Justification, maxChars int) (longestLine, lineCount int)
	Draw(canvas Canvas, singleLine bool, wordWrap WrapMode, ellipsize bool, justify Justification, vAlign VerticalAlignment) EventFlag
}

type CTextBuffer struct {
	raw   string
	input WordLine
	style Style

	sync.Mutex
}

func NewEmptyTextBuffer(style Style) TextBuffer {
	return &CTextBuffer{
		style: style,
	}
}

func NewTextBuffer(input string, style Style) TextBuffer {
	tb := &CTextBuffer{
		style: style,
	}
	tb.Set(input, style)
	return tb
}

func (b *CTextBuffer) Set(input string, style Style) {
	b.raw = input
	b.input = NewWordLine(input, style)
}

func (b *CTextBuffer) SetInput(input WordLine) {
	b.input = input
	b.raw = input.Value()
}

func (b *CTextBuffer) Style() Style {
	return b.style
}

func (b *CTextBuffer) SetStyle(style Style) {
	if b.style.String() != style.String() {
		b.style = style
		if b.input != nil {
			b.input = NewWordLine(b.raw, style)
		}
	}
}

func (b *CTextBuffer) CharacterCount() (cellCount int) {
	if b.input != nil {
		cellCount = b.input.CharacterCount()
	}
	return
}

func (b *CTextBuffer) WordCount() (wordCount int) {
	if b.input != nil {
		wordCount = b.input.WordCount()
	}
	return
}

func (b *CTextBuffer) PlainText(wordWrap WrapMode, ellipsize bool, justify Justification, maxChars int) (plain string) {
	lines := b.input.Make(wordWrap, ellipsize, justify, maxChars, b.style)
	for _, line := range lines {
		if len(plain) > 0 {
			plain += "\n"
		}
		for _, word := range line.Words() {
			for _, char := range word.Characters() {
				plain += string(char.Value())
			}
		}
	}
	return
}

func (b *CTextBuffer) PlainTextInfo(wordWrap WrapMode, ellipsize bool, justify Justification, maxChars int) (longestLine, lineCount int) {
	lines := b.input.Make(wordWrap, ellipsize, justify, maxChars, b.style)
	lineCount = len(lines)
	for _, line := range lines {
		lcc := line.CharacterCount()
		if longestLine < lcc {
			longestLine = lcc
		}
	}
	return
}

func (b *CTextBuffer) Draw(canvas Canvas, singleLine bool, wordWrap WrapMode, ellipsize bool, justify Justification, vAlign VerticalAlignment) EventFlag {
	b.Lock()
	defer b.Unlock()
	if b.input == nil || b.input.CharacterCount() == 0 {
		// non-operation
		return EVENT_PASS
	}

	if singleLine {
		wordWrap = WRAP_NONE
	}

	maxChars := canvas.Width()
	lines := b.input.Make(wordWrap, ellipsize, justify, maxChars, b.style)
	size := canvas.GetSize()
	if size.W <= 0 || size.H <= 0 {
		return EVENT_PASS
	}
	if len(lines) == 0 {
		return EVENT_PASS
	}
	lenLines := len(lines)
	if singleLine && lenLines > 1 {
		lenLines = 1
	}

	var atCanvasLine, fromInputLine = 0, 0
	switch vAlign {
	case ALIGN_BOTTOM:
		numLines := lenLines
		if numLines > size.H {
			delta := utils.FloorI(numLines-size.H, 0)
			fromInputLine = delta
		} else {
			delta := size.H - numLines
			atCanvasLine = delta
		}
	case ALIGN_MIDDLE:
		numLines := lenLines
		halfLines := numLines / 2
		halfCanvas := size.H / 2
		delta := utils.FloorI(halfCanvas-halfLines, 0)
		if numLines > size.H {
			fromInputLine = delta
		} else {
			atCanvasLine = delta
		}
	case ALIGN_TOP:
	default:
	}

	y := atCanvasLine
	for lid := fromInputLine; lid < lenLines; lid++ {
		if lid >= len(lines) {
			break
		}
		if y >= size.H {
			break
		}
		x := 0
		for _, word := range lines[lid].Words() {
			for _, c := range word.Characters() {
				if x <= size.W {
					_ = canvas.SetRune(x, y, c.Value(), c.Style())
					x++
				}
			}
		}
		y++
		if singleLine {
			break
		}
	}

	return EVENT_STOP
}
