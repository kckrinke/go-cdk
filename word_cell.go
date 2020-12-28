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

// CTextBuffer holds rows of words
// Views can Draw text using this, along with other primitives

type WordCell struct {
	characters []*CTextCell
}

func NewWordCell(word string, style Style) *WordCell {
	wc := &WordCell{
		characters: make([]*CTextCell, len(word)),
	}
	wc.Set(word, style)
	return wc
}

func (w *WordCell) Set(word string, style Style) {
	for i, c := range word {
		w.characters[i] = NewRuneCell(c, style)
	}
}

func (w WordCell) Len() int {
	return len(w.characters)
}

func (w WordCell) String() (s string) {
	s = ""
	for _, c := range w.characters {
		s += c.String()
	}
	return
}
