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

package cdk

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/kckrinke/go-cdk/utils"
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
			return !utils.LexicographicLess(f[i].Names()[0], f[j].Names()[0])
		}
		return false
	} else if j_is_ctk {
		return true
	}
	return utils.LexicographicLess(f[i].Names()[0], f[j].Names()[0])
}

func (f CliFlagSorter) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
