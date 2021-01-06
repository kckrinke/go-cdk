package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWordLine(t *testing.T) {
	Convey("Word Lines with...", t, func() {
		Convey("Basic checks", func() {
			wl := NewEmptyWordLine()
			So(wl, ShouldNotBeNil)
			So(wl.Len(), ShouldEqual, 0)
			So(wl.CharacterCount(), ShouldEqual, 0)
			So(wl.WordCount(), ShouldEqual, 0)
			wl.AppendWord("word", DefaultMonoCdkStyle)
			So(wl.Len(), ShouldEqual, 1)
			So(wl.CharacterCount(), ShouldEqual, 4)
			So(wl.WordCount(), ShouldEqual, 1)
			wl.AppendWordCell(NewWordCell(" ", DefaultMonoCdkStyle))
			So(wl.Len(), ShouldEqual, 2)
			So(wl.CharacterCount(), ShouldEqual, 5)
			So(wl.WordCount(), ShouldEqual, 1)
			wl.SetLine("word\nline", DefaultMonoCdkStyle)
			So(wl.Len(), ShouldEqual, 3)
			So(wl.CharacterCount(), ShouldEqual, 9)
			So(wl.WordCount(), ShouldEqual, 2)
			w := wl.GetWord(1)
			So(w, ShouldNotBeNil)
			So(w.Len(), ShouldEqual, 1)
			w = wl.GetWord(3)
			So(w, ShouldBeNil)
			wl.RemoveWord(0)
			So(wl.Len(), ShouldEqual, 2)
			So(wl.CharacterCount(), ShouldEqual, 5)
			So(wl.WordCount(), ShouldEqual, 1)
			c := wl.GetCharacter(2)
			So(c, ShouldNotBeNil)
			So(c.Value(), ShouldEqual, 'i')
			So(c.Style().String(), ShouldEqual, DefaultMonoCdkStyle.String())
			So(c.Width(), ShouldEqual, 1)
			c = wl.GetCharacter(6)
			So(c, ShouldBeNil)
			So(wl.HasSpace(), ShouldEqual, true)
			wl.RemoveWord(0)
			So(wl.HasSpace(), ShouldEqual, false)
		})
		Convey("Empty Strings", func() {
			wl := NewWordLine("", DefaultMonoCdkStyle)
			So(wl, ShouldNotBeNil)
			So(wl.Value(), ShouldEqual, "")
			So(wl.String(), ShouldEqual, "{}")
			So(wl.Words(), ShouldHaveLength, 0)
			So(wl.Len(), ShouldEqual, 0)
		})
		Convey("One Word", func() {
			wl := NewWordLine("word", DefaultMonoCdkStyle)
			So(wl, ShouldNotBeNil)
			So(wl.Value(), ShouldEqual, "word")
			So(wl.Words(), ShouldHaveLength, 1)
			So(wl.CharacterCount(), ShouldEqual, 4)
			So(wl.WordCount(), ShouldEqual, 1)
			err := wl.AppendWordRune(0, '!', DefaultMonoCdkStyle)
			So(err, ShouldBeNil)
			So(wl.Value(), ShouldEqual, "word!")
			err = wl.AppendWordRune(1, '#', DefaultMonoCdkStyle)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "word at index 1 not found")
		})
		Convey("More Than One Word", func() {
			wl := NewWordLine("more than words", DefaultMonoCdkStyle)
			So(wl, ShouldNotBeNil)
			So(wl.Value(), ShouldEqual, "more than words")
			So(wl.Words(), ShouldHaveLength, 5)
			So(wl.CharacterCount(), ShouldEqual, 15)
		})
		// maxChars, fillerStyle
		// wrap: word, word-char, char, none
		// justify: full, center, right, left
		Convey("Making things...", func() {
			Convey("wrapping on...", func() {
				Convey("word", func() {
					wl := NewWordLine("one two three", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines := wl.Make(WRAP_WORD, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Len(), ShouldEqual, 4)
					So(lines[1].Len(), ShouldEqual, 1)
					wl = NewWordLine("one two three\nfour five six", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_WORD, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 4)
					So(lines[0].Len(), ShouldEqual, 4)
					So(lines[0].Value(), ShouldEqual, "one two ")
					So(lines[1].Value(), ShouldEqual, "three")
					So(lines[1].Len(), ShouldEqual, 1)
					So(lines[2].Len(), ShouldEqual, 3)
					So(lines[2].Value(), ShouldEqual, "four five")
					So(lines[3].Len(), ShouldEqual, 1)
					So(lines[3].Value(), ShouldEqual, "six")
					wl = NewWordLine("1234567890ABCDEF", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_WORD, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 1)
					So(lines[0].Value(), ShouldEqual, "1234567890")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 1)
					wl = NewWordLine("1234567890   DEF", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_WORD, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Value(), ShouldEqual, "1234567890")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 1)
					So(lines[1].Value(), ShouldEqual, "DEF")
					So(lines[1].CharacterCount(), ShouldEqual, 3)
					So(lines[1].Len(), ShouldEqual, 1)
				})
				Convey("char", func() {
					wl := NewWordLine("one two three", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines := wl.Make(WRAP_CHAR, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Value(), ShouldEqual, "one two th")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 5)
					So(lines[1].Value(), ShouldEqual, "ree")
					So(lines[1].CharacterCount(), ShouldEqual, 3)
					So(lines[1].Len(), ShouldEqual, 1)
					wl = NewWordLine("1234567890   DEF", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_CHAR, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Value(), ShouldEqual, "1234567890")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 1)
					So(lines[1].Value(), ShouldEqual, "DEF")
					So(lines[1].CharacterCount(), ShouldEqual, 3)
					So(lines[1].Len(), ShouldEqual, 1)
				})
				Convey("word-char", func() {
					wl := NewWordLine("one two three", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines := wl.Make(WRAP_WORD_CHAR, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Value(), ShouldEqual, "one two ")
					So(lines[0].CharacterCount(), ShouldEqual, 8)
					So(lines[0].Len(), ShouldEqual, 4)
					So(lines[1].Value(), ShouldEqual, "three")
					So(lines[1].CharacterCount(), ShouldEqual, 5)
					So(lines[1].Len(), ShouldEqual, 1)
					wl = NewWordLine("1234567890ABCDEF", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_WORD_CHAR, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Value(), ShouldEqual, "1234567890")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 1)
					So(lines[1].Value(), ShouldEqual, "ABCDEF")
					So(lines[1].CharacterCount(), ShouldEqual, 6)
					So(lines[1].Len(), ShouldEqual, 1)
				})
				Convey("none", func() {
					wl := NewWordLine("one two three", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines := wl.Make(WRAP_NONE, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 1)
					So(lines[0].Value(), ShouldEqual, "one two th")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 5)
					wl = NewWordLine("1234567890ABCDEF", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_NONE, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 1)
					So(lines[0].Value(), ShouldEqual, "1234567890")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 1)
					wl = NewWordLine("one two three\n1234567890ABCDEF", DefaultMonoCdkStyle)
					So(wl, ShouldNotBeNil)
					lines = wl.Make(WRAP_NONE, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
					So(lines, ShouldHaveLength, 2)
					So(lines[0].Value(), ShouldEqual, "one two th")
					So(lines[0].CharacterCount(), ShouldEqual, 10)
					So(lines[0].Len(), ShouldEqual, 5)
					So(lines[1].Value(), ShouldEqual, "1234567890")
					So(lines[1].CharacterCount(), ShouldEqual, 10)
					So(lines[1].Len(), ShouldEqual, 1)
				})
			})
			Convey("aligning to...", func() {
				wl := NewWordLine("one", DefaultMonoCdkStyle)
				So(wl, ShouldNotBeNil)
				lines := wl.Make(WRAP_NONE, JUSTIFY_LEFT, 10, DefaultMonoCdkStyle)
				So(lines, ShouldHaveLength, 1)
				So(lines[0].Value(), ShouldEqual, "one")
				So(lines[0].CharacterCount(), ShouldEqual, 3)
				So(lines[0].Len(), ShouldEqual, 1)
				lines = wl.Make(WRAP_NONE, JUSTIFY_RIGHT, 10, DefaultMonoCdkStyle)
				So(lines, ShouldHaveLength, 1)
				So(lines[0].Value(), ShouldEqual, "       one")
				So(lines[0].CharacterCount(), ShouldEqual, 10)
				So(lines[0].Len(), ShouldEqual, 8)
				lines = wl.Make(WRAP_NONE, JUSTIFY_CENTER, 10, DefaultMonoCdkStyle)
				So(lines, ShouldHaveLength, 1)
				So(lines[0].Value(), ShouldEqual, "    one")
				So(lines[0].CharacterCount(), ShouldEqual, 7)
				So(lines[0].Len(), ShouldEqual, 5)
				wl = NewWordLine("one two", DefaultMonoCdkStyle)
				So(wl, ShouldNotBeNil)
				lines = wl.Make(WRAP_NONE, JUSTIFY_FILL, 10, DefaultMonoCdkStyle)
				So(lines, ShouldHaveLength, 1)
				So(lines[0].Value(), ShouldEqual, "one    two")
				So(lines[0].CharacterCount(), ShouldEqual, 10)
				So(lines[0].Len(), ShouldEqual, 3)
				wl = NewWordLine("1 2 3 4", DefaultMonoCdkStyle)
				So(wl, ShouldNotBeNil)
				lines = wl.Make(WRAP_NONE, JUSTIFY_FILL, 10, DefaultMonoCdkStyle)
				So(lines, ShouldHaveLength, 1)
				So(lines[0].Value(), ShouldEqual, "1  2  3  4")
				So(lines[0].CharacterCount(), ShouldEqual, 10)
				So(lines[0].Len(), ShouldEqual, 7)
			})
		})
	})
}
