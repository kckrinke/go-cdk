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
	"sync"
)

type TextCell interface {
	Dirty() bool
	Set(r rune)
	SetByte(b []byte)
	SetStyle(style Style)
	Width() int
	Value() rune
	String() string
	Style() Style
	IsSpace() bool
}

type CTextCell struct {
	char  TextChar
	style Style
	dirty bool

	sync.Mutex
}

func NewRuneCell(char rune, style Style) TextCell {
	return NewTextCell(NewTextChar([]byte(string(char))), style)
}

func NewTextCell(char TextChar, style Style) TextCell {
	return &CTextCell{
		char:  char,
		style: style,
		dirty: false,
	}
}

func (t *CTextCell) Dirty() bool {
	return t.dirty
}

func (t *CTextCell) Set(r rune) {
	t.Lock()
	defer t.Unlock()
	t.char.Set(r)
	t.dirty = true
}

func (t *CTextCell) SetByte(b []byte) {
	t.Lock()
	defer t.Unlock()
	t.char.SetByte(b)
	t.dirty = true
}

func (t *CTextCell) SetStyle(style Style) {
	t.Lock()
	defer t.Unlock()
	t.style = style
	t.dirty = true
}

func (t *CTextCell) Width() int {
	return t.char.Width()
}

func (t *CTextCell) Value() rune {
	return t.char.Value()
}

func (t *CTextCell) String() string {
	return fmt.Sprintf(
		"{Char=%s,Style=%s}",
		t.char.String(),
		t.style.String(),
	)
}

func (t *CTextCell) Style() Style {
	return t.style
}

func (t *CTextCell) IsSpace() bool {
	return t.char.IsSpace()
}
