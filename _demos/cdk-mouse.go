// +build cdkmouse

// Copyright 2015 The TCell Authors
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

// mouse displays a text box and tests mouse interaction.  As you click
// and drag, boxes are displayed on screen.  Other events are reported in
// the box.  Press ESC twice to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/kckrinke/go-cdk"

	"github.com/mattn/go-runewidth"
)

var defStyle cdk.Style

func emitStr(s *cdk.Canvas, x, y int, style cdk.Style, str string) {
	for _, c := range str {
		w := runewidth.RuneWidth(c)
		if w == 0 {
			c = ' '
			w = 1
		}
		s.SetRune(x, y, c, style)
		x += w
	}
}

func drawBox(s *cdk.Canvas, x1, y1, x2, y2 int, style cdk.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetRune(col, y1, cdk.RuneHLine, style)
		s.SetRune(col, y2, cdk.RuneHLine, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetRune(x1, row, cdk.RuneVLine, style)
		s.SetRune(x2, row, cdk.RuneVLine, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetRune(x1, y1, cdk.RuneULCorner, style)
		s.SetRune(x2, y1, cdk.RuneURCorner, style)
		s.SetRune(x1, y2, cdk.RuneLLCorner, style)
		s.SetRune(x2, y2, cdk.RuneLRCorner, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetRune(col, row, r, style)
		}
	}
}

func drawSelect(s *cdk.Canvas, x1, y1, x2, y2 int, sel bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			mainc, style, width := s.GetContent(col, row)
			if style == cdk.StyleDefault {
				style = defStyle
			}
			style = style.Reverse(sel)
			s.SetRune(col, row, mainc, style)
			col += width - 1
		}
	}
}

type DrawBoxArgs struct {
	x1, y1, x2, y2 int
	style          cdk.Style
	r              rune
}

type CdkMouseWindow struct {
	cdk.CWindow

	mx, my          int
	ox, oy          int
	bx, by          int
	w, h            int
	lchar           rune
	bstr, lks, pstr string
	ecnt            int
	pasting         bool
	st, up          cdk.Style
	mark            rune
	lkp             rune
	want_clear      bool

	drawBoxQueue []DrawBoxArgs
}

func (w *CdkMouseWindow) Init() (already bool) {
	if w.CWindow.Init() {
		return true
	}
	w.mx, w.my = -1, -1
	w.ox, w.oy = -1, -1
	w.bx, w.by = -1, -1
	w.w, w.h = 0, 0
	w.lchar = '*'
	w.bstr = ""
	w.lks = ""
	w.pstr = ""
	w.ecnt = 0
	w.pasting = false
	w.st = cdk.StyleDefault.Background(cdk.ColorRed)
	w.up = cdk.StyleDefault.
		Background(cdk.ColorBlue).
		Foreground(cdk.ColorBlack)
	w.mark = 0
	w.lkp = 0
	w.want_clear = false
	w.drawBoxQueue = []DrawBoxArgs{}
	return false
}

func (w *CdkMouseWindow) Draw(s *cdk.Canvas) cdk.EventFlag {
	w.LogInfo("Draw: %s", s)

	size := s.GetSize()
	w.w, w.h = size.W, size.H

	white := cdk.StyleDefault.
		Foreground(cdk.ColorWhite).Background(cdk.ColorRed)

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	pastefmt := "Paste: [%d] %s"

	drawBox(s, 1, 1, 42, 7, white, ' ')
	emitStr(s, 2, 2, white, "Press ESC twice to exit, C to clear.")
	emitStr(s, 2, 3, white, fmt.Sprintf(posfmt, w.mx, w.my))
	emitStr(s, 2, 4, white, fmt.Sprintf(btnfmt, w.bstr))
	emitStr(s, 2, 5, white, fmt.Sprintf(keyfmt, w.lks))

	ps := w.pstr
	if len(ps) > 26 {
		ps = "..." + ps[len(ps)-24:]
	}
	emitStr(s, 2, 6, white, fmt.Sprintf(pastefmt, len(w.pstr), ps))

	w.bstr = ""

	// always clear any old selection box
	if w.ox >= 0 && w.oy >= 0 && w.bx >= 0 {
		drawSelect(s, w.ox, w.oy, w.bx, w.by, false)
	}

	for i := 0; i < len(w.drawBoxQueue); i++ {
		a := w.drawBoxQueue[i]
		drawBox(s, a.x1, a.y1, a.x2, a.y2, a.style, a.r)
	}
	w.drawBoxQueue = []DrawBoxArgs{}

	s.SetRune(w.w-1, w.h-1, w.mark, w.st)

	if w.ox >= 0 && w.bx >= 0 {
		drawSelect(s, w.ox, w.oy, w.bx, w.by, true)
	}
	return cdk.EVENT_STOP
}

func (w *CdkMouseWindow) ProcessEvent(evt cdk.Event) cdk.EventFlag {
	w.LogInfo("ProcessEvent: %v", evt)
	switch ev := evt.(type) {
	case *cdk.EventResize:
		w.mark = 'R'
	case *cdk.EventKey:
		w.lkp = ev.Rune()
		if w.pasting {
			w.mark = 'P'
			if ev.Key() == cdk.KeyRune {
				w.pstr = w.pstr + string(ev.Rune())
			} else {
				w.pstr = w.pstr + "\ufffd" // replacement for now
			}
			w.lks = ""
			return cdk.EVENT_STOP
		}
		w.pstr = ""
		w.mark = 'K'
		if ev.Key() == cdk.KeyEscape {
			w.ecnt++
			if w.ecnt > 1 {
				w.GetDisplay().RequestQuit()
				return cdk.EVENT_STOP
			}
		} else if ev.Key() == cdk.KeyCtrlL {
			w.GetDisplay().RequestSync()
		} else {
			w.ecnt = 0
			if ev.Rune() == 'C' || ev.Rune() == 'c' {
				w.want_clear = true
				return cdk.EVENT_STOP
			}
		}
		w.lks = ev.Name()
	case *cdk.EventPaste:
		w.pasting = ev.Start()
		if w.pasting {
			w.pstr = ""
		}
	case *cdk.EventMouse:
		x, y := ev.Position()
		button := ev.Buttons()
		for i := uint(0); i < 8; i++ {
			if int(button)&(1<<i) != 0 {
				w.bstr += fmt.Sprintf(" Button%d", i+1)
			}
		}
		if button&cdk.WheelUp != 0 {
			w.bstr += " WheelUp"
		}
		if button&cdk.WheelDown != 0 {
			w.bstr += " WheelDown"
		}
		if button&cdk.WheelLeft != 0 {
			w.bstr += " WheelLeft"
		}
		if button&cdk.WheelRight != 0 {
			w.bstr += " WheelRight"
		}
		// Only buttons, not wheel events
		button &= cdk.ButtonMask(0xff)
		ch := '*'

		if button != cdk.ButtonNone && w.ox < 0 {
			w.ox, w.oy = x, y
		}
		switch ev.Buttons() {
		case cdk.ButtonNone:
			if w.ox >= 0 {
				bg := cdk.Color((w.lchar-'0')*2) | cdk.ColorValid
				w.drawBoxQueue = append(
					w.drawBoxQueue,
					DrawBoxArgs{
						w.ox,
						w.oy,
						x,
						y,
						w.up.Background(bg),
						w.lchar,
					},
				)
				w.ox, w.oy = -1, -1
				w.bx, w.by = -1, -1
			}
		case cdk.Button1:
			ch = '1'
		case cdk.Button2:
			ch = '2'
		case cdk.Button3:
			ch = '3'
		case cdk.Button4:
			ch = '4'
		case cdk.Button5:
			ch = '5'
		case cdk.Button6:
			ch = '6'
		case cdk.Button7:
			ch = '7'
		case cdk.Button8:
			ch = '8'
		default:
			ch = '*'

		}
		if button != cdk.ButtonNone {
			w.bx, w.by = x, y
		}
		w.lchar = ch
		w.mark = 'M'
		w.mx, w.my = x, y
	default:
		w.mark = 'X'
	}
	return cdk.EVENT_STOP
}

func main() {
	app := cdk.NewApp(
		"cdk-mouse",
		"The tcell mouse demo as a formal CDK Application",
		"0.0.1",
		"mouse",
		"CDK Mouse",
		"/dev/tty",
		func(d cdk.Display) error {
			cdk.Debugf("cdk-mouse initFn hit")
			d.CaptureCtrlC()
			w := &CdkMouseWindow{}
			w.Init()
			d.SetActiveWindow(w)
			return nil
		},
	)
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
