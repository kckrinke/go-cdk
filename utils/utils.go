// Copyright 2020 The CDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// see: github.com/urfave/cli/v2
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

func SumInts(ints []int) (sum int) {
	for _, v := range ints {
		sum += v
	}
	return
}

func SolveGaps(n, max int) (gaps []int) {
	// for n gaps, arrange max space
	for i := 0; i < n; i++ {
		gaps = append(gaps, 0)
	}
	if n < max {
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
	}
	return
}
