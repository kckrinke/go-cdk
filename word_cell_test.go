package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWordCell(t *testing.T) {
	Convey("Empty Strings", t, func() {
		wc, err := NewWordCell("", DefaultMonoCdkStyle)
		So(err, ShouldBeNil)
		So(wc, ShouldNotBeNil)
		So(wc.String(), ShouldEqual, "")
		So(wc.Characters(), ShouldHaveLength, 0)
		So(wc.Len(), ShouldEqual, 0)
	})
	Convey("One Word", t, func() {
		wc, err := NewWordCell("word", DefaultMonoCdkStyle)
		So(err, ShouldBeNil)
		So(wc, ShouldNotBeNil)
		So(wc.Value(), ShouldEqual, "word")
		So(wc.Characters(), ShouldHaveLength, 4)
		So(wc.Len(), ShouldEqual, 4)
	})
	Convey("More Than One Word", t, func() {
		wc, err := NewWordCell("another word", DefaultMonoCdkStyle)
		So(err, ShouldNotBeNil)
		So(wc, ShouldBeNil)
		So(err.Error(), ShouldEqual, "words cannot contain spaces")
	})
}
