package utils

func SolveSpaceAlloc(nChildren, nSpace, minSpacing int) (increment int, gaps []int) {
	numGaps := nChildren - 1
	availableSpace := nSpace - (minSpacing * numGaps)
	remainder := availableSpace % nChildren
	increment = (availableSpace - remainder) / nChildren
	extra := (minSpacing * numGaps) + remainder
	for i := 0; i <= numGaps; i++ {
		gaps = append(gaps, minSpacing)
	}
	front := true
	first, last := 0, numGaps
	fw, bw := 0, numGaps
	for SumInts(gaps) < extra {
		if front {
			gaps[fw]++
			front = false
			fw++
			if fw > last || fw == bw {
				fw = first
			}
		} else {
			gaps[bw]++
			front = true
			bw--
			if bw < first || bw == fw {
				bw = last
			}
		}
	}
	return
}
