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

func TestTextBuffer(t *testing.T) {
	Convey("Text Buffers with...", t, func() {
		Convey("Basic checks", func() {
			tb := NewEmptyTextBuffer(DefaultMonoStyle)
			So(tb, ShouldNotBeNil)
			So(tb.Style().String(), ShouldEqual, DefaultMonoStyle.String())
			So(tb.CharacterCount(), ShouldEqual, 0)
			So(tb.WordCount(), ShouldEqual, 0)
			tb = NewTextBuffer("test", DefaultMonoStyle)
			So(tb, ShouldNotBeNil)
			So(tb.CharacterCount(), ShouldEqual, 4)
			So(tb.WordCount(), ShouldEqual, 1)
		})
		Convey("Draw checks", func() {
			tb := NewEmptyTextBuffer(DefaultMonoStyle)
			So(tb, ShouldNotBeNil)
			canvas := NewCanvas(Point2I{}, Rectangle{10, 3}, DefaultMonoTheme)
			f := tb.Draw(canvas, true, WRAP_NONE, JUSTIFY_LEFT, ALIGN_TOP)
			So(f, ShouldEqual, EVENT_PASS)

			tb = NewTextBuffer("test", DefaultMonoStyle)
			So(tb, ShouldNotBeNil)
			canvas = NewCanvas(Point2I{}, Rectangle{10, 3}, DefaultMonoTheme)
			f = tb.Draw(canvas, true, WRAP_NONE, JUSTIFY_LEFT, ALIGN_TOP)
			So(f, ShouldEqual, EVENT_STOP)
			val := ""
			numSpaces := 0
			for x := 0; x < 10; x++ {
				if c := canvas.GetContent(x, 0); !c.IsSpace() {
					val += string(c.Value())
				} else {
					numSpaces++
				}
			}
			So(val, ShouldEqual, "test")
			So(numSpaces, ShouldEqual, 6)

			canvas = NewCanvas(Point2I{}, Rectangle{10, 3}, DefaultMonoTheme)
			f = tb.Draw(canvas, true, WRAP_NONE, JUSTIFY_LEFT, ALIGN_BOTTOM)
			So(f, ShouldEqual, EVENT_STOP)
			val = ""
			numSpaces = 0
			for x := 0; x < 10; x++ {
				if c := canvas.GetContent(x, 2); !c.IsSpace() {
					val += string(c.Value())
				} else {
					numSpaces++
				}
			}
			So(val, ShouldEqual, "test")
			So(numSpaces, ShouldEqual, 6)

			canvas = NewCanvas(Point2I{}, Rectangle{10, 3}, DefaultMonoTheme)
			f = tb.Draw(canvas, true, WRAP_NONE, JUSTIFY_LEFT, ALIGN_MIDDLE)
			So(f, ShouldEqual, EVENT_STOP)
			val = ""
			numSpaces = 0
			for x := 0; x < 10; x++ {
				if c := canvas.GetContent(x, 1); !c.IsSpace() {
					val += string(c.Value())
				} else {
					numSpaces++
				}
			}
			So(val, ShouldEqual, "test")
			So(numSpaces, ShouldEqual, 6)
		})
	})
}
