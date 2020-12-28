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
	return &Region{
		Point2I{x, y},
		Rectangle{w, h},
	}
}

func (r Region) String() string {
	return fmt.Sprintf("%v, %v", r.Point2I.String(), r.Rectangle.String())
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
