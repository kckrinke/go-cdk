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
	"io/ioutil"
	"os"
)

var (
	_cdk__prev_exit            = cdkLogger.ExitFunc
	_cdk__fake_exited      int = -1
	_cdk__last_fake_logged     = ""
	_cdk__last_fake_exited     = -1
	_cdk__last_fake_error  error
)

func FakeExiting() {
	ResetFakeExited()
	_cdk__prev_exit = cdkLogger.ExitFunc
	cdkLogger.ExitFunc = func(code int) {
		_cdk__fake_exited = code
	}
}

func DidFakeExit() bool {
	return _cdk__fake_exited >= 0
}

func ResetFakeExited() {
	_cdk__fake_exited = -1
}

func RestoreExiting() {
	cdkLogger.ExitFunc = _cdk__prev_exit
	_cdk__prev_exit = nil
	_cdk__fake_exited = -1
}

func DoWithFakeIO(fn func() error) (string, bool, error) {
	real_stdout := os.Stdout
	real_stderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	FakeExiting()
	fn_err := fn()
	faked_it := DidFakeExit()
	RestoreExiting()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = real_stdout // restoring the real stdout
	os.Stderr = real_stderr // restoring the real stdout
	_cdk__last_fake_error = fn_err
	_cdk__last_fake_logged = string(out)
	_cdk__last_fake_exited = _cdk__fake_exited
	if fn_err != nil {
		return "", faked_it, fn_err
	}
	return string(out), faked_it, nil
}

func GetLastFakeIO() (string, int, error) {
	return _cdk__last_fake_logged, _cdk__last_fake_exited, _cdk__last_fake_error
}

type AppFn func(app App)
type DisplayManagerFn func(d DisplayManager)

func WithApp(initFn DisplayInitFn, action AppFn) func() {
	return func() {
		app := NewApp(
			"AppName", "AppUsage", "v0.0.0",
			"app-tag", "AppTitle",
			OffscreenDisplayTtyPath,
			initFn,
		)
		defer func() {
			if app.DisplayManager() != nil {
				app.DisplayManager().ReleaseDisplay()
				app.DisplayManager().Destroy()
			}
			app = nil
		}()
		action(app)
	}
}

func WithDisplayManager(action DisplayManagerFn) func() {
	return func() {
		d := NewDisplayManager("testing", OffscreenDisplayTtyPath)
		d.CaptureDisplay(OffscreenDisplayTtyPath)
		defer d.ReleaseDisplay()
		action(d)
	}
}

func TestingMakesNoContent(d DisplayManager) error {
	return nil
}

func TestingMakesActiveWindow(d DisplayManager) error {
	w := NewOffscreenWindow(d.GetTitle())
	d.SetActiveWindow(w)
	return nil
}
