package utils

func SolveIncrements(nChildren, nSpace int) (increment, remainder int) {
	remainder = nSpace % nChildren
	increment = (nSpace-remainder) / nChildren
	return
}

func SolveGaps(n, max int) (gaps []int) {
	// for n gaps, arrange max space
	for i := 0; i < n; i++ {
		gaps = append(gaps, 0)
	}
	front := false
	fw, bw := 0, n-1
	for SumInts(gaps) < max {
		if front {
			gaps[fw]++
			front = false
			fw++
			if fw > n-1 {
				fw = 0
			}
		} else {
			gaps[bw]++
			front = true
			bw--
			if bw < 1 {
				bw = n - 1
			}
		}
	}
	return
}
