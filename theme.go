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
	DefaultBorderRune          = BorderRune{
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
	DefaultNilTheme  = Theme{}
	DefaultMonoTheme = Theme{
		Normal:      DefaultMonoCdkStyle,
		Border:      DefaultMonoCdkStyle.Dim(true),
		Focused:     DefaultMonoCdkStyle.Dim(false),
		Active:      DefaultMonoCdkStyle.Dim(false).Reverse(true),
		FillRune:    DefaultFillRune,
		BorderRunes: DefaultBorderRune,
		Overlay:     false,
	}
	DefaultColorTheme = Theme{
		Normal:      DefaultColorCdkStyle.Dim(true),
		Border:      DefaultColorCdkStyle.Dim(true),
		Focused:     DefaultColorCdkStyle.Dim(false),
		Active:      DefaultColorCdkStyle.Dim(false).Reverse(true),
		FillRune:    DefaultFillRune,
		BorderRunes: DefaultBorderRune,
		Overlay:     false,
	}
)

type BorderRune struct {
	TopLeft     rune
	Top         rune
	TopRight    rune
	Left        rune
	Right       rune
	BottomLeft  rune
	Bottom      rune
	BottomRight rune
}

func (b BorderRune) String() string {
	return fmt.Sprintf(
		"{BorderRunes=%v,%v,%v,%v,%v,%v,%v,%v}",
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

type Theme struct {
	Normal      Style
	Border      Style
	Focused     Style
	Active      Style
	FillRune    rune
	BorderRunes BorderRune
	Overlay     bool // keep existing background
}

func (t Theme) String() string {
	return fmt.Sprintf(
		"{Normal=%v,Border=%v,Focused=%v,Active=%v,FillRune=%v,BorderRunes=%v,Overlay=%v}",
		t.Normal,
		t.Border,
		t.Focused,
		t.Active,
		t.FillRune,
		t.BorderRunes,
		t.Overlay,
	)
}
