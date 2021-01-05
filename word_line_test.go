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
		So(wl.Len(), ShouldEqual, 0)
	})
	Convey("One Word", t, func() {
		wl := NewWordLine("word", DefaultMonoCdkStyle)
		So(wl, ShouldNotBeNil)
		So(wl.Value(), ShouldEqual, "word")
		So(wl.Words(), ShouldHaveLength, 1)
		So(wl.CharacterCount(), ShouldEqual, 4)
		So(wl.WordCount(), ShouldEqual, 1)
	})
	Convey("More Than One Word", t, func() {
		wl := NewWordLine("more than words", DefaultMonoCdkStyle)
		So(wl, ShouldNotBeNil)
		So(wl.Value(), ShouldEqual, "more   than   words")
		So(wl.Words(), ShouldHaveLength, 5)
		So(wl.CharacterCount(), ShouldEqual, 15)
	})
}
