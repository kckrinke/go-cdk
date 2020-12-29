package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegion(t *testing.T) {
	Convey("Basic Region Features", t, func() {
		r := NewRegion(1, 1, 2, 2)
		So(r, ShouldNotBeNil)
		So(r, ShouldHaveSameTypeAs, &Region{})
		So(r.String(), ShouldEqual, "x:1,y:1,w:2,h:2")
		So(r.Origin(), ShouldHaveSameTypeAs, Point2I{})
		So(r.Origin().X, ShouldEqual, 1)
		So(r.Origin().Y, ShouldEqual, 1)
		So(r.Size(), ShouldHaveSameTypeAs, Rectangle{})
		So(r.Size().W, ShouldEqual, 2)
		So(r.Size().H, ShouldEqual, 2)
		So(r.FarPoint().X, ShouldEqual, 3)
		So(r.FarPoint().Y, ShouldEqual, 3)
		So(r.HasPoint(Point2I{0, 0}), ShouldEqual, false)
		So(r.HasPoint(Point2I{1, 1}), ShouldEqual, true)
		r.SetRegion(2, 2, 4, 4)
		So(r.String(), ShouldEqual, "x:2,y:2,w:4,h:4")
	})
}
