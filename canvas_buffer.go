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

// provide an underlying buffer for Canvases
type CanvasBuffer interface {
	String() string
	Size() (size Rectangle)
	Width() (width int)
	Height() (height int)
	Resize(size Rectangle, style Style)
	Cell(x int, y int) TextCell
	GetDim(x, y int) bool
	GetBgColor(x, y int) (bg Color)
	GetContent(x, y int) (textCell TextCell)
	SetContent(x int, y int, r rune, style Style) error
	LoadData(d [][]TextCell)
}

// concrete implementation of the CanvasBuffer interface
type CCanvasBuffer struct {
	data  [][]TextCell
	size  Rectangle
	style Style

	sync.Mutex
}

// construct a new canvas buffer
func NewCanvasBuffer(size Rectangle, style Style) CanvasBuffer {
	size.Floor(0, 0)
	b := &CCanvasBuffer{
		data: make([][]TextCell, size.W),
		size: MakeRectangle(0, 0),
	}
	b.Resize(size, style)
	return b
}

// return a string describing the buffer, only useful for debugging purposes
func (b *CCanvasBuffer) String() string {
	return fmt.Sprintf(
		"{Size=%s}",
		b.size,
	)
}

// return the rectangle size of the buffer
func (b *CCanvasBuffer) Size() (size Rectangle) {
	return b.size
}

// return just the width of the buffer
func (b *CCanvasBuffer) Width() (width int) {
	return b.size.W
}

// return just the height of the buffer
func (b *CCanvasBuffer) Height() (height int) {
	return b.size.H
}

// resize the buffer
func (b *CCanvasBuffer) Resize(size Rectangle, style Style) {
	b.Lock()
	defer b.Unlock()
	size.Floor(0, 0)
	if b.size.W == size.W && b.size.H == size.H {
		return
	}
	for x := 0; x < size.W; x++ {
		if len(b.data) <= x {
			b.data = append(b.data, make([]TextCell, size.H))
		}
		for y := 0; y < size.H; y++ {
			if len(b.data[x]) <= y {
				b.data[x] = append(b.data[x], NewRuneCell(' ', style))
			} else if b.data[x][y] == nil {
				b.data[x][y] = NewRuneCell(' ', style)
			}
		}
	}
	if b.size.W > size.W {
		b.data = b.data[:size.W]
	}
	if b.size.H > size.H {
		for x := 0; x < size.W; x++ {
			if len(b.data) <= x {
				b.data = append(b.data, make([]TextCell, size.H))
			}
			if len(b.data[x]) >= size.H {
				b.data[x] = b.data[x][:size.H]
			}
		}
	}
	b.size = size
}

// return the text cell at the given coordinates, nil if not found
func (b *CCanvasBuffer) Cell(x int, y int) TextCell {
	if x >= 0 && y >= 0 && x < b.size.W && y < b.size.H {
		return b.data[x][y]
	}
	return nil
}

// return true if the given coordinates are styled 'dim', false otherwise
func (b *CCanvasBuffer) GetDim(x, y int) bool {
	c := b.GetContent(x, y)
	s := c.Style()
	_, _, a := s.Decompose()
	return a.IsDim()
}

// return the background color at the given coordinates
func (b *CCanvasBuffer) GetBgColor(x, y int) (bg Color) {
	c := b.GetContent(x, y)
	s := c.Style()
	_, bg, _ = s.Decompose()
	return
}

// convenience method, returns the results of calling Cell() with the given
// coordinates
func (b *CCanvasBuffer) GetContent(x, y int) (textCell TextCell) {
	textCell = b.Cell(x, y)
	return
}

// set the cell content at the given coordinates
func (b *CCanvasBuffer) SetContent(x int, y int, r rune, style Style) error {
	dLen := len(b.data)
	if x >= 0 && x < dLen {
		dxLen := len(b.data[x])
		if y >= 0 && y < dxLen {
			b.data[x][y].Set(r)
			b.data[x][y].SetStyle(style)
			return nil
		}
		return fmt.Errorf("y=%v not in range [0-%d]", y, len(b.data[x])-1)
	}
	return fmt.Errorf("x=%v not in range [0-%d]", x, len(b.data)-1)
}

// given matrix array of text cells, load that data in this canvas space
func (b *CCanvasBuffer) LoadData(d [][]TextCell) {
	for x := 0; x < len(d); x++ {
		for y := 0; y < len(d[x]); y++ {
			if y >= len(b.data[x]) {
				b.data[x] = append(b.data[x], NewRuneCell(d[x][y].Value(), d[x][y].Style()))
			} else {
				b.data[x][y].Set(d[x][y].Value())
			}
		}
	}
}
