// +build windows

// Copyright 2015 The TCell Authors
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

// On Windows we don't have support for termios.  We probably could, and
// may should, in a cygwin type environment.  Its not clear how to make
// this all work nicely with both cygwin and Windows console, so we
// decline to do so here.

func (t *cDisplay) termioInit(ttyPath string) error {
	return ErrNoDisplay
}

func (t *cDisplay) termioClose() {
	return
}

func (t *cDisplay) getWinSize() (int, int, error) {
	return 0, 0, ErrNoDisplay
}

func (t *cDisplay) getCharset() string {
	return "UTF-16LE"
}

func (t *cDisplay) Beep() error {
	return ErrNoDisplay
}

type termiosPrivate struct{}