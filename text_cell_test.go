package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTextCell(t *testing.T) {
	Convey("Basic checks", t, func() {
		tc := NewRuneCell('*', DefaultMonoCdkStyle)
		So(tc, ShouldNotBeNil)
		So(tc.IsSpace(), ShouldEqual, false)
		So(tc.Width(), ShouldEqual, 1)
		So(tc.Style().String(), ShouldEqual, DefaultMonoCdkStyle.String())
		So(tc.Dirty(), ShouldEqual, false)
		So(tc.Value(), ShouldEqual, '*')
		So(tc.String(), ShouldEqual, "{Char=*,Style={fg=unnamed[-1],bg=unnamed[-1],attrs=16}}")
		tc.Set(' ')
		So(tc.IsSpace(), ShouldEqual, true)
		So(tc.Width(), ShouldEqual, 1)
		So(tc.Dirty(), ShouldEqual, true)
		tc.SetByte([]byte{'0'})
		So(tc.IsSpace(), ShouldEqual, false)
		So(tc.Dirty(), ShouldEqual, true)
		tc.SetStyle(DefaultColorCdkStyle)
		So(tc.Dirty(), ShouldEqual, true)
		So(tc.Style().String(), ShouldEqual, DefaultColorCdkStyle.String())
	})
}
