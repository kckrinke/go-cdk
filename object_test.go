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

func TestObject(t *testing.T) {
	Convey("Basic Object Features", t, func() {
		o := &CObject{}
		So(o, ShouldImplement, (*Object)(nil))
		So(o.Init(), ShouldEqual, false)
		So(o.Init(), ShouldEqual, true)
		// normal testing
		So(o.GetTheme().String(), ShouldEqual, DefaultColorTheme.String())
		o.SetTheme(DefaultMonoTheme)
		So(o.GetTheme().String(), ShouldEqual, DefaultMonoTheme.String())
		o.SetProperty("testing", nil)
		So(o.GetProperty("testing"), ShouldBeNil)
		So(o.GetPropertyAsBool("testing", true), ShouldEqual, true)
		So(o.GetPropertyAsInt("testing", 1), ShouldEqual, 1)
		So(o.GetPropertyAsString("testing", "one"), ShouldEqual, "one")
		So(o.GetPropertyAsFloat("testing", 1.0), ShouldEqual, 1.0)
		o.SetProperty("testing", true)
		So(o.GetPropertyAsBool("testing", false), ShouldEqual, true)
		o.SetProperty("testing", 1)
		So(o.GetPropertyAsInt("testing", 0), ShouldEqual, 1)
		o.SetProperty("testing", "one")
		So(o.GetPropertyAsString("testing", "nil"), ShouldEqual, "one")
		o.SetProperty("testing", 1.0)
		So(o.GetPropertyAsFloat("testing", 0.0), ShouldEqual, 1.0)
		So(o.GetProperty("test"), ShouldBeNil)
		// destruction testing
		hit0 := false
		o.Connect(SignalDestroy, "basic-destroy", func(data []interface{}, argv ...interface{}) EventFlag {
			hit0 = true
			return EVENT_STOP
		})
		o.Destroy()
		So(hit0, ShouldEqual, true)
		So(o.IsValid(), ShouldEqual, true)
		hit0 = false
		o.Disconnect(SignalDestroy, "basic-destroy")
		o.Connect(SignalDestroy, "basic-destroy", func(data []interface{}, argv ...interface{}) EventFlag {
			hit0 = true
			return EVENT_PASS
		})
		o.Destroy()
		So(hit0, ShouldEqual, true)
		So(o.IsValid(), ShouldEqual, false)
	})
}
