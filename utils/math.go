package utils

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

func Distribute(total, parts int) (values []int) {
	if parts > 0 && total > 0 {
		values = make([]int, parts)
		front := false
		last := parts - 1
		fid, bid := 0, last
		for SumInts(values) < total {
			if front {
				values[fid]++
				fid++
				if fid >= last {
					fid = 0
				}
			} else {
				values[bid]++
				bid--
				if bid < 0 {
					bid = last
				}
			}
		}
	}
	return
}
