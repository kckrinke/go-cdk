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
)

type WordLine struct {
	words []*WordCell
}

func NewWordLine(line string, style Style) *WordLine {
	words := strings.Fields(line)
	wl := &WordLine{
		words: make([]*WordCell, len(words)),
	}
	for i, word := range words {
		word, _ := NewWordCell(word, style)
		wl.words[i] = word
	}
	return wl
}

func (w WordLine) Words() []*WordCell {
	return w.words
}

func (w WordLine) LetterCount(spaces bool) int {
	c := 0
	for i, word := range w.words {
		if i != 0 && spaces {
			c += 1
		}
		c += word.Len()
	}
	return c
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
