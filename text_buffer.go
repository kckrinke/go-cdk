package cdk

import (
	"strings"
	"sync"
)

var (
	TAB_SPACE = "    "
)

type TextBuffer interface {
	Set(input string, style Style)
	LetterCount(spaces bool) int
	Draw(canvas *Canvas, singleLineMode bool, wordWrap WrapMode, justify Justification, align VerticalAlignment) EventFlag
}

type CTextBuffer struct {
	lines []*WordLine
	style Style

	sync.Mutex
}

func NewTextBuffer(input string, style Style) TextBuffer {
	tb := &CTextBuffer{
		lines: make([]*WordLine, 0),
	}
	tb.style = style
	tb.Set(input, style)
	return tb
}

func (b *CTextBuffer) Set(input string, style Style) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		b.lines = append(b.lines, NewWordLine(line, style))
	}
}

func (b *CTextBuffer) LetterCount(spaces bool) int {
	c := 0
	for i, word := range b.lines {
		if i != 0 && spaces {
			c += 1
		}
		c += word.LetterCount(spaces)
	}
	return c
}

func (b *CTextBuffer) getPosAtChar(atLine, n int) (lid, wid, cid int) {
	c := 0
	var line *WordLine
	var word *WordCell
	for lid, line = range b.lines {
		// TODO: sort this out better
		if lid != atLine {
			continue
		}
		for wid, word = range line.words {
			for cid, _ = range word.characters {
				if c == n {
					return
				}
				c++
			}
		}
	}
	return
}

func (b *CTextBuffer) Draw(canvas *Canvas, singleLine bool, wordWrap WrapMode, justify Justification, align VerticalAlignment) EventFlag {
	if !singleLine && canvas.size.H == 1 {
		singleLine = true
	}
	maxChars := canvas.size.W
	if singleLine {
		var atLine int
		switch align {
		case ALIGN_MIDDLE:
			atLine = (canvas.size.H / 2) - (len(b.lines) / 2)
		case ALIGN_BOTTOM:
			atLine = canvas.size.H - len(b.lines)
		case ALIGN_TOP:
		default:
			atLine = 0
		}
		spaces := b.LetterCount(true)
		if spaces > maxChars {
			b.truncateSingleLine(atLine, 0, canvas, wordWrap, maxChars)
			return EVENT_STOP
		}
		b.justifySingleLine(atLine, 0, canvas, wordWrap, maxChars, justify)
		return EVENT_STOP
	}
	b.wrapMultiLine(canvas, wordWrap, justify, align)
	return EVENT_STOP
}

func (b *CTextBuffer) wrapMultiLine(canvas *Canvas, wordWrap WrapMode, justify Justification, align VerticalAlignment) {
	var atLine int
	switch align {
	case ALIGN_MIDDLE:
		atLine = (canvas.size.H / 2) - (len(b.lines) / 2)
	case ALIGN_BOTTOM:
		atLine = canvas.size.H - len(b.lines)
	case ALIGN_TOP:
	default:
		atLine = 0
	}
	maxChars := canvas.size.W
	sorted := b.wrapSort(maxChars, wordWrap)
	origLines := b.lines
	b.lines = sorted
	for y, _ := range sorted {
		b.justifySingleLine(atLine+y, y, canvas, wordWrap, maxChars, justify)
	}
	b.lines = origLines
}

func (b *CTextBuffer) wrapSort(maxChars int, wordWrap WrapMode) []*WordLine {
	var sorted []*WordLine
	sorted = append(sorted, NewWordLine("", b.style))
	lid, wid := 0, 0
	for _, line := range b.lines {
		if line.LetterCount(true) < maxChars {
			if lid >= len(sorted) {
				sorted = append(sorted, NewWordLine("", b.style))
			}
			for _, word := range line.words {
				sorted[lid].words = append(sorted[lid].words, word)
			}
			lid++
			wid = 0
			continue
		}
		switch wordWrap {
		case WRAP_WORD:
			fallthrough
		case WRAP_WORD_CHAR:
			// attempt wrap word, if no spaces, fallback to char
			if strings.Contains(line.String(), " ") {
				for _, word := range line.words {
					if sorted[lid].LetterCount(true)+1+word.Len() < maxChars {
						sorted[lid].words = append(sorted[lid].words, word)
					} else {
						lid++
						sorted = append(sorted, NewWordLine("", b.style))
						sorted[lid].words = append(sorted[lid].words, word)
					}
				}
				break
			}
			fallthrough
		case WRAP_CHAR:
			for _, word := range line.words {
				if wid >= len(sorted[lid].words) {
					sorted[lid].words = append(sorted[lid].words, &WordCell{})
				}
				for _, char := range word.characters {
					if sorted[lid].LetterCount(true)+1+1 < maxChars {
						sorted[lid].words[wid].characters = append(sorted[lid].words[wid].characters, char)
					} else {
						lid++
						sorted = append(sorted, NewWordLine("", b.style))
						wid = 0
						if len(sorted[lid].words) < wid {
							sorted[lid].words = append(sorted[lid].words, NewWordCell("", b.style))
						}
						sorted[lid].words[wid].characters = append(sorted[lid].words[wid].characters, char)
					}
				}
				wid++
			}
		case WRAP_NONE:
			fallthrough
		default:
			// truncate
		truncate_loop:
			for wid, word := range line.words {
				if wid >= len(sorted[lid].words) {
					sorted[lid].words = append(sorted[lid].words, &WordCell{})
				}
				for _, char := range word.characters {
					if sorted[lid].LetterCount(true)+1+1 < maxChars {
						sorted[lid].words[wid].characters = append(sorted[lid].words[wid].characters, char)
					} else {
						lid++
						sorted = append(sorted, NewWordLine("", b.style))
						break truncate_loop
					}
				}
			}
		}
	}
	return sorted
}

func (b *CTextBuffer) justifySingleLine(atLine, forLine int, canvas *Canvas, wordWrap WrapMode, maxChars int, justify Justification) {
	if len(b.lines) < forLine || len(b.lines[forLine].words) == 0 {
		return
	}
	switch justify {
	case JUSTIFY_RIGHT:
		count := b.lines[forLine].LetterCount(true)
		delta := maxChars - count
		x := 0
		for x < maxChars-delta {
			for _, word := range b.lines[forLine].words {
				for _, char := range word.characters {
					canvas.SetRune(x+delta, atLine, char.Value(), char.Style())
					x++
				}
				// consume one blank space between words
				x++
			}
		}
	case JUSTIFY_CENTER:
		count := b.lines[forLine].LetterCount(true)
		half := count / 2
		halfway := canvas.size.W / 2
		delta := halfway - half
		x := 0
		for x < count {
			for _, word := range b.lines[forLine].words {
				for _, char := range word.characters {
					canvas.SetRune(x+delta, atLine, char.Value(), char.Style())
					x++
				}
				// consume one blank space between words
				x++
			}
		}
	case JUSTIFY_FILL:
		// word_count := len(b.lines[atLine].words)
		numgap := len(b.lines[forLine].words) - 1
		if numgap == 0 {
			return
		}
		var gaps []string
		for i := 0; i < numgap; i++ {
			gaps = append(gaps, " ")
		}
		// fmt.Printf("gaps: %v, words: %v, spaced:\"%v\"\n", gaps, words, spaced)
		s := b.joinGaps(b.lines[forLine], gaps)
		forward := true
		inner := 0
		outer := numgap - 1
		for {
			if forward {
				gaps[inner] += " "
				inner += 1
				if inner > numgap-1 {
					inner = 0
				}
				forward = false
			} else {
				gaps[outer] += " "
				outer -= 1
				if outer < 0 {
					outer = numgap - 1
				}
				forward = true
			}
			s = b.joinGaps(b.lines[forLine], gaps)
			if len(s) >= maxChars {
				// return s
				break
			}
		}
		// gaps ready for rendering
		x := 0
		for x < len(s) {
			for wid, word := range b.lines[forLine].words {
				for _, char := range word.characters {
					canvas.SetRune(x, atLine, char.Value(), char.Style())
					x++
				}
				if len(gaps) > wid {
					x += len(gaps[wid])
				}
			}
		}
	case JUSTIFY_LEFT:
		fallthrough
	default:
		x := 0
		for _, word := range b.lines[forLine].words {
			for _, char := range word.characters {
				canvas.SetRune(x, atLine, char.Value(), char.Style())
				x++
			}
			// consume one blank space between words
			x++
		}
	}
}

func (b *CTextBuffer) truncateSingleLine(atLine, forLine int, canvas *Canvas, wordWrap WrapMode, maxChars int) {
	switch wordWrap {
	case WRAP_WORD:
	case WRAP_WORD_CHAR:
		var lid, wid, cid int
		for i := maxChars; i > 0; i-- {
			lid, wid, cid = b.getPosAtChar(atLine, i)
			if b.lines[lid].words[wid].characters[cid].IsSpace() {
				break
			}
			lid, wid, cid = -1, -1, -1
		}
		if lid >= -1 {
		truncAtSpace:
			for i := 0; i < len(b.lines[lid].words); i++ {
				for j := 0; j < len(b.lines[lid].words[i].characters); j++ {
					char := b.lines[lid].words[i].characters[j]
					canvas.SetRune(atLine, j, char.Value(), char.Style())
					if i == wid && j == cid {
						break truncAtSpace
					}
				}
			}
			break
		}
		// the line has no spaces in it at all, must truncate
		// means WRAP_WORD_CHAR and WRAP_WORD are identical operations
		fallthrough
	case WRAP_CHAR:
		fallthrough
	case WRAP_NONE:
		fallthrough
	default:
		count := 0
	truncAtChar:
		for i := 0; i < len(b.lines[forLine].words); i++ {
			for j := 0; j < len(b.lines[forLine].words[i].characters); j++ {
				char := b.lines[forLine].words[i].characters[j]
				canvas.SetRune(j, atLine, char.Value(), char.Style())
				count++
				if count >= maxChars {
					break truncAtChar
				}
			}
		}
	}
	return
}

func (b *CTextBuffer) joinGaps(line *WordLine, gaps []string) string {
	output := ""
	last_idx := len(line.words) - 1
	for idx, word := range line.words {
		output += word.String()
		if idx < last_idx {
			output += gaps[idx]
		}
	}
	return output
}
