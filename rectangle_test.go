// Copyright 2020 The CDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRectangle(t *testing.T) {
	Convey("Basic Rectangle Features", t, func() {
		r := NewRectangle(2, 2)
		So(r, ShouldNotBeNil)
		So(r, ShouldHaveSameTypeAs, &Rectangle{})
		So(r.String(), ShouldEqual, "w:2,h:2")
		r.SetArea(1, 1)
		So(r.W, ShouldEqual, 1)
		So(r.H, ShouldEqual, 1)
		So(r.Volume(), ShouldEqual, 1)
		r.SetAreaR(Rectangle{2, 2})
		So(r.W, ShouldEqual, 2)
		So(r.H, ShouldEqual, 2)
		So(r.Volume(), ShouldEqual, 4)
		r.AddArea(1, 1)
		So(r.W, ShouldEqual, 3)
		So(r.H, ShouldEqual, 3)
		So(r.Volume(), ShouldEqual, 9)
		r.AddAreaR(Rectangle{1, 1})
		So(r.W, ShouldEqual, 4)
		So(r.H, ShouldEqual, 4)
		So(r.Volume(), ShouldEqual, 16)
		r.SubArea(1, 1)
		So(r.W, ShouldEqual, 3)
		So(r.H, ShouldEqual, 3)
		So(r.Volume(), ShouldEqual, 9)
		r.SubAreaR(Rectangle{1, 1})
		So(r.W, ShouldEqual, 2)
		So(r.H, ShouldEqual, 2)
		So(r.Volume(), ShouldEqual, 4)
		r.SetArea(10, 10)
		So(r.Volume(), ShouldEqual, 100)
		region := Region{Point2I{5, 5}, Rectangle{5, 5}}
		So(r.ClampRegion(region), ShouldEqual, false)
		So(r.Volume(), ShouldEqual, 100)
		region.SetRegion(0, 0, 2, 2)
		So(r.ClampRegion(region), ShouldEqual, true)
		So(r.Volume(), ShouldEqual, 4)
		region.SetRegion(5, 5, 2, 2)
		So(r.ClampRegion(region), ShouldEqual, true)
		So(r.Volume(), ShouldEqual, 25)
	})
}
