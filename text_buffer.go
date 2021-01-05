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
)

var (
	TapSpace = "    "
)

type TextBuffer interface {
	Set(input string, style Style)
	CharacterCount() (cellCount int)
	WordCount() (wordCount int)
	Draw(canvas *Canvas, singleLineMode bool, wordWrap WrapMode, justify Justification, align VerticalAlignment) EventFlag
}

type CTextBuffer struct {
	input *WordLine
	style Style

	sync.Mutex
}

func NewEmptyTextBuffer(style Style) *CTextBuffer {
	return &CTextBuffer{
		style: style,
	}
}

func NewTextBuffer(input string, style Style) *CTextBuffer {
	tb := &CTextBuffer{
		style: style,
	}
	tb.Set(input, style)
	return tb
}

func (b *CTextBuffer) Set(input string, style Style) {
	b.input = NewWordLine(input, style)
}

func (b *CTextBuffer) CharacterCount() (cellCount int) {
	cellCount = b.input.CharacterCount()
	return
}

func (b *CTextBuffer) WordCount() (wordCount int) {
	wordCount = b.input.WordCount()
	return
}

func (b *CTextBuffer) Draw(canvas *Canvas, singleLine bool, wordWrap WrapMode, justify Justification, align VerticalAlignment) EventFlag {
	if b.input.CharacterCount() == 0 {
		// non-operation
		return EVENT_PASS
	}

	maxChars := canvas.size.W
	lines := b.input.Make(wordWrap, justify, maxChars)

	for y, line := range lines {
		x := 0
		for _, word := range line.Words() {
			for _, c := range word.Characters() {
				canvas.SetRune(x, y, c.Value(), c.Style())
				x++
			}
		}
	}

	return EVENT_STOP
}
