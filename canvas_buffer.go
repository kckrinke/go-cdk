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

type CanvasBuffer struct {
	data  [][]*CTextCell
	size  Rectangle
	style Style

	sync.Mutex
}

func NewCanvasBuffer(size Rectangle, style Style) *CanvasBuffer {
	b := &CanvasBuffer{
		data: make([][]*CTextCell, size.W),
		size: MakeRectangle(0, 0),
	}
	b.style = style
	b.Resize(size)
	return b
}

func (b *CanvasBuffer) String() string {
	return fmt.Sprintf(
		"{Size=%s,Style=%s}",
		b.size,
		b.style.String(),
	)
}

func (b *CanvasBuffer) Resize(size Rectangle) {
	b.Lock()
	defer b.Unlock()
	for x := 0; x < size.W; x++ {
		if len(b.data) <= x {
			b.data = append(b.data, make([]*CTextCell, size.H))
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

func (b *CanvasBuffer) Cell(x int, y int) *CTextCell {
	if b.size.W > x && b.size.H > y {
		return b.data[x][y]
	}
	return nil
}

func (b *CanvasBuffer) GetDim(x, y int) bool {
	_, s, _ := b.GetContent(x, y)
	_, _, a := s.Decompose()
	return a.IsDim()
}

func (b *CanvasBuffer) GetBgColor(x, y int) (bg Color) {
	_, s, _ := b.GetContent(x, y)
	_, bg, _ = s.Decompose()
	return
}

func (b *CanvasBuffer) GetContent(x, y int) (mainc rune, style Style, width int) {
	if x >= 0 && y >= 0 && x < b.size.W && y < b.size.H {
		c := b.data[x][y]
		c.Lock()
		mainc, style = c.Value(), c.Style()
		if width = c.Width(); width == 0 || mainc < ' ' {
			width = 1
			mainc = ' '
		}
		c.Unlock()
	}
	return
}

func (b *CanvasBuffer) SetContent(x int, y int, mainc rune, style Style) error {
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

func (b *CanvasBuffer) LoadData(d [][]*CTextCell) {
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
