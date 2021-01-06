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

func TestTheme(t *testing.T) {
	Convey("Basic Theme Features", t, func() {
		So(
			DefaultMonoTheme.String(),
			ShouldEqual,
			"{Normal={fg=unnamed[-1],bg=unnamed[-1],attrs=16},Border={fg=unnamed[-1],bg=unnamed[-1],attrs=16},Focused={fg=unnamed[-1],bg=unnamed[-1],attrs=0},Active={fg=unnamed[-1],bg=unnamed[-1],attrs=4},FillRune=32,BorderRunes={BorderRunes=9488,9472,9484,9474,9492,9472,9496,9474},Overlay=false}",
		)
	})
}
