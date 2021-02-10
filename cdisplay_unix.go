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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Thing of this is as CDK "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *cDisplay) engage() error {
	stdin := int(t.in.Fd())
	if _, err := term.MakeRaw(stdin); err != nil {
		return err
	}
	if w, h, err := term.GetSize(stdin); err == nil && w != 0 && h != 0 {
		t.cells.Resize(w, h)
		t.PostEvent(NewEventResize(w, h))
	}
	return nil
}

// disengage is used to release the terminal back to support from the caller.
// Think of this as CDK disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *cDisplay) disengage() {
	if t.in != nil {
		term.Restore(int(t.in.Fd()), t.saved)
	}
}

// initialize is used at application startup, and sets up the initial values
// including file descriptors used for terminals and saving the initial state
// so that it can be restored when the application terminates.
func (t *cDisplay) initialize() error {
	var err error
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return fmt.Errorf("display is not a terminal")
	}
	t.in = os.Stdin
	t.out = os.Stdout
	t.saved, err = term.GetState(fd)
	if err != nil {
		return err
	}
	signal.Notify(t.sigWinch, syscall.SIGWINCH)

	if err := t.engage(); err != nil {
		return err
	}
	return nil
}

// finalize is used to at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *cDisplay) finalize() {

	signal.Stop(t.sigWinch)

	<-t.inDoneQ

	t.disengage()
}

// getWinSize is called to obtain the terminal dimensions.
func (t *cDisplay) getWinSize() (int, int, error) {
	return term.GetSize(int(t.in.Fd()))
}

// Beep emits a beep to the terminal.
func (t *cDisplay) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}
