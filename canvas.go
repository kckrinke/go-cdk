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
	"unicode/utf8"

	"github.com/kckrinke/go-cdk/utils"
)

// a Canvas is the primary means of drawing to the terminal display within CDK
type Canvas interface {
	String() string
	Resize(size Rectangle)
	GetContent(x, y int) (textCell TextCell)
	SetContent(x, y int, char string, s Style) error
	SetRune(x, y int, r rune, s Style) error
	SetOrigin(origin Point2I)
	GetOrigin() Point2I
	GetSize() Rectangle
	Width() (width int)
	Height() (height int)
	SetTheme(style Theme)
	GetTheme() Theme
	SetFill(fill rune)
	GetFill() rune
	Equals(onlyDirty bool, v Canvas) bool
	Composite(v Canvas) error
	Render(display Display) error
	ForEach(fn CanvasForEachFn) EventFlag
	DrawText(pos Point2I, size Rectangle, justify Justification, singleLineMode bool, wrap WrapMode, style Style, markup bool, text string)
	DrawSingleLineText(pos Point2I, maxChars int, justify Justification, style Style, markup bool, text string)
	DrawLine(pos Point2I, length int, orient Orientation, style Style)
	DrawHorizontalLine(pos Point2I, length int, style Style)
	DrawVerticalLine(pos Point2I, length int, style Style)
	Box(pos Point2I, size Rectangle, border bool, fill bool, theme Theme)
	DebugBox(color Color, format string, argv ...interface{})
	Fill(s Theme)
	FillBorder(dim bool, border bool)
	FillBorderTitle(dim bool, title string, justify Justification)
}

// concrete implementation of the Canvas interface
type CCanvas struct {
	buffer CanvasBuffer
	origin Point2I
	size   Rectangle
	theme  Theme
	fill   rune
}

// create a new canvas object with the given origin point, size and theme
func NewCanvas(origin Point2I, size Rectangle, theme Theme) Canvas {
	return &CCanvas{
		buffer: NewCanvasBuffer(size, theme.Normal),
		origin: origin,
		size:   size,
		theme:  theme,
		fill:   ' ',
	}
}

// return a string describing the canvas metadata, useful for debugging
func (c CCanvas) String() string {
	return fmt.Sprintf(
		"{Origin=%s,Size=%s,Theme=%s,Fill=%v,Buffer=%v}",
		c.origin,
		c.size,
		c.theme,
		c.fill,
		c.buffer.String(),
	)
}

// change the size of the canvas, not recommended to do this in practice
func (c *CCanvas) Resize(size Rectangle) {
	c.buffer.Resize(size)
	c.size = size
}

// get the text cell at the given coordinates
func (c *CCanvas) GetContent(x, y int) (textCell TextCell) {
	return c.buffer.GetContent(x, y)
}

// from the given string, set the character and style of the cell at the given
// coordinates. note that only the first UTF-8 byte is used
func (c *CCanvas) SetContent(x, y int, char string, s Style) error {
	r, _ := utf8.DecodeRune([]byte(char))
	return c.buffer.SetContent(x, y, r, s)
}

// set the rune and the style of the cell at the given coordinates
func (c *CCanvas) SetRune(x, y int, r rune, s Style) error {
	return c.buffer.SetContent(x, y, r, s)
}

// set the origin (top-left corner) position of the canvas, used when
// compositing one canvas with another
func (c *CCanvas) SetOrigin(origin Point2I) {
	c.origin = origin
}

// get the origin point of the canvas
func (c *CCanvas) GetOrigin() Point2I {
	return c.origin
}

// get the rectangle size of the canvas
func (c *CCanvas) GetSize() Rectangle {
	return c.size
}

// convenience method to get just the width of the canvas
func (c *CCanvas) Width() (width int) {
	return c.size.W
}

// convenience method to get just the height of the canvas
func (c *CCanvas) Height() (height int) {
	return c.size.H
}

// set the default theme for the canvas, used in cases where there's no inherent
// style available during a fill process for example
func (c *CCanvas) SetTheme(style Theme) {
	c.theme = style
}

// return the default theme for this canvas
func (c *CCanvas) GetTheme() Theme {
	return c.theme
}

// set the rune used to fill in areas, typically a space character " "
func (c *CCanvas) SetFill(fill rune) {
	c.fill = fill
}

// get the rune used to file in areas
func (c *CCanvas) GetFill() rune {
	return c.fill
}

// returns true if the given canvas is painted the same as this one, can compare
// for only cells that were "set" (dirty) or compare every cell of the two
// canvases
func (c *CCanvas) Equals(onlyDirty bool, v Canvas) bool {
	vOrigin := v.GetOrigin()
	vSize := v.GetSize()
	if c.origin.Equals2I(vOrigin) {
		if c.size.EqualsR(vSize) {
			for x := 0; x < vSize.W; x++ {
				for y := 0; y < vSize.H; y++ {
					ca := c.GetContent(x, y)
					va := v.GetContent(x, y)
					if !onlyDirty || (onlyDirty && va.Dirty()) {
						if ca.Style() != va.Style() {
							return false
						}
						if ca.Value() != va.Value() {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

// apply the given canvas to this canvas, at the given one's origin. returns
// an error if the underlying buffer write failed or if the given canvas is
// beyond the bounds of this canvas
func (c *CCanvas) Composite(v Canvas) error {
	vOrigin := v.GetOrigin()
	bSize := c.buffer.Size()
	for y := 0; y < v.Height(); y++ {
		for x := 0; x < v.Width(); x++ {
			cell := v.GetContent(x, y)
			if cell != nil {
				if cell.Dirty() {
					oX, oY := vOrigin.X+x, vOrigin.Y+y
					if oX >= 0 && oX < bSize.W && oY >= 0 && oY < bSize.H {
						if err := c.buffer.SetContent(
							oX,
							oY,
							cell.Value(),
							cell.Style(),
						); err != nil {
							return err
						}
					}
				}
			} else {
				return fmt.Errorf("cell is nil x=%v,y=%v", x, y)
			}
		}
	}
	return nil
}

// render this canvas upon the given display
func (c *CCanvas) Render(display Display) error {
	for x := 0; x < c.size.W; x++ {
		for y := 0; y < c.size.H; y++ {
			cell := c.buffer.Cell(x, y)
			display.SetContent(x, y, cell.Value(), nil, cell.Style())
		}
	}
	return nil
}

// func signature used when iterating over each cell
type CanvasForEachFn = func(x, y int, cell TextCell) EventFlag

// convenience method to iterate of each cell of the canvas, if the given fn
// returns EVENT_STOP then the iteration is halted, otherwise EVEN_PASS will
// allow for the next iteration to proceed
func (c *CCanvas) ForEach(fn CanvasForEachFn) EventFlag {
	for x := 0; x < c.buffer.Width(); x++ {
		for y := 0; y < c.buffer.Height(); y++ {
			if f := fn(x, y, c.buffer.Cell(x, y)); f == EVENT_STOP {
				return EVENT_STOP
			}
		}
	}
	return EVENT_PASS
}

// Write text to the canvas buffer
// origin is the top-left coordinate for the text area being rendered
// alignment is based on origin.X boxed by maxChars or canvas size.W
func (c *CCanvas) DrawText(pos Point2I, size Rectangle, justify Justification, singleLineMode bool, wrap WrapMode, style Style, markup bool, text string) {
	var tb TextBuffer
	if markup {
		m, err := NewMarkup(text, style)
		if err != nil {
			FatalDF(1, "failed to parse markup: %v", err)
		}
		tb = m.TextBuffer()
	} else {
		tb = NewTextBuffer(text, style)
	}
	if size.W == -1 || size.W >= c.size.W {
		size.W = c.size.W
	}
	v := NewCanvas(pos, size, c.theme)
	tb.Draw(v, singleLineMode, wrap, justify, ALIGN_TOP)
	if err := c.Composite(v); err != nil {
		ErrorF("composite error: %v", err)
	}
}

// write a single line of text to the canvas at the given position, of at most
// maxChars, with the text justified and styled. supports Tango markup content
func (c *CCanvas) DrawSingleLineText(position Point2I, maxChars int, justify Justification, style Style, markup bool, text string) {
	c.DrawText(position, MakeRectangle(maxChars, 1), justify, true, WRAP_NONE, style, markup, text)
}

// draw a line vertically or horizontally with the given style
func (c *CCanvas) DrawLine(pos Point2I, length int, orient Orientation, style Style) {
	TraceF("c.DrawLine(%v,%v,%v,%v)", pos, length, orient, style)
	switch orient {
	case ORIENTATION_HORIZONTAL:
		c.DrawHorizontalLine(pos, length, style)
	case ORIENTATION_VERTICAL:
		c.DrawVerticalLine(pos, length, style)
	}
}

// convenience method to draw a horizontal line
func (c *CCanvas) DrawHorizontalLine(pos Point2I, length int, style Style) {
	length = utils.ClampI(length, pos.X, c.size.W-pos.X)
	end := pos.X + length
	for i := pos.X; i < end; i++ {
		c.SetRune(i, pos.Y, RuneHLine, style)
	}
}

// convenience method to draw a vertical line
func (c *CCanvas) DrawVerticalLine(pos Point2I, length int, style Style) {
	length = utils.ClampI(length, pos.Y, c.size.H-pos.Y)
	end := pos.Y + length
	for i := pos.Y; i < end; i++ {
		c.SetRune(i, pos.Y, RuneVLine, style)
	}
}

// draw a box, at position, of size, with or without a border, with or without
// being filled in and following the given theme
func (c *CCanvas) Box(pos Point2I, size Rectangle, border bool, fill bool, theme Theme) {
	TraceDF(1, "c.Box(%v,%v,%v,%v)", pos, size, border, theme)
	endx := pos.X + size.W - 1
	endy := pos.Y + size.H - 1
	// for each column
	for ix := pos.X; ix < (pos.X + size.W); ix++ {
		// for each row
		for iy := pos.Y; iy < (pos.Y + size.H); iy++ {
			if theme.Overlay {
				theme.Border = theme.Border.
					Background(c.buffer.GetBgColor(ix, iy)).
					Dim(c.buffer.GetDim(ix, iy))
				theme.Normal = theme.Normal.
					Background(c.buffer.GetBgColor(ix, iy)).
					Dim(c.buffer.GetDim(ix, iy))
			}
			switch {
			case ix == pos.X:
				// left column
				switch {
				case iy == pos.Y && border:
					// top left corner
					c.SetRune(ix, iy, theme.BorderRunes.TopLeft, theme.Border)
				case iy == endy && border:
					// bottom left corner
					c.SetRune(ix, iy, theme.BorderRunes.BottomLeft, theme.Border)
				default:
					// left border
					if border {
						c.SetRune(ix, iy, theme.BorderRunes.Left, theme.Border)
					} else if fill {
						c.SetRune(ix, iy, theme.FillRune, theme.Normal)
					}
				} // left column switch
			case ix == endx:
				// right column
				switch {
				case iy == pos.Y && border:
					// top right corner
					c.SetRune(ix, iy, theme.BorderRunes.TopRight, theme.Border)
				case iy == endy && border:
					// bottom right corner
					c.SetRune(ix, iy, theme.BorderRunes.BottomRight, theme.Border)
				default:
					// right border
					if border {
						c.SetRune(ix, iy, theme.BorderRunes.Right, theme.Border)
					} else if fill {
						c.SetRune(ix, iy, theme.FillRune, theme.Normal)
					}
				} // right column switch
			default:
				// middle columns
				switch {
				case iy == pos.Y && border:
					// top middle
					c.SetRune(ix, iy, theme.BorderRunes.Top, theme.Border)
				case iy == endy && border:
					// bottom middle
					c.SetRune(ix, iy, theme.BorderRunes.Bottom, theme.Border)
				default:
					// middle middle
					if fill {
						c.SetRune(ix, iy, theme.FillRune, theme.Normal)
					}
				} // middle columns switch
			} // draw switch
		} // for iy
	} // for ix
}

// draw a box with Sprintf-formatted text along the top-left of the box, useful
// for debugging more than anything else as the normal draw primitives are far
// more flexible
func (c *CCanvas) DebugBox(color Color, format string, argv ...interface{}) {
	text := fmt.Sprintf(format, argv...)
	bs := DefaultMonoTheme
	bs.Border = bs.Border.Foreground(color)
	c.Box(
		MakePoint2I(0, 0),
		c.size,
		true,
		false,
		bs,
	)
	c.DrawSingleLineText(MakePoint2I(1, 0), c.size.W-2, JUSTIFY_LEFT, bs.Border, false, text)
}

// fill the entire canvas according to the given theme
func (c *CCanvas) Fill(s Theme) {
	TraceF("c.fill(%v,%v)", s)
	c.Box(MakePoint2I(0, 0), c.size, false, true, s)
}

// fill the entire canvas, with or without 'dim' styling, with or without a
// border
func (c *CCanvas) FillBorder(dim bool, border bool) {
	TraceF("c.FillBorder(%v,%v): origin=%v, size=%v", dim, border, c.origin, c.size)
	s := c.theme
	if dim {
		s.Normal = s.Normal.Dim(true)
		s.Border = s.Border.Dim(true)
	}
	c.Box(
		MakePoint2I(0, 0),
		c.size,
		border,
		true,
		c.theme,
	)
}

// fill the entire canvas, with or without 'dim' styling, with plain text
// justified across the top border
func (c *CCanvas) FillBorderTitle(dim bool, title string, justify Justification) {
	TraceF("c.FillBorderTitle(%v)", dim)
	s := c.theme
	if dim {
		s.Normal = s.Normal.Dim(true)
		s.Border = s.Border.Dim(true)
	}
	c.Box(
		MakePoint2I(0, 0),
		c.size,
		true,
		true,
		c.theme,
	)
	c.DrawSingleLineText(MakePoint2I(1, 0), c.size.W-2, justify, c.theme.Normal.Dim(dim), false, title)
}
