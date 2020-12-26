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

	"github.com/gobuffalo/envy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCdk(t *testing.T) {
	envy.Set("GO_CDK_LOG_FORMAT", "pretty")
	envy.Set("GO_CDK_LOG_LEVEL", "error")
	envy.Set("GO_CDK_LOG_OUTPUT", "file")
	Convey("CDK Main checks", t, func() {
		So(MainScreen(), ShouldBeNil)
		_, _, err := DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "error")
			envy.Set("GO_CDK_LOG_OUTPUT", "file")
			return MainInit()
		})
		So(err, ShouldBeNil)
		So(MainScreen(), ShouldNotBeNil)
		_, faked_it, err := DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "error")
			envy.Set("GO_CDK_LOG_OUTPUT", "file")
			err := MainInit()
			MainQuit()
			return err
		})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "non-operation: main already initialized")
		So(faked_it, ShouldEqual, true)
		So(MainScreen(), ShouldBeNil)
	})
}
