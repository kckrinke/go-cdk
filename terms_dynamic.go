// +build !tcell_minimal,!nacl,!js,!zos,!plan9,!windows,!android

// Copyright 2021 The CDK Authors
// Copyright 2019 The TCell Authors
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
	// This imports a dynamic version of the terminal database, which
	// is built using infocmp.  This relies on a working installation
	// of infocmp (typically supplied with ncurses).  We only do this
	// for systems likely to have that -- i.e. UNIX based hosts.  We
	// also don't support Android here, because you really don't want
	// to run external programs there.  Generally the android terminals
	// will be automatically included anyway.
	"github.com/kckrinke/go-terminfo"
	"github.com/kckrinke/go-terminfo/dynamic"
)

func loadDynamicTerminfo(term string) (*terminfo.Terminfo, error) {
	ti, _, e := dynamic.LoadTerminfo(term)
	if e != nil {
		return nil, e
	}
	return ti, nil
}
