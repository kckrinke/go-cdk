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
	"fmt"
)

var cdkMainScreen Screen

func MainScreen() Screen {
	return cdkMainScreen
}

func SetMainScreen(s Screen) {
	cdkMainScreen = s
}

var cdkDidBoot = false

func Bootstrap() error {
	if cdkDidBoot {
		return nil
	}
	return HandleSystemEvent(SYSTEM_EVENT_BOOT)
}

func MainInit() error {
	return MainInitWithTty("")
}

func MainInitWithTty(ttyPath string) error {
	Bootstrap()
	if cdkMainScreen == nil {
		s, e := NewScreen()
		if e != nil {
			return fmt.Errorf("error creating new screen: %v", e)
		}
		if e := s.Init(); e != nil {
			return fmt.Errorf("error initializing new screen: %v", e)
		}
		defStyle := StyleDefault.
			Background(ColorReset).
			Foreground(ColorReset)
		s.SetStyle(defStyle)
		s.EnableMouse()
		s.EnablePaste()
		s.Clear()
		SetMainScreen(s)
		HandleSystemEvent(SYSTEM_EVENT_INIT)
		return nil
	}
	return fmt.Errorf("non-operation: main already initialized")
}

func MainQuit() {
	HandleSystemEvent(SYSTEM_EVENT_SHUTDOWN)
	if cdkMainScreen != nil {
		cdkMainScreen.Close()
		cdkMainScreen = nil
	}
	Exit(0)
}
