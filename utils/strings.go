package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// basic word wrapping utility method
func BasicWordWrap(input string, maxChars int) (output string) {
	var words []string
	single := strings.Replace(input, "\n", " ", -1)
	if words = strings.Fields(single); len(words) == 0 {
		return input
	}
	output = ""
	count := 0
	for idx, word := range words {
		if count + 1 + len(word) >= maxChars {
			output += "\n"
			count = 0
		} else if idx > 0 {
			output += " "
		}
		output += word
		count += 1 + len(word)
	}
	return
}

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

func IsTrue(text string) bool {
	switch strings.ToLower(text) {
	case "1":
		fallthrough
	case "on":
		fallthrough
	case "yes":
		fallthrough
	case "true":
		return true
	}
	return false
}

func IsFalse(text string) bool {
	switch strings.ToLower(text) {
	case "0":
		fallthrough
	case "off":
		fallthrough
	case "no":
		fallthrough
	case "false":
		return true
	}
	return false
}

func IsUrl(str string) (isUrl bool) {
	if u, err := url.Parse(str); err == nil && u.Scheme != "" {
		return true
	}
	return
}

func StringSliceHasValue(slice []string, value string) (has bool) {
	for _, str := range slice {
		if str == value {
			return true
		}
	}
	return
}
