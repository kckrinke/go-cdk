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
	"unicode"
)

type WordLine struct {
	words []*WordCell
}

func NewEmptyWordLine() *WordLine {
	return &WordLine{
		words: make([]*WordCell, 0),
	}
}

func NewWordLine(line string, style Style) *WordLine {
	wl := &WordLine{}
	wl.SetLine(line, style)
	return wl
}

func (w *WordLine) SetLine(line string, style Style) {
	w.words = make([]*WordCell, 0)
	isWord := false
	wid := 0
	for _, c := range line {
		if unicode.IsSpace(c) {
			if isWord || len(w.words) == 0 {
				isWord = false
				w.words = append(w.words, NewEmptyWordCell())
				wid = len(w.words) - 1
			}
			// appending to the "space" word
			w.words[wid].AppendRune(c, style)
		} else if c == '\n' {
			// always a reset and single-char-word of it's own
			isWord = false
			w.words = append(w.words, NewWordCell("\n", style))
			wid = len(w.words) - 1
		} else {
			if !isWord || len(w.words) == 0 {
				isWord = true
				w.words = append(w.words, NewEmptyWordCell())
				wid = len(w.words) - 1
			}
			// appending to the "real" word
			w.words[wid].AppendRune(c, style)
		}
	}
}

func (w *WordLine) AppendWord(word string, style Style) {
	w.words = append(w.words, NewWordCell(word, style))
}

func (w *WordLine) AppendWordCell(word *WordCell) {
	w.words = append(w.words, word)
}

func (w WordLine) GetWord(index int) *WordCell {
	if index < len(w.words) {
		return w.words[index]
	}
	return nil
}

func (w *WordLine) RemoveWord(index int) {
	if index < len(w.words) {
		w.words = append(
			w.words[:index],
			w.words[index+1:]...,
		)
	}
}

func (w WordLine) GetCharacter(index int) *CTextCell {
	if index < w.CharacterCount() {
		count := 0
		for _, word := range w.words {
			for _, c := range word.characters {
				if count == index {
					return c
				}
				count++
			}
		}
	}
	return nil
}

func (w WordLine) Words() []*WordCell {
	return w.words
}

func (w WordLine) Len() (wordSpaceCount int) {
	return len(w.words)
}

func (w WordLine) CharacterCount() (count int) {
	for _, word := range w.words {
		count += word.Len()
	}
	return
}

func (w WordLine) WordCount() (wordCount int) {
	for _, word := range w.words {
		if !word.IsSpace() {
			wordCount++
		}
	}
	return
}

func (w WordLine) GapWordList() (gaps []*WordCell) {
	gaps = make([]*WordCell, len(w.words))
	for i, word := range w.words {
		if word.IsSpace() {
			gaps[i] = word
		} else {
			gaps[i] = nil
		}
	}
	return
}

func (w WordLine) GapCount() (count int) {
	for _, word := range w.words {
		if word.IsSpace() {
			count++
		}
	}
	return
}

func (w WordLine) HasSpace() bool {
	for _, word := range w.words {
		if word.IsSpace() {
			return true
		}
	}
	return false
}

func (w WordLine) Value() (s string) {
	s = ""
	for i, c := range w.words {
		if i > 0 {
			s += " "
		}
		s += c.Value()
	}
	return
}

func (w WordLine) String() (s string) {
	s = ""
	for i, c := range w.words {
		if i > 0 {
			s += " "
		}
		s += c.String()
	}
	return
}

// wrap, justify and align the set input,
func (w WordLine) Make(wrap WrapMode, justify Justification, maxChars int) (lines []*WordLine) {
	lines = append(lines, NewEmptyWordLine())
	cid, wid, lid := 0, 0, 0
	for _, word := range w.words {
		if wid >= lines[lid].Len() {
			lines[lid].AppendWordCell(NewEmptyWordCell())
		}
		for _, c := range word.characters {
			switch c.Value() {
			case '\n':
				lines = append(lines, NewEmptyWordLine())
				lid = len(lines) - 1
				wid = 0
			default:
				if wid >= lines[lid].Len() {
					lines[lid].AppendWordCell(NewEmptyWordCell())
				}
				lines[lid].words[wid].AppendRune(c.Value(), c.Style())
			}
			cid++
		}
		wid++
	}
	lines = w.applyTypography(wrap, justify, maxChars, lines)
	return
}

func (w WordLine) applyTypography(wrap WrapMode, justify Justification, maxChars int, in []*WordLine) (out []*WordLine) {
	// all space-words must be applied as 1 width
	switch wrap {
	case WRAP_WORD:
		// break onto inserted/new line at end gap
		// - if line has no breakpoints, truncate
		out = w.applyTypographicWrapWord(maxChars, in)
	case WRAP_WORD_CHAR:
		// break onto inserted/new line at end gap
		// - if line has no breakpoints, fallthrough
		out = w.applyTypographicWrapWordChar(maxChars, in)
	case WRAP_CHAR:
		// break onto inserted/new line at maxChars
		out = w.applyTypographicWrapChar(maxChars, in)
	case WRAP_NONE:
		// truncate each line to maxChars
		out = w.applyTypographicWrapNone(maxChars, in)
	}
	switch justify {
	case JUSTIFY_FILL:
		// each non-empty line is space-expanded to fill maxChars
	case JUSTIFY_CENTER:
		// each non-empty line is centered on halfway maxChars
	case JUSTIFY_RIGHT:
		// each non-empty line is left-padded to fill maxChars
	case JUSTIFY_LEFT:
		fallthrough
	default:
		// each non-empty line has leading space removed
	}
	return
}

func (w WordLine) applyTypographicWrapWord(maxChars int, input []*WordLine) (out []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(out) {
			out = append(out, NewEmptyWordLine())
		}
		if line.Len() > maxChars {
			if !line.HasSpace() {
				// nothing to break on, truncate on maxChars
			truncateWrapWord:
				for _, word := range line.Words() {
					if wid >= out[lid].Len() {
						out[lid].AppendWordCell(NewEmptyWordCell())
					}
					for _, c := range word.Characters() {
						if cid > maxChars {
							lid = len(out) // don't append trailing NEWLs
							break truncateWrapWord
						}
						out[lid].words[wid].AppendRune(c.Value(), c.Style())
						cid++
					}
					wid++
				}
				continue
			}
		}
		for _, word := range line.Words() {
			if wid >= out[lid].Len() {
				out[lid].AppendWordCell(NewEmptyWordCell())
			}
			wordLen := word.Len()
			if word.IsSpace() {
				wordLen = 1
			}
			if cid+wordLen >= maxChars {
				out = append(out, NewEmptyWordLine())
				lid = len(out) - 1
				wid = 0
				if !word.IsSpace() {
					out[lid].AppendWordCell(word)
					wid = out[lid].Len() - 1
				}
				continue
			}
			if word.IsSpace() {
				c := word.GetCharacter(0)
				wc := NewEmptyWordCell()
				wc.AppendRune(c.Value(), c.Style())
				out[lid].AppendWordCell(wc)
				cid += wc.Len()
			} else {
				out[lid].AppendWordCell(word)
				cid += word.Len()
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}

func (w WordLine) applyTypographicWrapWordChar(maxChars int, input []*WordLine) (out []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(out) {
			out = append(out, NewEmptyWordLine())
		}
		if line.Len() > maxChars {
			if !line.HasSpace() {
				// nothing to break on, wrap on maxChars
				for _, word := range line.Words() {
					if wid >= out[lid].Len() {
						out[lid].AppendWordCell(NewEmptyWordCell())
					}
					for _, c := range word.Characters() {
						if cid > maxChars {
							out = append(out, NewEmptyWordLine())
							lid = len(out) - 1
							out[lid].AppendWordCell(NewEmptyWordCell())
							wid = 0
							cid = 0
						}
						out[lid].words[wid].AppendRune(c.Value(), c.Style())
						cid++
					}
					wid++
				}
				lid++
				wid = 0
				cid = 0
				continue
			}
		}
		for _, word := range line.Words() {
			if wid >= out[lid].Len() {
				out[lid].AppendWordCell(NewEmptyWordCell())
			}
			wordLen := word.Len()
			if word.IsSpace() {
				wordLen = 1
			}
			if cid+wordLen >= maxChars {
				out = append(out, NewEmptyWordLine())
				lid = len(out) - 1
				wid = 0
				if !word.IsSpace() {
					out[lid].AppendWordCell(word)
					wid = out[lid].Len() - 1
				}
				continue
			}
			if word.IsSpace() {
				c := word.GetCharacter(0)
				wc := NewEmptyWordCell()
				wc.AppendRune(c.Value(), c.Style())
				out[lid].AppendWordCell(wc)
				cid += wc.Len()
			} else {
				out[lid].AppendWordCell(word)
				cid += word.Len()
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}

func (w WordLine) applyTypographicWrapChar(maxChars int, input []*WordLine) (out []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(out) {
			out = append(out, NewEmptyWordLine())
		}
		for _, word := range line.Words() {
			if word.IsSpace() {
				if cid >= maxChars {
					out = append(out, NewEmptyWordLine())
					lid = len(out) - 1
					wid = 0
					cid = 0
				}
				if c := word.GetCharacter(0); c != nil {
					wc := NewEmptyWordCell()
					wc.AppendRune(c.Value(), c.Style())
					out[lid].AppendWordCell(wc)
					cid += wc.Len()
				}
			} else {
				if cid+word.Len() > maxChars {
					firstHalf, secondHalf := NewEmptyWordCell(), NewEmptyWordCell()
					for _, c := range word.Characters() {
						if cid < maxChars {
							firstHalf.AppendRune(c.Value(), c.Style())
						} else {
							secondHalf.AppendRune(c.Value(), c.Style())
						}
						cid++
					}
					out[lid].AppendWordCell(firstHalf)
					out = append(out, NewEmptyWordLine())
					lid = len(out) - 1
					out[lid].AppendWordCell(secondHalf)
					cid = secondHalf.Len()
					wid = 0
				} else {
					out[lid].AppendWordCell(word)
					cid += word.Len()
				}
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}

func (w WordLine) applyTypographicWrapNone(maxChars int, input []*WordLine) (out []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(out) {
			out = append(out, NewEmptyWordLine())
		}
	forEachWordWrapNone:
		for _, word := range line.Words() {
			if wid >= out[lid].Len() {
				out[lid].AppendWordCell(NewEmptyWordCell())
			}
			if word.IsSpace() {
				if cid > maxChars {
					lid = len(out)
					wid = 0
					cid = 0
					break
				}
				if c := word.GetCharacter(0); c != nil {
					wc := NewEmptyWordCell()
					wc.AppendRune(c.Value(), c.Style())
					out[lid].AppendWordCell(wc)
					cid += wc.Len()
				}
			} else {
				if cid+word.Len() > maxChars {
					for _, c := range word.Characters() {
						if cid > maxChars {
							lid++
							wid = 0
							cid = 0
							break forEachWordWrapNone
						}
						out[lid].words[wid].AppendRune(c.Value(), c.Style())
						cid++
					}
				} else {
					out[lid].AppendWordCell(word)
					cid += word.Len()
				}
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}
