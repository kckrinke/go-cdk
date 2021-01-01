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

// unicode just displays a Unicode test on your screen.
// Press ESC to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/encoding"
	runewidth "github.com/mattn/go-runewidth"
)

var row = 0
var style = cdk.StyleDefault

func putln(s cdk.Display, str string) {

	puts(s, style, 1, row, str)
	row++
}

func puts(s cdk.Display, style cdk.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

func main() {

	s, e := cdk.NewDisplay()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	plain := cdk.StyleDefault
	bold := style.Bold(true)

	s.SetStyle(cdk.StyleDefault.
		Foreground(cdk.ColorBlack).
		Background(cdk.ColorWhite))
	s.Clear()

	quit := make(chan struct{})

	style = bold
	putln(s, "Press ESC to Exit")
	putln(s, "Character set: "+s.CharacterSet())
	style = plain

	putln(s, "English:   October")
	putln(s, "Icelandic: október")
	putln(s, "Arabic:    أكتوبر")
	putln(s, "Russian:   октября")
	putln(s, "Greek:     Οκτωβρίου")
	putln(s, "Chinese:   十月 (note, two double wide characters)")
	putln(s, "Combining: A\u030a (should look like Angstrom)")
	putln(s, "Emoticon:  \U0001f618 (blowing a kiss)")
	putln(s, "Airplane:  \u2708 (fly away)")
	putln(s, "Command:   \u2318 (mac clover key)")
	putln(s, "Enclose:   !\u20e3 (should be enclosed exclamation)")
	putln(s, "ZWJ:       \U0001f9db\u200d\u2640 (female vampire)")
	putln(s, "ZWJ:       \U0001f9db\u200d\u2642 (male vampire)")
	putln(s, "Family:    \U0001f469\u200d\U0001f467\u200d\U0001f467 (woman girl girl)\n")
	putln(s, "Region:    \U0001f1fa\U0001f1f8 (USA! USA!)\n")
	putln(s, "")
	putln(s, "Box:")
	putln(s, string([]rune{
		cdk.RuneULCorner,
		cdk.RuneHLine,
		cdk.RuneTTee,
		cdk.RuneHLine,
		cdk.RuneURCorner,
	}))
	putln(s, string([]rune{
		cdk.RuneVLine,
		cdk.RuneBullet,
		cdk.RuneVLine,
		cdk.RuneLantern,
		cdk.RuneVLine,
	})+"  (bullet, lantern/section)")
	putln(s, string([]rune{
		cdk.RuneLTee,
		cdk.RuneHLine,
		cdk.RunePlus,
		cdk.RuneHLine,
		cdk.RuneRTee,
	}))
	putln(s, string([]rune{
		cdk.RuneVLine,
		cdk.RuneDiamond,
		cdk.RuneVLine,
		cdk.RuneUArrow,
		cdk.RuneVLine,
	})+"  (diamond, up arrow)")
	putln(s, string([]rune{
		cdk.RuneLLCorner,
		cdk.RuneHLine,
		cdk.RuneBTee,
		cdk.RuneHLine,
		cdk.RuneLRCorner,
	}))

	s.Show()
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

	<-quit

	s.Close()
}
