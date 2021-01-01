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
	"unicode"
	"unicode/utf8"
)

type TextChar interface {
	Set(r rune)
	SetByte(b []byte)
	Width() int
	Value() rune
	String() string
	IsSpace() bool
}

type CTextChar struct {
	value rune
	width int

	sync.RWMutex
}

func NewTextChar(b []byte) *CTextChar {
	r, s := utf8.DecodeRune(b)
	return &CTextChar{
		value: r,
		width: s,
	}
}

func (c *CTextChar) Set(r rune) {
	c.SetByte([]byte(string(r)))
}

func (c *CTextChar) SetByte(b []byte) {
	c.value, c.width = utf8.DecodeRune(b)
}

func (c *CTextChar) Width() int {
	return c.width
}

func (c *CTextChar) Value() rune {
	return c.value
}

func (c *CTextChar) String() string {
	return fmt.Sprintf("%c", c.value)
}

func (c *CTextChar) IsSpace() bool {
	return unicode.IsSpace(c.Value())
}
