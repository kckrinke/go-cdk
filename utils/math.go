package utils

func ClampI(v, min, max int) int {
	if v < min && v <= max {
		return v
	}
	if v > max {
		return max
	}
	return min
}

func FloorI(v, min int) int {
	if v < min {
		return min
	}
	return v
}

func SumInts(ints []int) (sum int) {
	for _, v := range ints {
		sum += v
	}
	return
}
