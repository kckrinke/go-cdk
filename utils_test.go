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

func TestUtils(t *testing.T) {
	Convey("Utility func checks", t, func() {
		padded := pad_left("thing", " ", 6)
		So(padded[0], ShouldEqual, ' ')
		padded = pad_right("thing", " ", 6)
		So(padded[5], ShouldEqual, ' ')
		nlstr := "thing\n\r\n"
		So(clean_crlf(nlstr), ShouldEqual, "thing")
		So(nlsprintf("%s\n\r\n", "thing"), ShouldEqual, "thing")
		logged, faked_it, err := GetLastFakeIO()
		So(logged, ShouldHaveSameTypeAs, "")
		So(faked_it, ShouldEqual, -1)
		So(err, ShouldBeNil)
	})
}
