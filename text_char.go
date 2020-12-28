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

func NewTextChar(b []byte) CTextChar {
	r, s := utf8.DecodeRune(b)
	return CTextChar{
		value: r,
		width: s,
	}
}

func (c *CTextChar) Set(r rune) {
	c.SetByte([]byte(string(r)))
}

func (c *CTextChar) SetByte(b []byte) {
	c.Lock()
	defer c.Unlock()
	c.value, c.width = utf8.DecodeRune(b)
}

func (c CTextChar) Width() int {
	c.RLock()
	defer c.RUnlock()
	return c.width
}

func (c CTextChar) Value() rune {
	c.RLock()
	defer c.RUnlock()
	return c.value
}

func (c CTextChar) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("%c", c.value)
}

func (c CTextChar) IsSpace() bool {
	return unicode.IsSpace(c.Value())
}
