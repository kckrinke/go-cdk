package utils

import (
	"unicode"
)

func LexicographicLess(i, j string) bool {
	iRunes := []rune(i)
	jRunes := []rune(j)

	lenShared := len(iRunes)
	if lenShared > len(jRunes) {
		lenShared = len(jRunes)
	}

	for index := 0; index < lenShared; index++ {
		ir := iRunes[index]
		jr := jRunes[index]

		if lir, ljr := unicode.ToLower(ir), unicode.ToLower(jr); lir != ljr {
			return lir < ljr
		}

		if ir != jr {
			return ir < jr
		}
	}

	return i < j
}

func RotateSlice(a []interface{}, rotation int) (rotated []interface{}) {
	if size := len(a); size > 0 {
		var tmp []interface{}
		for i := 0; i < rotation; i++ {
			tmp = a[1:size]
			tmp = append(tmp, a[0])
			rotated = tmp
		}
	}
	return
}
