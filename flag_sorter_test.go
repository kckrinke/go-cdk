package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli/v2"
)

func TestFlagSorter(t *testing.T) {
	Convey("Basic FlagSorter Features", t, func() {
		flags := FlagSorter{
			&cli.BoolFlag{
				Name: "cdk-test",
			},
			&cli.BoolFlag{
				Name: "z-test",
			},
			&cli.BoolFlag{
				Name: "a-test",
			},
			&cli.BoolFlag{
				Name: "cdk-two-test",
			},
			&cli.BoolFlag{},
		}
		So(flags.Len(), ShouldEqual, 5)
		So(flags.Less(0, 1), ShouldEqual, false)
		So(flags.Less(1, 2), ShouldEqual, true)
		So(flags.Less(1, 0), ShouldEqual, true)
		So(flags.Less(2, 0), ShouldEqual, true)
		So(flags.Less(0, 3), ShouldEqual, true)
		So(flags.Less(3, 0), ShouldEqual, false)
		So(flags.Less(4, 0), ShouldEqual, true)
		So(flags.Less(0, 4), ShouldEqual, false)
		So(len(flags[4].Names()), ShouldEqual, 1)
		So(flags[4].Names()[0], ShouldEqual, "")
		flags.Swap(0, 1)
		So(flags[0].Names()[0], ShouldEqual, "z-test")
		So(flags[1].Names()[0], ShouldEqual, "cdk-test")
	})
}
