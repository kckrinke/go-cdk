// Copyright 2019 The TCell Authors
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

	"github.com/mattn/go-runewidth"
)

type cell struct {
	currMain  rune
	currComb  []rune
	currStyle Style
	lastMain  rune
	lastStyle Style
	lastComb  []rune
	width     int

	mutex *sync.Mutex
}

func newCell() *cell {
	c := &cell{}
	_ = c.init()
	return c
}

func (c *cell) init() bool {
	if c.mutex == nil {
		c.currComb = make([]rune, 0)
		c.lastComb = make([]rune, 0)
		c.mutex = &sync.Mutex{}
		return true
	}
	return false
}

func (c *cell) lock() {
	c.mutex.Lock()
}

func (c *cell) unlock() {
	c.mutex.Unlock()
}

// CellBuffer represents a two dimensional array of character cells.
// This is primarily intended for use by Screen implementors; it
// contains much of the common code they need.  To create one, just
// declare a variable of its type; no explicit initialization is necessary.
//
// CellBuffer should be thread safe, original tcell is not.
type CellBuffer struct {
	w     int
	h     int
	cells []*cell

	mutex *sync.Mutex
}

func NewCellBuffer() *CellBuffer {
	cb := &CellBuffer{}
	cb.init()
	return cb
}

func (cb *CellBuffer) init() bool {
	if cb.mutex == nil {
		cb.cells = make([]*cell, 0)
		cb.mutex = &sync.Mutex{}
		return true
	}
	return false
}

func (cb *CellBuffer) lock() {
	cb.mutex.Lock()
}

func (cb *CellBuffer) unlock() {
	cb.mutex.Unlock()
}

// SetContent sets the contents (primary rune, combining runes,
// and style) for a cell at a given location.
func (cb *CellBuffer) SetContent(x int, y int, mainc rune, combc []rune, style Style) {
	Tracef("x=%d, y=%d, rune=%v, style=%v", x, y, mainc, style)
	cb.lock()
	defer cb.unlock()
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := cb.cells[(y*cb.w)+x]
		c.lock()
		c.currComb = append([]rune{}, combc...)
		if c.currMain != mainc {
			c.width = runewidth.RuneWidth(mainc)
		}
		c.currMain = mainc
		c.currStyle = style
		c.unlock()
	}
}

// GetContent returns the contents of a character cell, including the
// primary rune, any combining character runes (which will usually be
// nil), the style, and the display width in cells.  (The width can be
// either 1, normally, or 2 for East Asian full-width characters.)
func (cb *CellBuffer) GetContent(x, y int) (mainc rune, combc []rune, style Style, width int) {
	cb.lock()
	defer cb.unlock()
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := cb.cells[(y*cb.w)+x]
		c.lock()
		mainc, combc, style = c.currMain, c.currComb, c.currStyle
		if width = c.width; width == 0 || mainc < ' ' {
			width = 1
			mainc = ' '
		}
		c.unlock()
	}
	return
}

// Size returns the (width, height) in cells of the buffer.
func (cb *CellBuffer) Size() (w, h int) {
	cb.lock()
	defer cb.unlock()
	w, h = cb.w, cb.h
	return
}

// Invalidate marks all characters within the buffer as dirty.
func (cb *CellBuffer) Invalidate() {
	Tracef("called")
	cb.lock()
	defer cb.unlock()
	for i := range cb.cells {
		cb.cells[i].lock()
		cb.cells[i].lastMain = rune(0)
		cb.cells[i].unlock()
	}
}

// Dirty checks if a character at the given location needs an
// to be refreshed on the physical display.  This returns true
// if the cell content is different since the last time it was
// marked clean.
func (cb *CellBuffer) Dirty(x, y int) bool {
	cb.lock()
	defer cb.unlock()
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := cb.cells[(y*cb.w)+x]
		if c.lastMain == rune(0) {
			return true
		}
		if c.lastMain != c.currMain {
			return true
		}
		if c.lastStyle != c.currStyle {
			return true
		}
		if len(c.lastComb) != len(c.currComb) {
			return true
		}
		for i := range c.lastComb {
			if c.lastComb[i] != c.currComb[i] {
				return true
			}
		}
	}
	return false
}

// SetDirty is normally used to indicate that a cell has
// been displayed (in which case dirty is false), or to manually
// force a cell to be marked dirty.
func (cb *CellBuffer) SetDirty(x, y int, dirty bool) {
	Tracef("x=%d, y=%d, dirty=%v", x, y, dirty)
	cb.lock()
	defer cb.unlock()
	if x >= 0 && y >= 0 && x < cb.w && y < cb.h {
		c := cb.cells[(y*cb.w)+x]
		c.lock()
		if dirty {
			c.lastMain = rune(0)
		} else {
			if c.currMain == rune(0) {
				c.currMain = ' '
			}
			c.lastMain = c.currMain
			c.lastComb = c.currComb
			c.lastStyle = c.currStyle
		}
		c.unlock()
	}
}

// Resize is used to resize the cells array, with different dimensions,
// while preserving the original contents.  The cells will be invalidated
// so that they can be redrawn.
func (cb *CellBuffer) Resize(w, h int) {
	Debugf("w=%d, h=%d", w, h)
	if cb.h == h && cb.w == w {
		return
	}
	cb.lock()
	defer cb.unlock()
	if w == 0 || h == 0 {
		cb.cells = make([]*cell, 0)
		return
	}
	newc := make([]*cell, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			nc := newCell()
			k := (y * cb.w) + x
			if len(cb.cells) > k {
				oc := cb.cells[k]
				oc.lock()
				if oc != nil {
					nc.currMain = oc.currMain
					nc.currComb = oc.currComb
					nc.currStyle = oc.currStyle
					nc.width = oc.width
				}
				oc.unlock()
			}
			nc.lastMain = rune(0)
			newc[(y*w)+x] = nc
		}
	}
	cb.cells = newc
	cb.h = h
	cb.w = w
}

// Fill fills the entire cell buffer array with the specified character
// and style.  Normally choose ' ' to clear the screen.  This API doesn't
// support combining characters, or characters with a width larger than one.
func (cb *CellBuffer) Fill(r rune, style Style) {
	Tracef("rune=%v, style=%v", r, style)
	cb.lock()
	defer cb.unlock()
	for _, c := range cb.cells {
		c.lock()
		c.currMain = r
		c.currComb = nil
		c.currStyle = style
		c.width = 1
		c.unlock()
	}
}
