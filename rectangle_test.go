package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRectangle(t *testing.T) {
	Convey("Basic Rectangle Features", t, func() {
		r := NewRectangle(2, 2)
		So(r, ShouldNotBeNil)
		So(r, ShouldHaveSameTypeAs, &Rectangle{})
		So(r.String(), ShouldEqual, "w:2,h:2")
		r.SetArea(1, 1)
		So(r.W, ShouldEqual, 1)
		So(r.H, ShouldEqual, 1)
		So(r.Volume(), ShouldEqual, 1)
		r.SetAreaR(Rectangle{2, 2})
		So(r.W, ShouldEqual, 2)
		So(r.H, ShouldEqual, 2)
		So(r.Volume(), ShouldEqual, 4)
		r.AddArea(1, 1)
		So(r.W, ShouldEqual, 3)
		So(r.H, ShouldEqual, 3)
		So(r.Volume(), ShouldEqual, 9)
		r.AddAreaR(Rectangle{1, 1})
		So(r.W, ShouldEqual, 4)
		So(r.H, ShouldEqual, 4)
		So(r.Volume(), ShouldEqual, 16)
		r.SubArea(1, 1)
		So(r.W, ShouldEqual, 3)
		So(r.H, ShouldEqual, 3)
		So(r.Volume(), ShouldEqual, 9)
		r.SubAreaR(Rectangle{1, 1})
		So(r.W, ShouldEqual, 2)
		So(r.H, ShouldEqual, 2)
		So(r.Volume(), ShouldEqual, 4)
		r.SetArea(10, 10)
		So(r.Volume(), ShouldEqual, 100)
		region := Region{Point2I{5, 5}, Rectangle{5, 5}}
		So(r.Clamp(region), ShouldEqual, false)
		So(r.Volume(), ShouldEqual, 100)
		region.SetRegion(0, 0, 2, 2)
		So(r.Clamp(region), ShouldEqual, true)
		So(r.Volume(), ShouldEqual, 4)
		region.SetRegion(5, 5, 2, 2)
		So(r.Clamp(region), ShouldEqual, true)
		So(r.Volume(), ShouldEqual, 25)
	})
}
