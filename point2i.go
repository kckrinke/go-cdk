package cdk

import (
	"fmt"
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

func (p Point2I) String() string {
	return fmt.Sprintf("x:%v,y:%v", p.X, p.Y)
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
