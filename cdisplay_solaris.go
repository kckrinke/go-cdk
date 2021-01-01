// +build solaris illumos

// Copyright 2020 The TCell Authors
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
	"os/signal"
)

func (t *cDisplay) termioInit(ttyPath string) error {
	if ttyPath == "" {
		ttyPath = "/dev/tty"
	}
	var err error
	if t.term, err = term.Open(ttyPath); err != nil {
		return err
	}
	term.RawMode(t.term)
	signal.Notify(t.sigwinch, syscall.SIGWINCH)
	if w, h, e := t.getWinSize(); e == nil && w != 0 && h != 0 {
		t.cells.Resize(w, h)
	}
	return nil
}

func (t *cDisplay) termioClose() {

	signal.Stop(t.sigwinch)

	<-t.indoneq

	if t.term != nil {
		term.CBreakMode(t.term)
		t.term.Restore()
		t.term.Close()
	}
}

func (t *cDisplay) getWinSize() (w, h int, err error) {
	w, h, err = t.term.Winsz()
	if err != nil {
		w, h = -1, -1
		return
	}
	if w == 0 {
		colsEnv := os.Getenv("COLUMNS")
		if colsEnv != "" {
			if w, err = strconv.Atoi(colsEnv); err != nil {
				w, h = -1, -1
				return
			}
		} else {
			w = t.ti.Columns
		}
	}
	if h == 0 {
		rowsEnv := os.Getenv("LINES")
		if rowsEnv != "" {
			if h, err = strconv.Atoi(rowsEnv); err != nil {
				w, h = -1, -1
				return
			}
		} else {
			h = t.ti.Lines
		}
	}
	return
}

func (t *cDisplay) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}
