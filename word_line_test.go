package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWordLine(t *testing.T) {
	Convey("Empty Strings", t, func() {
		wl := NewWordLine("", DefaultMonoCdkStyle)
		So(wl, ShouldNotBeNil)
		So(wl.String(), ShouldEqual, "")
		So(wl.Words(), ShouldHaveLength, 0)
		So(wl.LetterCount(true), ShouldEqual, 0)
	})
	Convey("One Word", t, func() {
		wl := NewWordLine("word", DefaultMonoCdkStyle)
		So(wl, ShouldNotBeNil)
		So(wl.Value(), ShouldEqual, "word")
		So(wl.Words(), ShouldHaveLength, 1)
		So(wl.LetterCount(true), ShouldEqual, 4)
	})
	Convey("More Than One Word", t, func() {
		wl := NewWordLine("more than words", DefaultMonoCdkStyle)
		So(wl, ShouldNotBeNil)
		So(wl.Value(), ShouldEqual, "more than words")
		So(wl.Words(), ShouldHaveLength, 3)
		So(wl.LetterCount(true), ShouldEqual, 15)
		So(wl.LetterCount(false), ShouldEqual, 13)
	})
}
