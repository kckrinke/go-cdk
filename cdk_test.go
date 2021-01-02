package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func makesNoContent(d DisplayManager) error {
	return nil
}

func TestCdk(t *testing.T) {
	Convey("Making a new app instance", t, func() {
		Convey("validating factory", func() {
			app := NewApp(
				"AppName", "AppUsage", "v0.0.0",
				"app-tag", "AppTitle",
				OffscreenDisplayTtyPath,
				makesNoContent,
			)
			So(app, ShouldNotBeNil)
			So(app.Name(), ShouldEqual, "AppName")
			So(app.Usage(), ShouldEqual, "AppUsage")
			So(app.Title(), ShouldEqual, "AppTitle")
			So(app.Version(), ShouldEqual, "v0.0.0")
			So(app.Tag(), ShouldEqual, "app-tag")
			So(app.GetContext(), ShouldBeNil)
			So(app.CLI(), ShouldNotBeNil)
			app.Destroy()
		})
		Convey("with no content", WithApp(
			makesNoContent,
			func(d App) {
				// do tests here?
				So(d.DisplayManager(), ShouldNotBeNil)
			},
		))
	})
}
