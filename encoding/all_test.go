package encoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/kckrinke/go-cdk"
)

func TestRegisterAll(t *testing.T) {
	Convey("Encoding Register All", t, func() {
		list := cdk.ListEncodings()
		So(len(list), ShouldEqual, 5)
		So(list, ShouldContain, "utf-8")
		So(list, ShouldContain, "utf8")
		So(list, ShouldContain, "us-ascii")
		So(list, ShouldContain, "ascii")
		So(list, ShouldContain, "iso646")
		Register()
		list = cdk.ListEncodings()
		So(len(list), ShouldEqual, 61)
		So(list, ShouldContain, "euc-jp")
	})
}