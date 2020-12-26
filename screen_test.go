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
	"os"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScreen(t *testing.T) {
	Convey("New Screen checks", t, func() {
		if runtime.GOOS == "windows" {
			t.Log("NewScreen() tests unimplemented on windows")
		} else {
			ns, err := NewScreen()
			So(err, ShouldEqual, nil)
			nts, err := NewTerminfoScreen()
			So(err, ShouldEqual, nil)
			So(nts, ShouldNotEqual, nil)
			ncs, err := NewConsoleScreen()
			So(err, ShouldNotEqual, nil)
			So(ncs, ShouldEqual, nil)
			So(ns, ShouldHaveSameTypeAs, nts)
			os.Setenv("TERM", "")
			ns, err = NewScreen()
			So(err, ShouldNotEqual, nil)
			So(ns, ShouldEqual, nil)
			ncs, err = NewConsoleScreen()
			So(err, ShouldNotEqual, nil)
			So(ncs, ShouldEqual, nil)
		}
	})
}
