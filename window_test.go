package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWindow(t *testing.T) {
	Convey("Basic Window Features", t, func() {
		So(TypesManager.HasType(TypeWindow), ShouldEqual, true)
		w := &CWindow{}
		So(w.valid, ShouldEqual, false)
		So(w.GetTitle(), ShouldEqual, "")
		So(w.GetDisplayManager(), ShouldBeNil)
		So(w.Init(), ShouldEqual, false)
		So(w.Init(), ShouldEqual, true)
		d := &CDisplayManager{}
		w.SetDisplayManager(d)
		So(w.GetDisplayManager(), ShouldEqual, d)
		w.SetTitle("testing")
		So(w.GetTitle(), ShouldEqual, "testing")
		So(w.Draw(&CCanvas{}), ShouldEqual, EVENT_PASS)
		So(w.ProcessEvent(&EventError{}), ShouldEqual, EVENT_PASS)
	})
}
