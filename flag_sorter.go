package cdk

import (
	"strings"

	"github.com/urfave/cli/v2"
)

// CliFlagSorter is a slice of Flag.
type CliFlagSorter []cli.Flag

func (f CliFlagSorter) Len() int {
	return len(f)
}

func (f CliFlagSorter) Less(i, j int) bool {
	if len(f[j].Names()) == 0 {
		return false
	} else if len(f[i].Names()) == 0 {
		return true
	}
	i_is_ctk := strings.HasPrefix(f[i].Names()[0], "ctk-")
	j_is_ctk := strings.HasPrefix(f[j].Names()[0], "ctk-")
	if i_is_ctk {
		if j_is_ctk {
			return !LexicographicLess(f[i].Names()[0], f[j].Names()[0])
		}
		return false
	} else if j_is_ctk {
		return true
	}
	return LexicographicLess(f[i].Names()[0], f[j].Names()[0])
}

func (f CliFlagSorter) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
