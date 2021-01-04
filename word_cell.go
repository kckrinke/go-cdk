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

type WordCell struct {
	characters []*CTextCell
}

func NewEmptyWordCell() *WordCell {
	return &WordCell{
		characters: make([]*CTextCell, 0),
	}
}

func NewWordCell(word string, style Style) *WordCell {
	w := &WordCell{}
	w.Set(word, style)
	return w
}

func (w WordCell) Characters() []*CTextCell {
	return w.characters
}

func (w *WordCell) Set(word string, style Style) {
	w.characters = make([]*CTextCell, len(word))
	for i, c := range word {
		w.characters[i] = NewRuneCell(c, style)
	}
	return
}

func (w *WordCell) GetCharacter(index int) (char *CTextCell) {
	if index < len(w.characters) {
		char = w.characters[index]
	}
	return
}

func (w *WordCell) AppendRune(r rune, style Style) {
	w.characters = append(
		w.characters,
		NewRuneCell(r, style),
	)
}

func (w *WordCell) IsSpace() bool {
	for _, c := range w.characters {
		if !c.IsSpace() {
			return false
		}
	}
	return true
}

func (w *WordCell) HasSpace() bool {
	for _, c := range w.characters {
		if c.IsSpace() {
			return true
		}
	}
	return false
}

func (w WordCell) Len() (count int) {
	count = 0
	for _, c := range w.characters {
		count += c.Width()
	}
	return
}

func (w WordCell) Value() (word string) {
	word = ""
	for _, c := range w.characters {
		word += string(c.Value())
	}
	return
}

func (w WordCell) String() (s string) {
	s = ""
	for _, c := range w.characters {
		s += c.String()
	}
	return
}
