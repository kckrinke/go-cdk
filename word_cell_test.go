package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWordCell(t *testing.T) {
	Convey("Word Cells with...", t, func() {
		Convey("Empty Strings", func() {
			wc := NewWordCell("", DefaultMonoCdkStyle)
			So(wc, ShouldNotBeNil)
			So(wc.String(), ShouldEqual, "")
			So(wc.Characters(), ShouldHaveLength, 0)
			So(wc.Len(), ShouldEqual, 0)
		})
		Convey("One Word", func() {
			wc := NewWordCell("word", DefaultMonoCdkStyle)
			So(wc, ShouldNotBeNil)
			So(wc.Value(), ShouldEqual, "word")
			So(wc.Characters(), ShouldHaveLength, 4)
			So(wc.Len(), ShouldEqual, 4)
		})
		Convey("More Than One Word", func() {
			wc := NewWordCell("another word", DefaultMonoCdkStyle)
			So(wc, ShouldNotBeNil)
			So(wc.CompactLen(), ShouldEqual, 12)
		})
		Convey("Basic checks", func() {
			wc := NewEmptyWordCell()
			So(wc, ShouldNotBeNil)
			So(wc.Len(), ShouldEqual, 0)
			So(wc.IsSpace(), ShouldEqual, true)
			So(wc.HasSpace(), ShouldEqual, false)
			So(wc.GetCharacter(0), ShouldBeNil)
			wc.AppendRune(' ', DefaultMonoCdkStyle)
			So(wc.Len(), ShouldEqual, 1)
			c := wc.GetCharacter(0)
			So(c, ShouldNotBeNil)
			So(c.Value(), ShouldEqual, ' ')
			So(c.Style().String(), ShouldEqual, DefaultMonoCdkStyle.String())
			So(c.Width(), ShouldEqual, 1)
			So(wc.HasSpace(), ShouldEqual, true)
			So(wc.CompactLen(), ShouldEqual, 1)
			wc.AppendRune(' ', DefaultMonoCdkStyle)
			So(wc.CompactLen(), ShouldEqual, 1)
			So(wc.Len(), ShouldEqual, 2)
		})
	})
}
