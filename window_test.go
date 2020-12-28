package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWindow(t *testing.T) {
	Convey("Basic Window Features", t, func() {
		w := &CWindow{}
		So(w.valid, ShouldEqual, false)
		So(w.GetTitle(), ShouldEqual, "")
		So(w.GetDisplay(), ShouldBeNil)
		So(w.Init(), ShouldEqual, false)
		So(w.Init(), ShouldEqual, true)
		d := &CDisplay{}
		w.SetDisplay(d)
		So(w.GetDisplay(), ShouldEqual, d)
		w.SetTitle("testing")
		So(w.GetTitle(), ShouldEqual, "testing")
		So(w.Draw(&Canvas{}), ShouldEqual, EVENT_PASS)
		So(w.ProcessEvent(&EventError{}), ShouldEqual, EVENT_PASS)
	})
}
