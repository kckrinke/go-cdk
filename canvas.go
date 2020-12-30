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

type Canvas struct {
	buffer *CanvasBuffer
	origin Point2I
	size   Rectangle
	theme  Theme
	fill   rune
}

func NewCanvas(origin Point2I, size Rectangle, theme Theme) *Canvas {
	return &Canvas{
		buffer: NewCanvasBuffer(size, theme.Normal),
		origin: origin,
		size:   size,
		theme:  theme,
		fill:   ' ',
	}
}

func (c Canvas) String() string {
	return fmt.Sprintf(
		"{Origin=%s,Size=%s,Theme=%s,Fill=%v,Buffer=%v}",
		c.origin,
		c.size,
		c.theme,
		c.fill,
		c.buffer.String(),
	)
}

/* Canvas Methods */

func (c *Canvas) Resize(size Rectangle) {
	c.buffer.Resize(size)
	c.size = size
}

func (c *Canvas) GetContent(x, y int) (mainc rune, style Style, width int) {
	return c.buffer.GetContent(x, y)
}

func (c *Canvas) SetContent(x, y int, char string, s Style) error {
	r, _ := utf8.DecodeRune([]byte(char))
	return c.buffer.SetContent(x, y, r, s)
}

func (c *Canvas) SetRune(x, y int, r rune, s Style) error {
	return c.buffer.SetContent(x, y, r, s)
}

func (c *Canvas) SetOrigin(origin Point2I) {
	c.origin = origin
}

func (c *Canvas) GetOrigin() Point2I {
	return c.origin
}

func (c *Canvas) SetSize(size Rectangle) {
	c.size = size
}

func (c *Canvas) GetSize() Rectangle {
	return c.size
}

func (c *Canvas) SetTheme(style Theme) {
	c.theme = style
}

func (c *Canvas) GetTheme() Theme {
	return c.theme
}

func (c *Canvas) SetFill(fill rune) {
	c.fill = fill
}

func (c *Canvas) GetFill() rune {
	return c.fill
}

func (c *Canvas) Composite(v *Canvas) error {
	for x := 0; x < v.buffer.size.W; x++ {
		for y := 0; y < v.buffer.size.H; y++ {
			cell := v.buffer.Cell(x, y)
			if cell != nil {
				if cell.dirty {
					if err := c.buffer.SetContent(
						v.origin.X+x,
						v.origin.Y+y,
						cell.Value(),
						cell.Style(),
					); err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("cell is nil x=%v,y=%v", x, y)
			}
		}
	}
	return nil
}

func (c *Canvas) Render(screen Screen) error {
	for x := 0; x < c.size.W; x++ {
		for y := 0; y < c.size.H; y++ {
			cell := c.buffer.Cell(x, y)
			screen.SetContent(x, y, cell.Value(), nil, cell.style)
		}
	}
	return nil
}

/* Draw Primitives */

// Write text to the canvas buffer
// origin is the top-left coordinate for the text area being rendered
// alignment is based on origin.X boxed by maxChars or canvas size.W
func (c *Canvas) DrawText(pos Point2I, size Rectangle, justify Justification, singleLineMode bool, wrap WrapMode, style Style, markup bool, text string) {
	var tb TextBuffer
	if markup {
		m, err := NewMarkup(text, style)
		if err != nil {
			Fataldf(1, "failed to parse markup: %v", err)
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
		Errorf("composite error: %v", err)
	}
}

func (c *Canvas) DrawSingleLineText(pos Point2I, maxChars int, justify Justification, style Style, markup bool, text string) {
	c.DrawText(pos, MakeRectangle(maxChars, 1), justify, true, WRAP_NONE, style, markup, text)
}

func (c *Canvas) DrawLine(pos Point2I, length int, orient Orientation, style Style) {
	Tracef("c.Line(%v,%v,%v,%v)", pos, length, orient, style)
	switch orient {
	case ORIENTATION_HORIZONTAL:
		c.DrawHorizontalLine(pos, length, style)
	case ORIENTATION_VERTICAL:
		c.DrawVerticalLine(pos, length, style)
	}
}

func (c *Canvas) DrawHorizontalLine(pos Point2I, length int, style Style) {
	length = utils.ClampI(length, pos.X, c.size.W-pos.X)
	end := pos.X + length
	for i := pos.X; i < end; i++ {
		c.SetRune(i, pos.Y, RuneHLine, style)
	}
}

func (c *Canvas) DrawVerticalLine(pos Point2I, length int, style Style) {
	length = utils.ClampI(length, pos.Y, c.size.H-pos.Y)
	end := pos.Y + length
	for i := pos.Y; i < end; i++ {
		c.SetRune(i, pos.Y, RuneVLine, style)
	}
}

func (c *Canvas) Box(pos Point2I, size Rectangle, border bool, fill bool, theme Theme) {
	Tracedf(1, "c.Box(%v,%v,%v,%v)", pos, size, border, theme)
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

/* Draw Features */

func (c *Canvas) Fill(s Theme) {
	Tracef("c.fill(%v,%v)", s)
	c.Box(MakePoint2I(0, 0), c.size, false, true, s)
}

func (c *Canvas) FillBorder(dim bool, border bool) {
	Tracef("c.FillBorder(%v,%v): origin=%v, size=%v", dim, border, c.origin, c.size)
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

func (c *Canvas) FillBorderTitle(dim bool, title string, justify Justification) {
	Tracef("c.FillBorderTitle(%v)", dim)
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
	label := fmt.Sprintf(" %v ", title)
	pos := MakePoint2I((c.size.W / 2) - (len(label) / 2), 0)
	c.DrawSingleLineText(pos, len(label), JUSTIFY_CENTER, c.theme.Normal.Dim(dim), false, label)
}
