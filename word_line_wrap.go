package cdk

// wrap the input lines on the nearest word to maxChars
func (w *CWordLine) applyTypographicWrapWord(maxChars int, input []WordLine) (output []WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		if line.CharacterCount() > maxChars {
			if !line.HasSpace() {
				// nothing to break on, truncate on maxChars
				for _, word := range line.Words() {
					if wid >= output[lid].Len() {
						output[lid].AppendWordCell(NewEmptyWordCell())
					}
					for _, c := range word.Characters() {
						if cid >= maxChars {
							lid = len(output) // don't append trailing NEWLs
							break
						}
						output[lid].AppendWordRune(wid, c.Value(), c.Style())
						cid++
					}
					wid++
				}
				continue
			}
		}
		for _, word := range line.Words() {
			wordLen := word.Len()
			if word.IsSpace() && wordLen > 1 {
				wordLen = 1
			}
			if cid+wordLen > maxChars {
				output = append(output, NewEmptyWordLine())
				lid = len(output) - 1
				wid = -1
				cid = 0
				if !word.IsSpace() {
					output[lid].AppendWordCell(word)
					wid = output[lid].Len() - 1
					cid += word.Len()
				}
			} else if word.IsSpace() && cid+wordLen+1 > maxChars {
				// continue
			} else {
				if word.IsSpace() {
					if c := word.GetCharacter(0); c != nil {
						wc := NewEmptyWordCell()
						wc.AppendRune(c.Value(), c.Style())
						output[lid].AppendWordCell(wc)
						cid += wc.Len()
					}
				} else {
					output[lid].AppendWordCell(word)
					cid += word.Len()
				}
			}
			wid++
		}
		lid++
		cid = 0
		wid = 0
	}
	return
}

// wrap the input lines on the nearest word to maxChars if the line has space,
// else, truncate at maxChars
func (w *CWordLine) applyTypographicWrapWordChar(maxChars int, input []WordLine) (output []WordLine) {
	for lid, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		if line.CharacterCount() > maxChars {
			if !line.HasSpace() {
				wrapped := w.applyTypographicWrapChar(maxChars, []WordLine{line})
				for wLid, wLine := range wrapped {
					id := lid + wLid
					if id >= len(output) {
						output = append(output, NewEmptyWordLine())
					}
					for _, wWord := range wLine.Words() {
						output[id].AppendWordCell(wWord)
					}
				}
				continue
			}
		}
		wrapped := w.applyTypographicWrapWord(maxChars, []WordLine{line})
		for wLid, wLine := range wrapped {
			id := lid + wLid
			if id >= len(output) {
				output = append(output, NewEmptyWordLine())
			}
			for _, wWord := range wLine.Words() {
				output[id].AppendWordCell(wWord)
			}
		}
	}
	return
}

// wrap the input lines on the nearest character to maxChars
func (w *CWordLine) applyTypographicWrapChar(maxChars int, input []WordLine) (output []WordLine) {
	cid, wid, lid := 0, 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		for _, word := range line.Words() {
			if word.IsSpace() {
				if cid+1 > maxChars {
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

// truncate the input lines on the nearest character to maxChars
func (w *CWordLine) applyTypographicWrapNone(maxChars int, input []WordLine) (output []WordLine) {
	cid, lid := 0, 0
	for _, line := range input {
		if lid >= len(output) {
			output = append(output, NewEmptyWordLine())
		}
		for _, word := range line.Words() {
			if word.IsSpace() {
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
						if cid+c.Width() > maxChars {
							break
						}
						wc.AppendRune(c.Value(), c.Style())
						cid += c.Width()
					}
					if wc.Len() > 0 {
						output[lid].AppendWordCell(wc)
						break
					}
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