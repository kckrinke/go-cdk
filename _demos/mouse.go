//+build ignore

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

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/encoding"

	"github.com/mattn/go-runewidth"
)

var defStyle cdk.Style

func emitStr(s cdk.Screen, x, y int, style cdk.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawBox(s cdk.Screen, x1, y1, x2, y2 int, style cdk.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, cdk.RuneHLine, nil, style)
		s.SetContent(col, y2, cdk.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, cdk.RuneVLine, nil, style)
		s.SetContent(x2, row, cdk.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, cdk.RuneULCorner, nil, style)
		s.SetContent(x2, y1, cdk.RuneURCorner, nil, style)
		s.SetContent(x1, y2, cdk.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, cdk.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

func drawSelect(s cdk.Screen, x1, y1, x2, y2 int, sel bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			mainc, combc, style, width := s.GetContent(col, row)
			if style == cdk.StyleDefault {
				style = defStyle
			}
			style = style.Reverse(sel)
			s.SetContent(col, row, mainc, combc, style)
			col += width - 1
		}
	}
}

// This program just shows simple mouse and keyboard events.  Press ESC twice to
// exit.
func main() {

	encoding.Register()
	if err := cdk.MainInit(); err != nil {
		cdk.Fatalf("MainInit() error: %v", err)
	}
	s := cdk.MainScreen()

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	pastefmt := "Paste: [%d] %s"
	white := cdk.StyleDefault.Foreground(cdk.ColorWhite).Background(cdk.ColorRed)

	mx, my := -1, -1
	ox, oy := -1, -1
	bx, by := -1, -1
	w, h := s.Size()
	lchar := '*'
	bstr := ""
	lks := ""
	pstr := ""
	ecnt := 0
	pasting := false
	st := cdk.StyleDefault.Background(cdk.ColorRed)
	up := cdk.StyleDefault.
		Background(cdk.ColorBlue).
		Foreground(cdk.ColorBlack)

	cdk.SetMainPasteHandler(func(ev *cdk.EventPaste) error {
		pasting = ev.Start()
		if pasting {
			pstr = ""
		}
		return nil
	})
	cdk.SetMainResizeHandler(func(ev *cdk.EventResize) error {
		s.Sync()
		s.SetContent(w-1, h-1, 'R', nil, st)
		return nil
	})
	cdk.SetMainKeyHandler(func(ev *cdk.EventKey) error {
		s.SetContent(w-2, h-2, ev.Rune(), nil, st)
		if pasting {
			s.SetContent(w-1, h-1, 'P', nil, st)
			if ev.Key() == cdk.KeyRune {
				pstr = pstr + string(ev.Rune())
			} else {
				pstr = pstr + "\ufffd" // replacement for now
			}
			lks = ""
			return nil
		}
		pstr = ""
		s.SetContent(w-1, h-1, 'K', nil, st)
		if ev.Key() == cdk.KeyEscape {
			ecnt++
			if ecnt > 1 {
				cdk.MainQuit()
			}
		} else if ev.Key() == cdk.KeyCtrlL {
			s.Sync()
		} else {
			ecnt = 0
			if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		}
		lks = ev.Name()
		return nil
	})
	cdk.SetMainMouseHandler(func(ev *cdk.EventMouse) error {
		x, y := ev.Position()
		button := ev.Buttons()
		for i := uint(0); i < 8; i++ {
			if int(button)&(1<<i) != 0 {
				bstr += fmt.Sprintf(" Button%d", i+1)
			}
		}
		if button&cdk.WheelUp != 0 {
			bstr += " WheelUp"
		}
		if button&cdk.WheelDown != 0 {
			bstr += " WheelDown"
		}
		if button&cdk.WheelLeft != 0 {
			bstr += " WheelLeft"
		}
		if button&cdk.WheelRight != 0 {
			bstr += " WheelRight"
		}
		// Only buttons, not wheel events
		button &= cdk.ButtonMask(0xff)
		ch := '*'

		if button != cdk.ButtonNone && ox < 0 {
			ox, oy = x, y
		}
		switch ev.Buttons() {
		case cdk.ButtonNone:
			if ox >= 0 {
				bg := cdk.Color((lchar-'0')*2) | cdk.ColorValid
				drawBox(s, ox, oy, x, y,
					up.Background(bg),
					lchar)
				ox, oy = -1, -1
				bx, by = -1, -1
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
			bx, by = x, y
		}
		lchar = ch
		s.SetContent(w-1, h-1, 'M', nil, st)
		mx, my = x, y

		return nil
	})
	cdk.SetMainDefaultHandler(func(ev cdk.Event) error {
		cdk.Infof("default event: %v\n", ev)
		s.SetContent(w-1, h-1, 'X', nil, st)
		return nil
	})
	cdk.SetMainPreLoopHandler(func() error {
		drawBox(s, 1, 1, 42, 7, white, ' ')
		emitStr(s, 2, 2, white, "Press ESC twice to exit, C to clear.")
		emitStr(s, 2, 3, white, fmt.Sprintf(posfmt, mx, my))
		emitStr(s, 2, 4, white, fmt.Sprintf(btnfmt, bstr))
		emitStr(s, 2, 5, white, fmt.Sprintf(keyfmt, lks))

		ps := pstr
		if len(ps) > 26 {
			ps = "..." + ps[len(ps)-24:]
		}
		emitStr(s, 2, 6, white, fmt.Sprintf(pastefmt, len(pstr), ps))

		s.Show()
		bstr = ""

		w, h = s.Size()

		// always clear any old selection box
		if ox >= 0 && oy >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, false)
		}
		return nil
	})
	cdk.SetMainPostLoopHandler(func() error {
		if ox >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, true)
		}
		return nil
	})
	if err := cdk.MainLoop(); err != nil {
		cdk.Errorf("exiting main loop, error - %v\n", err)
	}
	cdk.MainQuit()
}
