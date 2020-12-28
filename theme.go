package cdk

import (
	"fmt"
)

var (
	DefaultFillRune      rune  = ' '
	DefaultMonoCdkStyle  Style = StyleDefault.Dim(true)
	DefaultColorCdkStyle Style = DefaultMonoCdkStyle.Foreground(ColorWhite).Background(ColorNavy)
	DefaultCdkStyle      Style = DefaultMonoCdkStyle
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
	DefaultMonoTheme = Theme{
		Normal:      DefaultMonoCdkStyle,
		Border:      DefaultMonoCdkStyle.Dim(true),
		FillRune:    DefaultFillRune,
		BorderRunes: DefaultBorderRune,
		Overlay:     false,
	}
	DefaultColorTheme = Theme{
		Normal:      DefaultColorCdkStyle,
		Border:      DefaultColorCdkStyle.Dim(true),
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
	FillRune    rune
	BorderRunes BorderRune
	Overlay     bool // keep existing background
}

func (t Theme) String() string {
	return fmt.Sprintf(
		"{Normal=%v,Border=%v,FillRune=%v,BorderRunes:%v,Overlay=%v}",
		t.Normal,
		t.Border,
		t.FillRune,
		t.BorderRunes,
		t.Overlay,
	)
}
