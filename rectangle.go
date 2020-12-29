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

// Rectangle

import (
	"fmt"
)

type Rectangle struct {
	W, H int
}

func NewRectangle(w, h int) *Rectangle {
	return &Rectangle{w, h}
}

func (r Rectangle) String() string {
	return fmt.Sprintf("w:%v,h:%v", r.W, r.H)
}

func (r Rectangle) Volume() int {
	return r.W * r.H
}

func (r *Rectangle) SetArea(w, h int) {
	r.W = w
	r.H = h
}

func (r *Rectangle) SetAreaR(size Rectangle) {
	r.W = size.W
	r.H = size.H
}

func (r *Rectangle) AddArea(x, y int) {
	r.W += x
	r.H += y
}

func (r *Rectangle) AddAreaR(size Rectangle) {
	r.W += size.W
	r.H += size.H
}

func (r *Rectangle) SubArea(x, y int) {
	r.W -= x
	r.H -= y
}

func (r *Rectangle) SubAreaR(size Rectangle) {
	r.W -= size.W
	r.H -= size.H
}

func (r *Rectangle) Clamp(region Region) (clamped bool) {
	clamped = false
	min, max := region.Origin(), region.FarPoint()
	// is width within range?
	if r.W >= min.X && r.W <= max.X {
		// width is within range, NOP
	} else {
		// width is not within range
		if r.W < min.X {
			// width is too low, CLAMP
			r.W = min.X
		} else {
			// width is too high, CLAMP
			r.W = max.X
		}
		clamped = true
	}
	// is height within range?
	if r.H >= min.Y && r.H <= max.Y {
		// height is within range, NOP
	} else {
		// height is not within range
		if r.H < min.Y {
			// height is too low, CLAMP
			r.H = min.Y
		} else {
			// height is too high, CLAMP
			r.H = max.Y
		}
		clamped = true
	}
	return
}
