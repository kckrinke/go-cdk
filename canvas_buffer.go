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

type CanvasBuffer interface {
	String() string
	Size() (size Rectangle)
	Width() (width int)
	Height() (height int)
	Resize(size Rectangle)
	Cell(x int, y int) TextCell
	GetDim(x, y int) bool
	GetBgColor(x, y int) (bg Color)
	GetContent(x, y int) (textCell TextCell)
	SetContent(x int, y int, mainc rune, style Style) error
	LoadData(d [][]TextCell)
}

type CCanvasBuffer struct {
	data  [][]TextCell
	size  Rectangle
	style Style

	sync.Mutex
}

func NewCanvasBuffer(size Rectangle, style Style) CanvasBuffer {
	b := &CCanvasBuffer{
		data: make([][]TextCell, size.W),
		size: MakeRectangle(0, 0),
	}
	b.style = style
	b.Resize(size)
	return b
}

func (b *CCanvasBuffer) String() string {
	return fmt.Sprintf(
		"{Size=%s,Style=%s}",
		b.size,
		b.style.String(),
	)
}

func (b *CCanvasBuffer) Size() (size Rectangle) {
	return b.size
}

func (b *CCanvasBuffer) Width() (width int) {
	return b.size.W
}

func (b *CCanvasBuffer) Height() (height int) {
	return b.size.H
}

func (b *CCanvasBuffer) Resize(size Rectangle) {
	b.Lock()
	defer b.Unlock()
	for x := 0; x < size.W; x++ {
		if len(b.data) <= x {
			b.data = append(b.data, make([]TextCell, size.H))
		}
		for y := 0; y < size.H; y++ {
			if len(b.data[x]) <= y {
				b.data[x] = append(b.data[x], NewRuneCell(' ', b.style))
			}
		}
	}
	if b.size.W > size.W {
		b.data = b.data[:size.W]
	}
	if b.size.H > size.H {
		for x := 0; x < size.W; x++ {
			b.data[x] = b.data[x][:size.H]
		}
	}
	b.size = size
}

func (b *CCanvasBuffer) Cell(x int, y int) TextCell {
	if b.size.W > x && b.size.H > y {
		return b.data[x][y]
	}
	return nil
}

func (b *CCanvasBuffer) GetDim(x, y int) bool {
	c := b.GetContent(x, y)
	s := c.Style()
	_, _, a := s.Decompose()
	return a.IsDim()
}

func (b *CCanvasBuffer) GetBgColor(x, y int) (bg Color) {
	c := b.GetContent(x, y)
	s := c.Style()
	_, bg, _ = s.Decompose()
	return
}

func (b *CCanvasBuffer) GetContent(x, y int) (textCell TextCell) {
	if x >= 0 && y >= 0 && x < b.size.W && y < b.size.H {
		textCell = b.data[x][y]
	}
	return
}

func (b *CCanvasBuffer) SetContent(x int, y int, mainc rune, style Style) error {
	dlen := len(b.data)
	if x >= 0 && x < dlen {
		dxlen := len(b.data[x])
		if y >= 0 && y < dxlen {
			b.data[x][y].Set(mainc)
			b.data[x][y].SetStyle(style)
			return nil
		}
		return fmt.Errorf("y=%v not in range [0-%d]", y, len(b.data[x])-1)
	}
	return fmt.Errorf("x=%v not in range [0-%d]", x, len(b.data)-1)
}

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
