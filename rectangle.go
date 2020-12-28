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

func (r Rectangle) Area() int {
	return r.W * r.H
}

func (r *Rectangle) Set(w, h int) {
	r.W = w
	r.H = h
}
