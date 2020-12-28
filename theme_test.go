package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTheme(t *testing.T) {
	Convey("Basic Theme Features", t, func() {
		So(DefaultMonoTheme.String(), ShouldEqual, "{Normal={fg=unnamed[-1],bg=unnamed[-1],attrs=16},Border={fg=unnamed[-1],bg=unnamed[-1],attrs=16},FillRune=32,BorderRunes:{BorderRunes=9488,9472,9484,9474,9492,9472,9496,9474},Overlay=false}")
	})
}
