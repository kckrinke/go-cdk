package cdk

import (
	"fmt"
	"unicode/utf8"
)

type View struct {
	buffer *ViewBuffer
	origin Point2I
	size   Rectangle
	theme  Theme
	fill   rune
}

func NewCView(origin Point2I, size Rectangle, theme Theme) *View {
	return &View{
		buffer: NewViewBuffer(size, theme.Normal),
		origin: origin,
		size:   size,
		theme:  theme,
		fill:   ' ',
	}
}

/* View Methods */

func (v *View) Resize(size Rectangle) {
	v.buffer.Resize(size)
	v.size = size
}

func (view *View) SetContent(x, y int, char string, s Style) error {
	r, _ := utf8.DecodeRune([]byte(char))
	return view.buffer.SetContent(x, y, r, s)
}

func (view *View) SetRune(x, y int, r rune, s Style) error {
	return view.buffer.SetContent(x, y, r, s)
}

func (view *View) SetOrigin(origin Point2I) {
	view.origin = origin
}

func (view *View) GetOrigin() Point2I {
	return view.origin
}

func (view *View) SetSize(size Rectangle) {
	view.size = size
}

func (view *View) GetSize() Rectangle {
	return view.size
}

func (view *View) SetTheme(style Theme) {
	view.theme = style
}

func (view *View) GetTheme() Theme {
	return view.theme
}

func (view *View) SetFill(fill rune) {
	view.fill = fill
}

func (view *View) GetFill() rune {
	return view.fill
}

func (view *View) Composite(v *View) error {
	for x := 0; x < v.buffer.size.W; x++ {
		for y := 0; y < v.buffer.size.H; y++ {
			cell := v.buffer.Cell(x, y)
			if cell != nil {
				if cell.dirty {
					if err := view.buffer.SetContent(
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

/* Draw Primitives */

// Write text to the view buffer
// origin is the top-left coordinate for the text area being rendered
// alignment is based on origin.X boxed by maxChars or view size.W
func (view *View) DrawText(pos Point2I, size Rectangle, justify Justification, singleLineMode bool, wrap WrapMode, style Style, markup bool, text string) {
	var tb TextBuffer
	if markup {
		m := NewMarkup(text, style)
		tb = m.TextBuffer()
	} else {
		tb = NewTextBuffer(text, style)
	}
	if size.W == -1 || size.W >= view.size.W {
		size.W = view.size.W
	}
	v := NewCView(pos, size, view.theme)
	tb.Draw(v, singleLineMode, wrap, justify, ALIGN_TOP)
	if err := view.Composite(v); err != nil {
		Errorf("composite error: %v", err)
	}
}

func (view *View) DrawSingleLineText(pos Point2I, maxChars int, justify Justification, style Style, markup bool, text string) {
	view.DrawText(pos, Rectangle{H: 1, W: maxChars}, justify, true, WRAP_NONE, style, markup, text)
}

func (view *View) DrawLine(pos Point2I, length int, orient Orientation, style Style) {
	Tracef("view.Line(%v,%v,%v,%v)", pos, length, orient, style)
	switch orient {
	case ORIENTATION_HORIZONTAL:
		view.DrawHorizontalLine(pos, length, style)
	case ORIENTATION_VERTICAL:
		view.DrawVerticalLine(pos, length, style)
	}
}

func (view *View) DrawHorizontalLine(pos Point2I, length int, style Style) {
	length = ClampI(length, pos.X, view.size.W-pos.X)
	end := pos.X + length
	for i := pos.X; i < end; i++ {
		view.SetRune(i, pos.Y, RuneHLine, style)
	}
}

func (view *View) DrawVerticalLine(pos Point2I, length int, style Style) {
	length = ClampI(length, pos.Y, view.size.H-pos.Y)
	end := pos.Y + length
	for i := pos.Y; i < end; i++ {
		view.SetRune(i, pos.Y, RuneVLine, style)
	}
}

func (view *View) Box(pos Point2I, size Rectangle, border bool, fill bool, theme Theme) {
	Tracef("view.Box(%v,%v,%v,%v)", pos, size, border, theme)
	endx := pos.X + size.W - 1
	endy := pos.Y + size.H - 1
	// for each column
	for ix := pos.X; ix < (pos.X + size.W); ix++ {
		// for each row
		for iy := pos.Y; iy < (pos.Y + size.H); iy++ {
			if theme.Overlay {
				theme.Border = theme.Border.
					Background(view.buffer.GetBgColor(ix, iy)).
					Dim(view.buffer.GetDim(ix, iy))
			}
			switch {
			case ix == pos.X:
				// left column
				switch {
				case iy == pos.Y && border:
					// top left corner
					view.SetRune(ix, iy, theme.BorderRunes.TopLeft, theme.Border)
				case iy == endy && border:
					// bottom left corner
					view.SetRune(ix, iy, theme.BorderRunes.BottomLeft, theme.Border)
				default:
					// left border
					if border {
						view.SetRune(ix, iy, theme.BorderRunes.Left, theme.Border)
					} else if fill {
						view.SetRune(ix, iy, theme.FillRune, theme.Normal)
					}
				} // left column switch
			case ix == endx:
				// right column
				switch {
				case iy == pos.Y && border:
					// top right corner
					view.SetRune(ix, iy, theme.BorderRunes.TopRight, theme.Border)
				case iy == endy && border:
					// bottom right corner
					view.SetRune(ix, iy, theme.BorderRunes.BottomRight, theme.Border)
				default:
					// right border
					if border {
						view.SetRune(ix, iy, theme.BorderRunes.Right, theme.Border)
					} else if fill {
						view.SetRune(ix, iy, theme.FillRune, theme.Normal)
					}
				} // right column switch
			default:
				// middle columns
				switch {
				case iy == pos.Y && border:
					// top middle
					view.SetRune(ix, iy, theme.BorderRunes.Top, theme.Border)
				case iy == endy && border:
					// bottom middle
					view.SetRune(ix, iy, theme.BorderRunes.Bottom, theme.Border)
				default:
					// middle middle
					if fill {
						view.SetRune(ix, iy, theme.FillRune, theme.Normal)
					}
				} // middle columns switch
			} // draw switch
		} // for iy
	} // for ix
}

/* Draw Features */

func (view *View) Fill(s Theme) {
	Tracef("view.fill(%v,%v)", s)
	view.Box(Point2I{0, 0}, view.size, false, true, s)
}

func (view *View) FillBorder(dim bool, border bool) {
	Tracef("view.FillBorder(%v,%v): origin=%v, size=%v", dim, border, view.origin, view.size)
	s := view.theme
	if dim {
		s.Normal = s.Normal.Dim(true)
		s.Border = s.Border.Dim(true)
	}
	view.Box(
		Point2I{0, 0},
		view.size,
		border,
		true,
		view.theme,
	)
}

func (view *View) FillBorderTitle(dim bool, title string, justify Justification) {
	Tracef("view.FillBorderTitle(%v)", dim)
	s := view.theme
	if dim {
		s.Normal = s.Normal.Dim(true)
		s.Border = s.Border.Dim(true)
	}
	view.Box(
		Point2I{0, 0},
		view.size,
		true,
		true,
		view.theme,
	)
	label := fmt.Sprintf(" %v ", title)
	pos := Point2I{(view.size.W / 2) - (len(label) / 2), 0}
	view.DrawSingleLineText(pos, len(label), JUSTIFY_CENTER, view.theme.Normal.Dim(dim), false, label)
}
