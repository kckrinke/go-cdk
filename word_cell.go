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

// WordCell holds a list of characters making up a word or a gap (space)

type WordCell interface {
	Characters() []TextCell
	Set(word string, style Style)
	GetCharacter(index int) (char TextCell)
	AppendRune(r rune, style Style)
	IsNil() bool
	IsSpace() bool
	HasSpace() bool
	Len() (count int)
	CompactLen() (count int)
	Value() (word string)
	String() (s string)
}

type CWordCell struct {
	characters []TextCell
}

func NewEmptyWordCell() WordCell {
	return &CWordCell{
		characters: make([]TextCell, 0),
	}
}

func NewNilWordCell(style Style) WordCell {
	return &CWordCell{
		characters: []TextCell{NewRuneCell(rune(0), style)},
	}
}

func NewWordCell(word string, style Style) WordCell {
	w := &CWordCell{}
	w.Set(word, style)
	return w
}

func (w *CWordCell) Characters() []TextCell {
	return w.characters
}

func (w *CWordCell) Set(word string, style Style) {
	w.characters = make([]TextCell, len(word))
	for i, c := range word {
		w.characters[i] = NewRuneCell(c, style)
	}
	return
}

func (w *CWordCell) GetCharacter(index int) (char TextCell) {
	if index < len(w.characters) {
		char = w.characters[index]
	}
	return
}

func (w *CWordCell) AppendRune(r rune, style Style) {
	w.characters = append(
		w.characters,
		NewRuneCell(r, style),
	)
}

func (w *CWordCell) IsNil() bool {
	for _, c := range w.characters {
		if !c.IsNil() {
			return false
		}
	}
	return true
}

func (w *CWordCell) IsSpace() bool {
	for _, c := range w.characters {
		if !c.IsSpace() {
			return false
		}
	}
	return true
}

func (w *CWordCell) HasSpace() bool {
	for _, c := range w.characters {
		if c.IsSpace() {
			return true
		}
	}
	return false
}

// the total number of characters in this word
func (w *CWordCell) Len() (count int) {
	count = 0
	for _, c := range w.characters {
		count += c.Width()
	}
	return
}

// same as `Len()` with space-words being treated as 1 character wide rather
// than the literal number of spaces from the input string
func (w *CWordCell) CompactLen() (count int) {
	if w.IsSpace() {
		count = 1
		return
	}
	count = w.Len()
	return
}

// returns the literal string value of the word
func (w *CWordCell) Value() (word string) {
	word = ""
	for _, c := range w.characters {
		word += string(c.Value())
	}
	return
}

// returns the debuggable value of the word
func (w *CWordCell) String() (s string) {
	s = ""
	for _, c := range w.characters {
		s += c.String()
	}
	return
}
