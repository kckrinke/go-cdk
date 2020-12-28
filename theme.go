package cdk

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
		Focused:     DefaultMonoCdkStyle.Dim(false),
		Active:      DefaultMonoCdkStyle.Dim(false).Bold(true),
		Border:      DefaultMonoCdkStyle.Dim(true),
		FillRune:    DefaultFillRune,
		BorderRunes: DefaultBorderRune,
		Overlay:     false,
	}
	DefaultColorTheme = Theme{
		Normal:      DefaultColorCdkStyle,
		Focused:     DefaultColorCdkStyle.Dim(false),
		Active:      DefaultColorCdkStyle.Dim(false).Bold(true),
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

type Theme struct {
	Normal      Style
	Focused     Style
	Active      Style
	Border      Style
	FillRune    rune
	BorderRunes BorderRune
	Overlay     bool // keep existing background
}
