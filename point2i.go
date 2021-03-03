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
	"regexp"
	"strconv"
)

// Point2I

/*
Point2I in 2D space, denoted by X and Y integers
*/

type Point2I struct {
	X, Y int
}

func NewPoint2I(x, y int) *Point2I {
	r := MakePoint2I(x, y)
	return &r
}

func MakePoint2I(x, y int) Point2I {
	return Point2I{X: x, Y: y}
}

var rxParsePoint2I = regexp.MustCompile(`(?:i)^{??(?:x:)??(\d+),(?:y:)??(\d+)}??$`)

func ParsePoint2I(value string) (point Point2I, ok bool) {
	if rxParsePoint2I.MatchString(value) {
		m := rxParsePoint2I.FindStringSubmatch(value)
		if len(m) == 3 {
			x, _ := strconv.Atoi(m[1])
			y, _ := strconv.Atoi(m[2])
			return MakePoint2I(x, y), true
		}
	}
	return Point2I{}, false
}

func (p Point2I) String() string {
	return fmt.Sprintf("{x:%v,y:%v}", p.X, p.Y)
}

func (p Point2I) Clone() (clone Point2I) {
	clone.X = p.X
	clone.Y = p.Y
	return
}

func (p Point2I) Equals(x, y int) bool {
	return p.X == x && p.Y == y
}

func (p Point2I) Equals2I(o Point2I) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p *Point2I) SetPoint(x, y int) {
	p.X = x
	p.Y = y
}

func (p *Point2I) SetPoint2I(point Point2I) {
	p.X = point.X
	p.Y = point.Y
}

func (p *Point2I) AddPoint(x, y int) {
	p.X += x
	p.Y += y
}

func (p *Point2I) AddPoint2I(point Point2I) {
	p.X += point.X
	p.Y += point.Y
}

func (p *Point2I) SubPoint(x, y int) {
	p.X -= x
	p.Y -= y
}

func (p *Point2I) SubPoint2I(point Point2I) {
	p.X -= point.X
	p.Y -= point.Y
}

func (p *Point2I) Clamp(region Region) (clamped bool) {
	clamped = false
	min, max := region.Origin(), region.FarPoint()
	// is width within range?
	if p.X >= min.X && p.X <= max.X {
		// width is within range, NOP
	} else {
		// width is not within range
		if p.X < min.X {
			// width is too low, CLAMP
			p.X = min.X
		} else {
			// width is too high, CLAMP
			p.X = max.X
		}
		clamped = true
	}
	// is height within range?
	if p.Y >= min.Y && p.Y <= max.Y {
		// height is within range, NOP
	} else {
		// height is not within range
		if p.Y < min.Y {
			// height is too low, CLAMP
			p.Y = min.Y
		} else {
			// height is too high, CLAMP
			p.Y = max.Y
		}
		clamped = true
	}
	return
}
