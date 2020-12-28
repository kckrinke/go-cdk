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
	char  CTextChar
	style Style
	dirty bool

	sync.Mutex
}

func NewByteCell(char []byte, style Style) *CTextCell {
	return NewTextCell(NewTextChar(char), style)
}

func NewStringCell(char string, style Style) *CTextCell {
	return NewTextCell(NewTextChar([]byte(char)), style)
}

func NewRuneCell(char rune, style Style) *CTextCell {
	return NewTextCell(NewTextChar([]byte(string(char))), style)
}

func NewTextCell(char CTextChar, style Style) *CTextCell {
	return &CTextCell{
		char:  char,
		style: style,
		dirty: false,
	}
}

func (t *CTextCell) Dirty() bool {
	return t.Dirty()
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
