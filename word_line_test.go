package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWordLine(t *testing.T) {
	Convey("Word Lines with...", t, func() {
		Convey("Empty Strings", func() {
			wl := NewWordLine("", DefaultMonoCdkStyle)
			So(wl, ShouldNotBeNil)
			So(wl.String(), ShouldEqual, "")
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
		})
		Convey("More Than One Word", func() {
			wl := NewWordLine("more than words", DefaultMonoCdkStyle)
			So(wl, ShouldNotBeNil)
			So(wl.Value(), ShouldEqual, "more   than   words")
			So(wl.Words(), ShouldHaveLength, 5)
			So(wl.CharacterCount(), ShouldEqual, 15)
		})
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
		Convey("Make checks", func() {
			// wrap: word, word-char, char, none
			// justify: full, center, right, left
		})
	})
}
