// +build ignore

// Copyright 2019 The cdk Authors
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

// colors just displays a single centered rectangle that should pulse
// through available colors.  It uses the RGB color cube, bumping at
// predefined larger intervals (values of about 8) in order that the
// changes happen quickly enough to be appreciable.
//
// Press ESC to exit the program.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kckrinke/go-cdk"
)

var red = int32(rand.Int() % 256)
var grn = int32(rand.Int() % 256)
var blu = int32(rand.Int() % 256)
var inc = int32(8) // rate of color change
var redi = int32(inc)
var grni = int32(inc)
var blui = int32(inc)

func makebox(s cdk.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	glyphs := []rune{'@', '#', '&', '*', '=', '%', 'Z', 'A'}

	lh := h / 2
	lw := w / 2
	lx := w / 4
	ly := h / 4
	st := cdk.StyleDefault
	gl := ' '

	if s.Colors() == 0 {
		st = st.Reverse(rand.Int()%2 == 0)
		gl = glyphs[rand.Int()%len(glyphs)]
	} else {

		red += redi
		if (red >= 256) || (red < 0) {
			redi = -redi
			red += redi
		}
		grn += grni
		if (grn >= 256) || (grn < 0) {
			grni = -grni
			grn += grni
		}
		blu += blui
		if (blu >= 256) || (blu < 0) {
			blui = -blui
			blu += blui

		}
		st = st.Background(cdk.NewRGBColor(red, grn, blu))
	}
	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetCell(lx+col, ly+row, st, gl)
		}
	}
	s.Show()
}

func flipcoin() bool {
	if rand.Int()&1 == 0 {
		return false
	}
	return true
}

func main() {

	rand.Seed(time.Now().UnixNano())
	cdk.SetEncodingFallback(cdk.EncodingFallbackASCII)
	s, e := cdk.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(cdk.StyleDefault.
		Foreground(cdk.ColorBlack).
		Background(cdk.ColorWhite))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *cdk.EventKey:
				switch ev.Key() {
				case cdk.KeyEscape, cdk.KeyEnter:
					close(quit)
					return
				case cdk.KeyCtrlL:
					s.Sync()
				}
			case *cdk.EventResize:
				s.Sync()
			}
		}
	}()

	cnt := 0
	dur := time.Duration(0)
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
		start := time.Now()
		makebox(s)
		dur += time.Now().Sub(start)
		cnt++
		if cnt%(256/int(inc)) == 0 {
			if flipcoin() {
				redi = -redi
			}
			if flipcoin() {
				grni = -grni
			}
			if flipcoin() {
				blui = -blui
			}
		}
	}

	s.Close()
	fmt.Printf("Finished %d boxes in %s\n", cnt, dur)
	fmt.Printf("Average is %0.3f ms / box\n", (float64(dur)/float64(cnt))/1000000.0)
}
