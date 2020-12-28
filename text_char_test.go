package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTextChar(t *testing.T) {
	Convey("Basic checks", t, func() {
		tc := NewTextChar([]byte{})
		So(tc, ShouldNotBeNil)
		So(tc.Width(), ShouldEqual, 0)
		So(tc.IsSpace(), ShouldEqual, false)
		tc.Set('*')
		So(tc, ShouldNotBeNil)
		So(tc.Width(), ShouldEqual, 1)
		So(tc.Value(), ShouldEqual, '*')
		So(tc.String(), ShouldEqual, "*")
		So(tc.IsSpace(), ShouldEqual, false)
		tc.SetByte([]byte{' '})
		So(tc, ShouldNotBeNil)
		So(tc.IsSpace(), ShouldEqual, true)
		So(tc.Value(), ShouldEqual, ' ')
		So(tc.String(), ShouldEqual, " ")
	})
}
