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
	"regexp"
	"strconv"
)

type Rectangle struct {
	W, H int
}

func NewRectangle(w, h int) *Rectangle {
	r := MakeRectangle(w, h)
	return &r
}

func MakeRectangle(w, h int) Rectangle {
	return Rectangle{W: w, H: h}
}

var rxParseRectangle = regexp.MustCompile(`(?:i)^{??(?:w:)??(\d+),(?:h:)??(\d+)}??$`)

func ParseRectangle(value string) (point Rectangle, ok bool) {
	if rxParseRectangle.MatchString(value) {
		m := rxParseRectangle.FindStringSubmatch(value)
		if len(m) == 3 {
			w, _ := strconv.Atoi(m[1])
			h, _ := strconv.Atoi(m[2])
			return MakeRectangle(w, h), true
		}
	}
	return Rectangle{}, false
}

func (r Rectangle) String() string {
	return fmt.Sprintf("{w:%v,h:%v}", r.W, r.H)
}

func (r Rectangle) Clone() (clone Rectangle) {
	clone.W = r.W
	clone.H = r.H
	return
}

func (r Rectangle) Equals(w, h int) bool {
	return r.W == w && r.H == h
}

func (r Rectangle) EqualsR(o Rectangle) bool {
	return r.W == o.W && r.H == o.H
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

func (r *Rectangle) AddArea(w, h int) {
	r.W += w
	r.H += h
}

func (r *Rectangle) AddAreaR(size Rectangle) {
	r.W += size.W
	r.H += size.H
}

func (r *Rectangle) SubArea(w, h int) {
	r.W -= w
	r.H -= h
}

func (r *Rectangle) SubAreaR(size Rectangle) {
	r.W -= size.W
	r.H -= size.H
}

func (r *Rectangle) Floor(minWidth, minHeight int) {
	if r.W < minWidth {
		r.W = minWidth
	}
	if r.H < minHeight {
		r.H = minHeight
	}
}

func (r *Rectangle) Clamp(minWidth, minHeight, maxWidth, maxHeight int) {
	if r.W < minWidth {
		r.W = minWidth
	}
	if r.H < minHeight {
		r.H = minHeight
	}
	if r.W > maxWidth {
		r.W = maxWidth
	}
	if r.H > maxHeight {
		r.H = maxHeight
	}
}

func (r *Rectangle) ClampRegion(region Region) (clamped bool) {
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
