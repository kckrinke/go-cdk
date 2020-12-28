package utils

import (
	"fmt"
	"regexp"
	"unicode"
)

func PadLeft(src, pad string, length int) string {
	for {
		if len(src) > length {
			return src[0 : length+1]
		}
		src = pad + src
	}
}

func PadRight(src, pad string, length int) string {
	for {
		if len(src) > length {
			return src[0 : length+1]
		}
		src += pad
	}
}

func CleanCRLF(s string) string {
	length := len(s)
	var last int
	for last = length - 1; last >= 0; last-- {
		if s[last] != '\r' && s[last] != '\n' {
			break
		}
	}
	return s[:last+1]
}

func NLSprintf(format string, argv ...interface{}) string {
	return CleanCRLF(fmt.Sprintf(format, argv...))
}

var _rxIsEmpty = regexp.MustCompile(`^\s*$`)

func IsEmpty(text string) bool {
	return len(text) == 0 || _rxIsEmpty.MatchString(text)
}

func HasSpace(text string) bool {
	for _, c := range text {
		if unicode.IsSpace(c) {
			return true
		}
	}
	return false
}
