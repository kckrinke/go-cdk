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
	DefaultFillRune            rune  = ' '
	DefaultMonoStyle           Style = StyleDefault.Reverse(false).Dim(true)
	DefaultMonoAccessoryStyle  Style = StyleDefault.Reverse(true).Dim(true)
	DefaultColorStyle          Style = StyleDefault.Foreground(ColorWhite).Background(ColorNavy)
	DefaultColorAccessoryStyle Style = StyleDefault.Foreground(ColorWhite).Background(ColorDarkSlateBlue)
	DefaultBorderRune                = BorderRuneSet{
		TopLeft:     RuneULCorner,
		Top:         RuneHLine,
		TopRight:    RuneURCorner,
		Left:        RuneVLine,
		Right:       RuneVLine,
		BottomLeft:  RuneLLCorner,
		Bottom:      RuneHLine,
		BottomRight: RuneLRCorner,
	}
	DefaultArrowRune = ArrowRuneSet{
		Up:    RuneUArrow,
		Left:  RuneLArrow,
		Down:  RuneDArrow,
		Right: RuneRArrow,
	}
	FancyArrowRune = ArrowRuneSet{
		Up:    RuneBlackMediumUpPointingTriangleCentred,
		Left:  RuneBlackMediumLeftPointingTriangleCentred,
		Down:  RuneBlackMediumDownPointingTriangleCentred,
		Right: RuneBlackMediumRightPointingTriangleCentred,
	}
)

var (
	DefaultNilTheme  = Theme{}
	DefaultMonoTheme = Theme{
		Normal:      DefaultMonoStyle,
		Border:      DefaultMonoStyle.Dim(true),
		Focused:     DefaultMonoStyle.Dim(false),
		Active:      DefaultMonoStyle.Dim(false).Reverse(true),
		Accessory:   DefaultMonoAccessoryStyle,
		FillRune:    DefaultFillRune,
		BorderRunes: DefaultBorderRune,
		ArrowRunes:  DefaultArrowRune,
		Overlay:     false,
	}
	DefaultColorTheme = Theme{
		Normal:      DefaultColorStyle.Dim(true),
		Border:      DefaultColorStyle.Dim(true),
		Focused:     DefaultColorStyle.Dim(false),
		Active:      DefaultColorStyle.Dim(false).Reverse(true),
		Accessory:   DefaultColorAccessoryStyle,
		FillRune:    DefaultFillRune,
		BorderRunes: DefaultBorderRune,
		ArrowRunes:  DefaultArrowRune,
		Overlay:     false,
	}
)

type BorderRuneSet struct {
	TopLeft     rune
	Top         rune
	TopRight    rune
	Left        rune
	Right       rune
	BottomLeft  rune
	Bottom      rune
	BottomRight rune
}

func (b BorderRuneSet) String() string {
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

type ArrowRuneSet struct {
	Up    rune
	Left  rune
	Down  rune
	Right rune
}

func (b ArrowRuneSet) String() string {
	return fmt.Sprintf(
		"{ArrowRunes=%v,%v,%v,%v}",
		b.Up,
		b.Left,
		b.Down,
		b.Right,
	)
}

type Theme struct {
	Normal      Style
	Border      Style
	Focused     Style
	Active      Style
	Accessory   Style
	FillRune    rune
	BorderRunes BorderRuneSet
	ArrowRunes  ArrowRuneSet
	Overlay     bool // keep existing background
}

func (t Theme) String() string {
	return fmt.Sprintf(
		"{Normal=%v,Border=%v,Focused=%v,Active=%v,Accessory=%v,FillRune=%v,BorderRunes=%v,ArrowRunes=%v,Overlay=%v}",
		t.Normal,
		t.Border,
		t.Focused,
		t.Active,
		t.Accessory,
		t.FillRune,
		t.BorderRunes,
		t.ArrowRunes,
		t.Overlay,
	)
}
