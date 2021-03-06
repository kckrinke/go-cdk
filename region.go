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

// Point2I and Rectangle combined

import (
	"fmt"
)

type Region struct {
	Point2I
	Rectangle
}

func NewRegion(x, y, w, h int) *Region {
	r := MakeRegion(x, y, w, h)
	return &r
}

func MakeRegion(x, y, w, h int) Region {
	return Region{
		Point2I{X: x, Y: y},
		Rectangle{W: w, H: h},
	}
}

func (r Region) String() string {
	return fmt.Sprintf("%v,%v", r.Point2I.String(), r.Rectangle.String())
}

func (r *Region) SetRegion(x, y, w, h int) {
	r.X, r.Y, r.W, r.H = x, y, w, h
}

func (r Region) HasPoint(pos Point2I) bool {
	if r.X <= pos.X {
		if r.Y <= pos.Y {
			if (r.X + r.W) >= pos.X {
				if (r.Y + r.H) >= pos.Y {
					return true
				}
			}
		}
	}
	return false
}

func (r Region) Origin() Point2I {
	return Point2I{r.X, r.Y}
}

func (r Region) FarPoint() Point2I {
	return Point2I{
		r.X + r.W,
		r.Y + r.H,
	}
}

func (r Region) Size() Rectangle {
	return Rectangle{r.W, r.H}
}
