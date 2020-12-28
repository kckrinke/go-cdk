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
	vb := &CanvasBuffer{
		data: make([][]*CTextCell, size.W),
		size: Rectangle{0, 0},
	}
	vb.style = style
	vb.Resize(size)
	return vb
}

func (vb CanvasBuffer) String() string {
	return fmt.Sprintf(
		"{Size=%s,Style=%s}",
		vb.size,
		vb.style.String(),
	)
}

func (vb *CanvasBuffer) Resize(size Rectangle) {
	vb.Lock()
	defer vb.Unlock()
	for x := 0; x < size.W; x++ {
		if len(vb.data) <= x {
			vb.data = append(vb.data, make([]*CTextCell, size.H))
		}
		for y := 0; y < size.H; y++ {
			if len(vb.data[x]) <= y {
				vb.data[x] = append(vb.data[x], NewRuneCell(' ', vb.style))
			}
		}
	}
	if vb.size.W > size.W {
		vb.data = vb.data[:size.W]
	}
	if vb.size.H > size.H {
		for x := 0; x < size.W; x++ {
			vb.data[x] = vb.data[x][:size.H]
		}
	}
	vb.size = size
}

func (vb *CanvasBuffer) Cell(x int, y int) *CTextCell {
	if vb.size.W > x && vb.size.H > y {
		return vb.data[x][y]
	}
	return nil
}

func (vb *CanvasBuffer) GetDim(x, y int) bool {
	_, s, _ := vb.GetContent(x, y)
	_, _, a := s.Decompose()
	return a.IsDim()
}

func (vb *CanvasBuffer) GetBgColor(x, y int) (bg Color) {
	_, s, _ := vb.GetContent(x, y)
	_, bg, _ = s.Decompose()
	return
}

func (vb *CanvasBuffer) GetContent(x, y int) (mainc rune, style Style, width int) {
	if x >= 0 && y >= 0 && x < vb.size.W && y < vb.size.H {
		c := vb.data[x][y]
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

func (vb *CanvasBuffer) SetContent(x int, y int, mainc rune, style Style) error {
	dlen := len(vb.data)
	if x >= 0 && x < dlen {
		dxlen := len(vb.data[x])
		if y >= 0 && y < dxlen {
			vb.data[x][y].Set(mainc)
			vb.data[x][y].SetStyle(style)
			return nil
		}
		return fmt.Errorf("y=%v not in range [0-%d]", y, len(vb.data[x])-1)
	}
	return fmt.Errorf("x=%v not in range [0-%d]", x, len(vb.data)-1)
}

func (vb *CanvasBuffer) LoadData(d [][]*CTextCell) {
	for x := 0; x < len(d); x++ {
		for y := 0; y < len(d[x]); y++ {
			if y >= len(vb.data[x]) {
				vb.data[x] = append(vb.data[x], NewRuneCell(d[x][y].Value(), d[x][y].Style()))
			} else {
				vb.data[x][y].Set(d[x][y].Value())
			}
		}
	}
}
