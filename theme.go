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
)

var (
	DefaultFillRune      rune  = ' '
	DefaultMonoCdkStyle  Style = StyleDefault.Dim(true)
	DefaultColorCdkStyle Style = StyleDefault.Foreground(ColorWhite).Background(ColorNavy)
	DefaultBorderRune          = FrameRunes{
		TopLeft:     RuneULCorner,
		Top:         RuneHLine,
		TopRight:    RuneURCorner,
		Left:        RuneVLine,
		Right:       RuneVLine,
		BottomLeft:  RuneLLCorner,
		Bottom:      RuneHLine,
		BottomRight: RuneLRCorner,
	}
)

var (
	DefaultNilTheme  = &CTheme{}
	DefaultMonoTheme = NewTheme(
		DefaultMonoCdkStyle,
		DefaultMonoCdkStyle.Dim(true),
		DefaultFillRune,
		DefaultBorderRune,
		false,
	)
	DefaultColorTheme = NewTheme(
		DefaultColorCdkStyle,
		DefaultColorCdkStyle.Dim(true),
		DefaultFillRune,
		DefaultBorderRune,
		false,
	)
)

type FrameRunes struct {
	TopLeft     rune
	Top         rune
	TopRight    rune
	Left        rune
	Right       rune
	BottomLeft  rune
	Bottom      rune
	BottomRight rune
}

func (b FrameRunes) String() string {
	return fmt.Sprintf(
		"{Frame=%v,%v,%v,%v,%v,%v,%v,%v}",
		b.TopRight,
		b.Top,
		b.TopLeft,
		b.Left,
		b.BottomLeft,
		b.Bottom,
		b.BottomRight,
		b.Right,
	)
}

type Theme interface {
	GetNormal() Style
	SetNormal(normal Style) Theme
	GetBorder() Style
	SetBorder(border Style) Theme
	GetFillRune() rune
	SetFillRune(fill rune) Theme
	GetFrame() FrameRunes
	SetFrame(frame FrameRunes) Theme
	GetOverlay() bool
	SetOverlay(overlay bool) Theme
	String() string
}

type CTheme struct {
	Normal   Style
	Border   Style
	FillRune rune
	Frame    FrameRunes
	Overlay  bool // keep existing background
}

func NewTheme(normal, border Style, fill rune, frame FrameRunes, overlay bool) Theme {
	return &CTheme{
		Normal:   normal,
		Border:   border,
		FillRune: fill,
		Frame:    frame,
		Overlay:  overlay,
	}
}

func (t CTheme) GetNormal() Style {
	return t.Normal
}

func (t *CTheme) SetNormal(normal Style) Theme {
	t.Normal = normal
	return t
}

func (t CTheme) GetBorder() Style {
	return t.Border
}

func (t *CTheme) SetBorder(border Style) Theme {
	t.Border = border
	return t
}

func (t CTheme) GetFillRune() rune {
	return t.FillRune
}

func (t *CTheme) SetFillRune(fill rune) Theme {
	t.FillRune = fill
	return t
}

func (t CTheme) GetFrame() FrameRunes {
	return t.Frame
}

func (t *CTheme) SetFrame(frame FrameRunes) Theme {
	t.Frame = frame
	return t
}

func (t CTheme) GetOverlay() bool {
	return t.Overlay
}

func (t *CTheme) SetOverlay(overlay bool) Theme {
	t.Overlay = overlay
	return t
}

func (t CTheme) String() string {
	return fmt.Sprintf(
		"{Normal=%v,Border=%v,FillRune=%v,Frame:%v,Overlay=%v}",
		t.Normal,
		t.Border,
		t.FillRune,
		t.Frame,
		t.Overlay,
	)
}
