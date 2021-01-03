package utils

import (
	"fmt"
)

// Returns the `value` given unless it's smaller than `min` or greater than
// `max`. If it's less than `min`, `min` is returned and if it's greater than
// `max` it returns max.
func ClampI(value, min, max int) int {
	if value >= min && value <= max {
		return value
	}
	if value > max {
		return max
	}
	return min
}

// Returns the `value` given unless it's less than `min`, in which case it
// returns `min`.
func FloorI(v, min int) int {
	if v < min {
		return min
	}
	return v
}

// Add the given list of integers up and return the result.
func SumInts(ints []int) (sum int) {
	for _, v := range ints {
		sum += v
	}
	return
}

// Round the given floating point number to the nearest larger integer and
// return that as an integer.
func CeilF2I(v float64) int {
	delta := v - float64(int(v))
	if delta > 0 {
		return int(v) + 1
	}
	return int(v)
}

func DistInts(max int, in []int) (out []int) {
	out = in
	front := true
	first, last := 0, len(out)-1
	fw, bw := 0, last
	for SumInts(out) < max {
		if front {
			out[fw]++
			front = false
			fw++
			if fw > last || fw == bw {
				fw = first
			}
		} else {
			out[bw]++
			front = true
			bw--
			if bw < first || bw == fw {
				bw = last
			}
		}
	}
	return
}

func SolveSpaceAlloc(nChildren, nSpace, minSpacing int) (increment int, gaps []int) {
	numGaps := nChildren - 1
	totalMinSpacing := minSpacing * numGaps
	availableSpace := nSpace - totalMinSpacing
	remainder := availableSpace % nChildren
	increment = (availableSpace - remainder) / nChildren
	extra := totalMinSpacing + remainder
	for i := 0; i < numGaps; i++ {
		gaps = append(gaps, minSpacing)
	}
	gaps = DistInts(extra, gaps)
	return
}

func Distribute(total, available, parts, nChildren, spacing int) (values, gaps []int, err error) {
	numGaps := nChildren - 1
	if numGaps > 0 {
		gaps = make([]int, numGaps)
		for i := 0; i < numGaps; i++ {
			gaps[i] = spacing
		}
	} else {
		gaps = make([]int, 0)
	}
	available -= SumInts(gaps)
	values = make([]int, parts)
	if parts > 0 {
		values = DistInts(available, values)
	}
	totalValues := SumInts(values)
	totalGaps := SumInts(gaps)
	totalDist := totalValues + totalGaps
	if totalDist > total {
		err = fmt.Errorf("totalDist[%d] > total[%d]", totalDist, total)
	} else if totalDist < total {
		delta := total - totalDist
		values = DistInts(SumInts(values) + delta, values)
	}
	return
}
