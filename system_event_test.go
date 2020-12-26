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
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const SYSTEM_EVENT_TEST SystemEventType = "test"
const SYSTEM_EVENT_UNKNOWN SystemEventType = "unknown"

func TestSystemEvent(t *testing.T) {
	Convey("System Event checks", t, func() {
		err := AddSystemEventHandler(SYSTEM_EVENT_TEST, "test.event", func() error {
			return fmt.Errorf("test.event.error")
		})
		So(err, ShouldBeNil)
		err = AddSystemEventHandler(SYSTEM_EVENT_TEST, "test.event", func() error {
			return fmt.Errorf("test.event.fail")
		})
		So(err, ShouldNotBeNil)
		err = HandleSystemEvent(SYSTEM_EVENT_TEST)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "system event test error: [test.event] test.event.error")
		_ = DelSystemEventHandler(SYSTEM_EVENT_TEST, "test.event")
		err = HandleSystemEvent(SYSTEM_EVENT_TEST)
		So(err, ShouldBeNil)
		err = DelSystemEventHandler(SYSTEM_EVENT_TEST, "test.event")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "system test event, tag not found: test.event")
		err = DelSystemEventHandler(SYSTEM_EVENT_UNKNOWN, "test.unknown")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "system event not found: unknown")
	})
}
