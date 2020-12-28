package cdk

import (
	"strings"
)

type WordLine struct {
	words []*WordCell
}

func NewWordLine(line string, style Style) *WordLine {
	words := strings.Fields(line)
	wl := &WordLine{
		words: make([]*WordCell, len(words)),
	}
	for i, word := range words {
		wl.words[i] = NewWordCell(word, style)
	}
	return wl
}

func (w WordLine) LetterCount(spaces bool) int {
	c := 0
	for i, word := range w.words {
		if i != 0 && spaces {
			c += 1
		}
		c += word.Len()
	}
	return c
}

func (w WordLine) String() (s string) {
	s = ""
	for i, c := range w.words {
		if i > 0 {
			s += " "
		}
		s += c.String()
	}
	return
}
