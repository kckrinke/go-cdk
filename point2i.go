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
	return &Point2I{x, y}
}

func (p Point2I) String() string {
	return fmt.Sprintf("x:%v,y:%v", p.X, p.Y)
}

func (p *Point2I) Set(x, y int) {
	p.X = x
	p.Y = y
}

func (p *Point2I) Add(v Point2I) {
	p.X += v.X
	p.Y += v.Y
}

func (p *Point2I) Sub(v Point2I) {
	p.X -= v.X
	p.Y -= v.Y
}
