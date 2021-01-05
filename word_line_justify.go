package cdk

import (
	"github.com/kckrinke/go-cdk/utils"
)

func (w WordLine) applyTypographicJustifyFill(maxChars int, fillerStyle Style, input []*WordLine) (output []*WordLine) {
	// trim left/right space for each line, maximize gaps
	lid := 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		width := line.CharacterCount()
		gaps := make([]int, 0)
		for _, word := range line.Words() {
			if word.IsSpace() {
				gaps = append(gaps, 1)
			}
		}
		widthMinusGaps := width - len(gaps)
		gaps = utils.DistInts(maxChars-widthMinusGaps, gaps)
		gid := 0
		for _, word := range line.Words() {
			if word.IsSpace() {
				wc := NewEmptyWordCell()
				for i := 0; i < gaps[gid]; i++ {
					wc.AppendRune(' ', fillerStyle)
				}
				gid++
				output[lid].AppendWordCell(wc)
			} else {
				output[lid].AppendWordCell(word)
			}
		}
		lid++
	}
	return
}

func (w WordLine) applyTypographicJustifyCenter(maxChars int, fillerStyle Style, input []*WordLine) (output []*WordLine) {
	// trim left space for each line
	wid, lid := 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		width := line.CharacterCount()
		halfWidth := width / 2
		halfWay := maxChars / 2
		delta := halfWay - halfWidth
		if delta > 0 {
			for i := 0; i < delta; i++ {
				output[lid].AppendWordCell(NewWordCell(" ", fillerStyle))
			}
		}
		for _, word := range line.Words() {
			output[lid].AppendWordCell(word)
			wid++
		}
		lid++
	}
	return
}

func (w WordLine) applyTypographicJustifyRight(maxChars int, fillerStyle Style, input []*WordLine) (output []*WordLine) {
	// trim left space for each line, assume no line needs wrapping or truncation
	wid, lid := 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		charCount := line.CharacterCount()
		delta := maxChars - charCount
		if delta > 0 {
			for i := 0; i < delta; i++ {
				output[lid].AppendWordCell(NewWordCell(" ", fillerStyle))
			}
		}
		for _, word := range line.Words() {
			output[lid].AppendWordCell(word)
			wid++
		}
		lid++
	}
	return
}

func (w WordLine) applyTypographicJustifyLeft(maxChars int, fillerStyle Style, input []*WordLine) (output []*WordLine) {
	// trim left space for each line
	wid, lid := 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		start := true
		for _, word := range line.Words() {
			if start {
				if word.IsSpace() {
					continue
				}
				start = false
			}
			output[lid].AppendWordCell(word)
			wid++
		}
		lid++
	}
	return
}
