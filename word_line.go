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
	"fmt"
	"unicode"
)

type WordLine interface {
	SetLine(line string, style Style)
	AppendWord(word string, style Style)
	AppendWordCell(word WordCell)
	AppendWordRune(wordIndex int, char rune, style Style) error
	GetWord(index int) WordCell
	RemoveWord(index int)
	GetCharacter(index int) TextCell
	Words() []WordCell
	Len() (wordSpaceCount int)
	CharacterCount() (count int)
	WordCount() (wordCount int)
	HasSpace() bool
	Value() (s string)
	String() (s string)
	Make(wrap WrapMode, justify Justification, maxChars int, fillerStyle Style) (lines []WordLine)
}

type CWordLine struct {
	words []WordCell
}

func NewEmptyWordLine() WordLine {
	return &CWordLine{
		words: make([]WordCell, 0),
	}
}

func NewWordLine(line string, style Style) WordLine {
	wl := &CWordLine{}
	wl.SetLine(line, style)
	return wl
}

func (w *CWordLine) SetLine(line string, style Style) {
	w.words = make([]WordCell, 0)
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

func (w *CWordLine) AppendWord(word string, style Style) {
	w.words = append(w.words, NewWordCell(word, style))
}

func (w *CWordLine) AppendWordCell(word WordCell) {
	w.words = append(w.words, word)
}

func (w *CWordLine) AppendWordRune(wordIndex int, char rune, style Style) error {
	if wordIndex < len(w.words) {
		w.words[wordIndex].AppendRune(char, style)
		return nil
	}
	return fmt.Errorf("word at index %d not found", wordIndex)
}

func (w *CWordLine) GetWord(index int) WordCell {
	if index < len(w.words) {
		return w.words[index]
	}
	return nil
}

func (w *CWordLine) RemoveWord(index int) {
	if index < len(w.words) {
		w.words = append(
			w.words[:index],
			w.words[index+1:]...,
		)
	}
}

func (w *CWordLine) GetCharacter(index int) TextCell {
	if index < w.CharacterCount() {
		count := 0
		for _, word := range w.words {
			for _, c := range word.Characters() {
				if count == index {
					return c
				}
				count++
			}
		}
	}
	return nil
}

func (w *CWordLine) Words() []WordCell {
	return w.words
}

func (w *CWordLine) Len() (wordSpaceCount int) {
	return len(w.words)
}

func (w *CWordLine) CharacterCount() (count int) {
	for _, word := range w.words {
		count += word.Len()
	}
	return
}

func (w *CWordLine) WordCount() (wordCount int) {
	for _, word := range w.words {
		if !word.IsSpace() {
			wordCount++
		}
	}
	return
}

func (w *CWordLine) HasSpace() bool {
	for _, word := range w.words {
		if word.IsSpace() {
			return true
		}
	}
	return false
}

func (w *CWordLine) Value() (s string) {
	for _, c := range w.words {
		s += c.Value()
	}
	return
}

func (w *CWordLine) String() (s string) {
	s = "{"
	for i, c := range w.words {
		if i > 0 {
			s += ","
		}
		s += c.String()
	}
	s += "}"
	return
}

// wrap, justify and align the set input, with filler style
func (w *CWordLine) Make(wrap WrapMode, justify Justification, maxChars int, fillerStyle Style) (lines []WordLine) {
	lines = append(lines, NewEmptyWordLine())
	cid, wid, lid := 0, 0, 0
	for _, word := range w.words {
		for _, c := range word.Characters() {
			switch c.Value() {
			case '\n':
				lines = append(lines, NewEmptyWordLine())
				lid = len(lines) - 1
				wid = -1
			default:
				if wid >= lines[lid].Len() {
					lines[lid].AppendWordCell(NewEmptyWordCell())
				}
				lines[lid].AppendWordRune(wid, c.Value(), c.Style())
			}
			cid++
		}
		wid++
	}
	lines = w.applyTypography(wrap, justify, maxChars, fillerStyle, lines)
	return
}

func (w *CWordLine) applyTypography(wrap WrapMode, justify Justification, maxChars int, fillerStyle Style, input []WordLine) (output []WordLine) {
	output = w.applyTypographicWrap(wrap, maxChars, input)
	output = w.applyTypographicJustify(justify, maxChars, fillerStyle, output)
	return
}

func (w *CWordLine) applyTypographicWrap(wrap WrapMode, maxChars int, input []WordLine) (output []WordLine) {
	// all space-words must be applied as 1 width
	switch wrap {
	case WRAP_WORD:
		// break onto inserted/new line at end gap
		// - if line has no breakpoints, truncate
		output = w.applyTypographicWrapWord(maxChars, input)
	case WRAP_WORD_CHAR:
		// break onto inserted/new line at end gap
		// - if line has no breakpoints, fallthrough
		output = w.applyTypographicWrapWordChar(maxChars, input)
	case WRAP_CHAR:
		// break onto inserted/new line at maxChars
		output = w.applyTypographicWrapChar(maxChars, input)
	case WRAP_NONE:
		// truncate each line to maxChars
		output = w.applyTypographicWrapNone(maxChars, input)
	}
	return
}

func (w *CWordLine) applyTypographicJustify(justify Justification, maxChars int, fillerStyle Style, input []WordLine) (output []WordLine) {
	switch justify {
	case JUSTIFY_FILL:
		// each non-empty line is space-expanded to fill maxChars
		output = w.applyTypographicJustifyFill(maxChars, fillerStyle, input)
	case JUSTIFY_CENTER:
		// each non-empty line is centered on halfway maxChars
		output = w.applyTypographicJustifyCenter(maxChars, fillerStyle, input)
	case JUSTIFY_RIGHT:
		// each non-empty line is left-padded to fill maxChars
		output = w.applyTypographicJustifyRight(maxChars, fillerStyle, input)
	case JUSTIFY_LEFT:
		fallthrough
	default:
		// each non-empty line has leading space removed
		output = w.applyTypographicJustifyLeft(input)
	}
	return
}
