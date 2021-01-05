package cdk

func (w WordLine) applyTypographicWrapWord(maxChars int, input []*WordLine) (output []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		if line.Len() > maxChars {
			if !line.HasSpace() {
				// nothing to break on, truncate on maxChars
			truncateWrapWord:
				for _, word := range line.Words() {
					if wid >= output[lid].Len() {
						output[lid].AppendWordCell(NewEmptyWordCell())
					}
					for _, c := range word.Characters() {
						if cid > maxChars {
							lid = len(output) // don't append trailing NEWLs
							break truncateWrapWord
						}
						output[lid].words[wid].AppendRune(c.Value(), c.Style())
						cid++
					}
					wid++
				}
				continue
			}
		}
		for _, word := range line.Words() {
			if wid >= output[lid].Len() {
				output[lid].AppendWordCell(NewEmptyWordCell())
			}
			wordLen := word.Len()
			if word.IsSpace() {
				wordLen = 1
			}
			if cid+wordLen >= maxChars {
				output = append(output, NewEmptyWordLine())
				lid = len(output) - 1
				wid = 0
				if !word.IsSpace() {
					output[lid].AppendWordCell(word)
					wid = output[lid].Len() - 1
				}
				continue
			}
			if word.IsSpace() {
				c := word.GetCharacter(0)
				wc := NewEmptyWordCell()
				wc.AppendRune(c.Value(), c.Style())
				output[lid].AppendWordCell(wc)
				cid += wc.Len()
			} else {
				output[lid].AppendWordCell(word)
				cid += word.Len()
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}

func (w WordLine) applyTypographicWrapWordChar(maxChars int, input []*WordLine) (output []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		if line.Len() > maxChars {
			if !line.HasSpace() {
				// nothing to break on, wrap on maxChars
				for _, word := range line.Words() {
					if wid >= output[lid].Len() {
						output[lid].AppendWordCell(NewEmptyWordCell())
					}
					for _, c := range word.Characters() {
						if cid > maxChars {
							output = append(output, NewEmptyWordLine())
							lid = len(output) - 1
							output[lid].AppendWordCell(NewEmptyWordCell())
							wid = 0
							cid = 0
						}
						output[lid].words[wid].AppendRune(c.Value(), c.Style())
						cid++
					}
					wid++
				}
				lid++
				wid = 0
				cid = 0
				continue
			}
		}
		for _, word := range line.Words() {
			if wid >= output[lid].Len() {
				output[lid].AppendWordCell(NewEmptyWordCell())
			}
			wordLen := word.Len()
			if word.IsSpace() {
				wordLen = 1
			}
			if cid+wordLen >= maxChars {
				output = append(output, NewEmptyWordLine())
				lid = len(output) - 1
				wid = 0
				if !word.IsSpace() {
					output[lid].AppendWordCell(word)
					wid = output[lid].Len() - 1
				}
				continue
			}
			if word.IsSpace() {
				c := word.GetCharacter(0)
				wc := NewEmptyWordCell()
				wc.AppendRune(c.Value(), c.Style())
				output[lid].AppendWordCell(wc)
				cid += wc.Len()
			} else {
				output[lid].AppendWordCell(word)
				cid += word.Len()
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}

func (w WordLine) applyTypographicWrapChar(maxChars int, input []*WordLine) (output []*WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		for _, word := range line.Words() {
			if word.IsSpace() {
				if cid > maxChars {
					output = append(output, NewEmptyWordLine())
					lid = len(output) - 1
					wid = 0
					cid = 0
				}
				if c := word.GetCharacter(0); c != nil {
					wc := NewEmptyWordCell()
					wc.AppendRune(c.Value(), c.Style())
					output[lid].AppendWordCell(wc)
					cid += wc.Len()
				}
			} else {
				if cid+word.Len() > maxChars {
					firstHalf, secondHalf := NewEmptyWordCell(), NewEmptyWordCell()
					for _, c := range word.Characters() {
						if cid < maxChars {
							firstHalf.AppendRune(c.Value(), c.Style())
						} else {
							secondHalf.AppendRune(c.Value(), c.Style())
						}
						cid++
					}
					output[lid].AppendWordCell(firstHalf)
					output = append(output, NewEmptyWordLine())
					lid = len(output) - 1
					output[lid].AppendWordCell(secondHalf)
					cid = secondHalf.Len()
					wid = 0
				} else {
					output[lid].AppendWordCell(word)
					cid += word.Len()
				}
			}
			wid++
		}
		lid++
		cid = 0
	}
	return
}

func (w WordLine) applyTypographicWrapNone(maxChars int, input []*WordLine) (output []*WordLine) {
	cid, lid := 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		for _, word := range line.Words() {
			if word.IsSpace() {
				if cid+1 > maxChars {
					lid = len(output) - 1
					cid = 0
					break
				}
				if c := word.GetCharacter(0); c != nil {
					wc := NewEmptyWordCell()
					wc.AppendRune(c.Value(), c.Style())
					output[lid].AppendWordCell(wc)
					cid += wc.Len()
				}
			} else {
				if cid+word.Len() > maxChars {
					wc := NewEmptyWordCell()
					for _, c := range word.Characters() {
						if cid > maxChars {
							output = append(output, NewEmptyWordLine())
							lid = len(output) - 1
							cid = 0
							break
						}
						wc.AppendRune(c.Value(), c.Style())
						cid++
					}
					output[lid].AppendWordCell(wc)
					cid += wc.Len()
				} else {
					output[lid].AppendWordCell(word)
					cid += word.Len()
				}
			}
		}
		lid++
		cid = 0
	}
	return
}
