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

	"github.com/kckrinke/go-cdk/utils"
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

func (w WordLine) MakeFill(maxChars int) (line *WordLine) {
	line = NewEmptyWordLine()
	if maxChars > w.CharacterCount() {
		// find gaps and increase them by duplicating cells
		var gapSizes []int
		sumGaps := 0
		trackGaps := []int{}
		for _, word := range w.words {
			line.AppendWordCell(word)
			if word.IsSpace() {
				gapSizes = append(gapSizes, word.Len())
				trackGaps = append(trackGaps, word.Len())
				sumGaps += word.Len()
			} else {
				gapSizes = append(gapSizes, 0)
			}
		}
		delta := maxChars - line.CharacterCount()
		trackGaps = utils.DistInts(delta+sumGaps, trackGaps)
		gapIndex := 0
		for i, v := range gapSizes {
			if v > 0 {
				gapSizes[i] = trackGaps[gapIndex]
				gapIndex++
			}
		}
		gapWords := w.GapWordList()
		// this is all weird lol
		for i, word := range line.words {
			if gapWords[i] != nil {
				origWordLen := word.Len()
				currentCharacterIndex := 0
				sizeRequested := gapSizes[i]
				for word.Len() < sizeRequested {
					if currentCharacterIndex >= origWordLen {
						currentCharacterIndex = 0
					}
					c := word.GetCharacter(currentCharacterIndex)
					word.AppendRune(c.Value(), c.Style())
					currentCharacterIndex++
				}
			}
		}
	} else {
		// truncate justify-left
		line = w.MakeLeft(maxChars)
	}
	return
}

func (w WordLine) MakeCenter(maxChars int, indentStyle Style) (line *WordLine) {
	line = NewEmptyWordLine()
	halfLine := maxChars / 2
	halfText := w.CharacterCount() / 2
	if maxChars > w.CharacterCount() {
		// center
		delta := halfLine - halfText
		indent := NewEmptyWordCell()
		for i := 0; i < delta; i++ {
			indent.AppendRune(' ', indentStyle)
		}
		line.AppendWordCell(indent)
		for _, word := range w.words {
			line.AppendWordCell(word)
		}
	} else {
		// truncate-center
		delta := halfText - halfLine
		count, wid := 0, 0
		for _, word := range w.words {
			if count > delta {
				line.words = append(line.words, NewEmptyWordCell())
				wid = len(line.words) - 1
			}
			for _, c := range word.characters {
				if count > delta {
					line.words[wid].AppendRune(c.Value(), c.Style())
				}
				count++
				if count > maxChars {
					// inner
					break
				}
			}
			if count > maxChars {
				// outer
				break
			}
		}
	}
	return
}

func (w WordLine) MakeRight(maxChars int, indentStyle Style) (line *WordLine) {
	line = NewEmptyWordLine()
	if maxChars > w.CharacterCount() {
		// indent
		delta := maxChars - w.CharacterCount()
		indent := NewEmptyWordCell()
		for i := 0; i < delta; i++ {
			indent.AppendRune(' ', indentStyle)
		}
		line.AppendWordCell(indent)
		for _, word := range w.words {
			line.AppendWordCell(word)
		}
	} else {
		// truncate-left
		delta := w.CharacterCount() - maxChars
		count, wid := 0, 0
		for _, word := range w.words {
			if count > delta {
				line.words = append(line.words, NewEmptyWordCell())
				wid = len(line.words) - 1
			}
			for _, c := range word.characters {
				if count > delta {
					line.words[wid].AppendRune(c.Value(), c.Style())
				}
				count++
			}
		}
	}
	return
}

func (w WordLine) MakeLeft(maxChars int) (line *WordLine) {
	line = NewEmptyWordLine()
	for _, word := range w.words {
		if maxChars > line.CharacterCount()+word.Len() {
			line.AppendWordCell(word)
		} else {
			partial := NewEmptyWordCell()
			delta := line.CharacterCount() + word.Len() - maxChars
			for i := 0; i < delta; i++ {
				c := word.GetCharacter(i)
				partial.AppendRune(c.Value(), c.Style())
			}
			line.AppendWordCell(partial)
			break
		}
	}
	return
}

func (w WordLine) MakeTruncated(maxChars int, onSpace bool) (line *WordLine, didTruncate bool) {
	line = NewEmptyWordLine()
	for _, word := range w.words {
		if maxChars > line.CharacterCount()+word.Len() {
			line.AppendWordCell(word)
		} else {
			didTruncate = true
			if onSpace && word.IsSpace() {
				// this word is a gap, we can stop here
				break
			} else if onSpace {
				// word is real, let's remove the last space?
				foundSpace := false
				for !foundSpace {
					lastIndex := line.Len() - 1
					lastWord := line.GetWord(lastIndex)
					if lastWord.IsSpace() {
						foundSpace = true
					}
					line.RemoveWord(lastIndex)
				}
				break
			}
			partial := NewEmptyWordCell()
			delta := maxChars - line.CharacterCount()
			for i := 0; i < delta; i++ {
				c := word.GetCharacter(i)
				partial.AppendRune(c.Value(), c.Style())
			}
			line.AppendWordCell(partial)
			break
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
