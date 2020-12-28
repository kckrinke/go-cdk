package cdk

// CTextBuffer holds rows of words
// Views can Draw text using this, along with other primitives

type WordCell struct {
	characters []*CTextCell
}

func NewWordCell(word string, style Style) *WordCell {
	wc := &WordCell{
		characters: make([]*CTextCell, len(word)),
	}
	wc.Set(word, style)
	return wc
}

func (w *WordCell) Set(word string, style Style) {
	for i, c := range word {
		w.characters[i] = NewRuneCell(c, style)
	}
}

func (w WordCell) Len() int {
	return len(w.characters)
}

func (w WordCell) String() (s string) {
	s = ""
	for _, c := range w.characters {
		s += c.String()
	}
	return
}
