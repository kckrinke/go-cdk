// +build ignore

// Copyright 2020 The cdk Authors
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

package main

import (
	"fmt"
	"os"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/encoding"

	"github.com/mattn/go-runewidth"
)

func emitStr(s cdk.Display, x, y int, style cdk.Style, str string) {
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

func displayHelloWorld(s cdk.Display) {
	w, h := s.Size()
	s.Clear()
	emitStr(s, w/2-7, h/2, cdk.StyleDefault, "Hello, World!")
	emitStr(s, w/2-9, h/2+1, cdk.StyleDefault, "Press ESC to exit.")
	s.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	encoding.Register()

	s, e := cdk.NewDisplay()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := cdk.StyleDefault.
		Background(cdk.ColorBlack).
		Foreground(cdk.ColorWhite)
	s.SetStyle(defStyle)

	displayHelloWorld(s)

	for {
		switch ev := s.PollEvent().(type) {
		case *cdk.EventResize:
			s.Sync()
			displayHelloWorld(s)
		case *cdk.EventKey:
			if ev.Key() == cdk.KeyEscape {
				s.Close()
				os.Exit(0)
			}
		}
	}
}
