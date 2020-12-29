package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPoint2I(t *testing.T) {
	Convey("Basic Point2I Features", t, func() {
		p := NewPoint2I(2, 2)
		So(p, ShouldNotBeNil)
		So(p, ShouldHaveSameTypeAs, &Point2I{})
		So(p.String(), ShouldEqual, "x:2,y:2")
		p.SetPoint(1, 1)
		So(p.X, ShouldEqual, 1)
		So(p.Y, ShouldEqual, 1)
		p.SetPoint2I(Point2I{2, 2})
		So(p.X, ShouldEqual, 2)
		So(p.Y, ShouldEqual, 2)
		p.AddPoint(1, 1)
		So(p.X, ShouldEqual, 3)
		So(p.Y, ShouldEqual, 3)
		p.AddPoint2I(Point2I{1, 1})
		So(p.X, ShouldEqual, 4)
		So(p.Y, ShouldEqual, 4)
		p.SubPoint(1, 1)
		So(p.X, ShouldEqual, 3)
		So(p.Y, ShouldEqual, 3)
		p.SubPoint2I(Point2I{1, 1})
		So(p.X, ShouldEqual, 2)
		So(p.Y, ShouldEqual, 2)
		p.SetPoint(10, 10)
		region := Region{Point2I{5, 5}, Rectangle{5, 5}}
		So(p.Clamp(region), ShouldEqual, false)
		So(p.String(), ShouldEqual, "x:10,y:10")
		region.SetRegion(0, 0, 2, 2)
		So(p.Clamp(region), ShouldEqual, true)
		So(p.String(), ShouldEqual, "x:2,y:2")
		region.SetRegion(5, 5, 2, 2)
		So(p.Clamp(region), ShouldEqual, true)
		So(p.String(), ShouldEqual, "x:5,y:5")
	})
}
