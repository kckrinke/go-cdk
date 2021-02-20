// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

// Copyright 2021 The TCell Authors
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
	// "fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kckrinke/go-term"
)

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Thing of this is as CDK "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *cDisplay) engage() error {
	if err := term.RawMode(t.term); err != nil {
		return err
	}
	if w, h, err := t.term.Winsz(); err == nil && w > 0 && h > 0 {
		t.cells.Resize(w, h)
		_ = t.PostEvent(NewEventResize(w, h))
	}
	return nil
}

// disengage is used to release the terminal back to support from the caller.
// Think of this as CDK disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *cDisplay) disengage() {
	if err := t.term.Restore(); err != nil {
		ErrorF("error restoring terminal: %v", err)
	}
}

// initialize is used at application startup, and sets up the initial values
// including file descriptors used for terminals and saving the initial state
// so that it can be restored when the application terminates.
func (t *cDisplay) initialize() error {
	var err error
	if t.term, err = term.Open("/dev/tty"); err != nil {
		return err
	}
	if err = term.RawMode(t.term); err != nil {
		return err
	}
	signal.Notify(t.sigWinch, syscall.SIGWINCH)
	if w, h, e := t.getWinSize(); e == nil && w != 0 && h != 0 {
		t.cells.Resize(w, h)
		_ = t.PostEvent(NewEventResize(w, h))
	}
	return nil
}

// finalize is used to at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *cDisplay) finalize() {
	signal.Stop(t.sigWinch)
	<-t.inDoneQ
	if t.term != nil {
		if err := term.CBreakMode(t.term); err != nil {
			ErrorF("error setting CBreakMode: %v", err)
		}
		if err := t.term.Restore(); err != nil {
			ErrorF("error restoring terminal: %v", err)
		}
		if err := t.term.Close(); err != nil {
			ErrorF("error closing terminal: %v", err)
		}
	}
}

// getWinSize is called to obtain the terminal dimensions.
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

// Beep emits a beep to the terminal.
func (t *cDisplay) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}
