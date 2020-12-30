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

// display is really a wrapper around Screen
// and Simulation screens

// basically a wrapper around Screen()
// manages one or more windows backed by viewports
// viewports manage the allocation of space
// drawables within viewports render the space

var (
	DisplayCallQueueCapacity = 16
)

const (
	TypeDisplay       CTypeTag = "display"
	SignalDisplayInit Signal   = "display-init"
	SignalInterrupt   Signal   = "ctrl-c"
	SignalEvent       Signal   = "event"
	SignalEventError  Signal   = "event-error"
	SignalEventKey    Signal   = "event-key"
	SignalEventMouse  Signal   = "event-mouse"
	SignalEventResize Signal   = "event-resize"
)

func init() {
	TypesManager.AddType(TypeDisplay)
}

type DisplayCallbackFn = func(d Display) error

type Display interface {
	Object

	GetTitle() string
	SetTitle(title string)

	Screen() Screen
	ScreenCaptured() bool
	CaptureScreen(ttyPath string)
	ReleaseScreen()

	CaptureCtrlC()
	ReleaseCtrlC()

	DefaultTheme() Theme

	ActiveWindow() Window
	SetActiveWindow(w Window)
	AddWindow(w Window) int
	GetWindows() []Window

	App() *CApp
	ProcessEvent(evt Event) EventFlag
	DrawScreen() EventFlag

	RequestDraw()
	RequestShow()
	RequestSync()
	RequestQuit()
	PostEvent(evt Event) error
	AsyncCall(fn DisplayCallbackFn) error
	AwaitCall(fn DisplayCallbackFn) error

	IsRunning() bool
	Run() error
}

// Basic display type
type CDisplay struct {
	CObject

	title string

	captureCtrlC bool

	active  int
	windows []Window

	app      *CApp
	ttyPath  string
	screen   Screen
	captured bool

	running  bool
	done     chan bool
	queue    chan DisplayCallbackFn
	events   chan Event
	process  chan Event
	requests chan ScreenStateReq
}

func NewDisplay(title string, ttyPath string) *CDisplay {
	d := new(CDisplay)
	d.title = title
	d.ttyPath = ttyPath
	d.Init()
	return d
}

// Initialization
func (d *CDisplay) Init() (already bool) {
	if d.InitTypeItem(TypeDisplay) {
		return true
	}
	d.CObject.Init()

	d.captured = false
	d.running = false
	d.done = make(chan bool)
	d.queue = make(chan DisplayCallbackFn, DisplayCallQueueCapacity)
	d.events = make(chan Event, DisplayCallQueueCapacity)
	d.process = make(chan Event, DisplayCallQueueCapacity)
	d.requests = make(chan ScreenStateReq, DisplayCallQueueCapacity)

	d.windows = []Window{}
	d.active = -1
	d.SetTheme(DefaultColorTheme)
	d.Emit(SignalDisplayInit, d)
	return false
}

func (d *CDisplay) GetTitle() string {
	return d.title
}

func (d *CDisplay) SetTitle(title string) {
	d.title = title
}

func (d *CDisplay) Screen() Screen {
	return d.screen
}

func (d *CDisplay) ScreenCaptured() bool {
	return d.screen != nil && d.captured
}

func (d *CDisplay) CaptureScreen(ttyPath string) {
	d.Lock()
	defer d.Unlock()
	var err error
	d.screen, err = NewScreen()
	if err != nil {
		Fatalf("error getting new screen: %v", err)
	}
	if err = d.screen.InitWithTty(ttyPath); err != nil {
		Fatalf("error initializing new screen: %v", err)
	}
	defStyle := StyleDefault.
		Background(ColorReset).
		Foreground(ColorReset)
	d.screen.SetStyle(defStyle)
	d.screen.EnableMouse()
	d.screen.EnablePaste()
	d.screen.Clear()
	d.SetTheme(d.DefaultTheme())
	d.captured = true
}

func (d *CDisplay) ReleaseScreen() {
	d.Lock()
	defer d.Unlock()
	if d.screen != nil {
		d.screen.Close()
		d.screen = nil
	}
	d.captured = false
}

func (d *CDisplay) CaptureCtrlC() {
	d.Lock()
	defer d.Unlock()
	d.captureCtrlC = true
}

func (d *CDisplay) ReleaseCtrlC() {
	d.Lock()
	defer d.Unlock()
	d.captureCtrlC = false
}

func (d *CDisplay) DefaultTheme() Theme {
	if d.screen != nil && d.screen.Colors() > 0 {
		return DefaultColorTheme
	}
	return DefaultMonoTheme
}

func (d *CDisplay) ActiveWindow() Window {
	if len(d.windows) > d.active && d.active >= 0 {
		return d.windows[d.active]
	}
	if len(d.windows) == 0 {
		return nil
	}
	d.active = 0
	return d.windows[0]
}

func (d *CDisplay) SetActiveWindow(w Window) {
	d.Lock()
	var id int = -1
	var window Window
	for id, window = range d.windows {
		if window == w {
			break
		}
	}
	if id > -1 {
		d.active = id
		d.Unlock()
		return
	}
	d.Unlock()
	d.active = d.AddWindow(w)
}

func (d *CDisplay) AddWindow(w Window) int {
	d.Lock()
	defer d.Unlock()
	var id int = -1
	var window Window
	for id, window = range d.windows {
		if window == w {
			break
		}
	}
	if id > -1 {
		d.LogError("display has window already: %v", w)
		return id
	}
	d.windows = append(d.windows, w)
	w.SetDisplay(d)
	return len(d.windows) - 1
}

func (d *CDisplay) GetWindows() []Window {
	return d.windows
}

func (d *CDisplay) App() *CApp {
	return d.app
}

func (d *CDisplay) ProcessEvent(evt Event) EventFlag {
	switch e := evt.(type) {
	case *EventError:
		d.LogErr(e)
		if w := d.ActiveWindow(); w != nil {
			if f := w.ProcessEvent(evt); f == EVENT_STOP {
				return EVENT_STOP
			}
		}
		return d.Emit(SignalEventError, d, e)
	case *EventKey:
		if d.captureCtrlC {
			switch e.Key() {
			case KeyCtrlC:
				d.LogTrace("display captured CtrlC")
				if f := d.Emit(SignalInterrupt, d); f == EVENT_STOP {
					return EVENT_STOP
				}
				d.RequestQuit()
			}
		}
		if w := d.ActiveWindow(); w != nil {
			if f := w.ProcessEvent(evt); f == EVENT_STOP {
				return EVENT_STOP
			}
		}
		return d.Emit(SignalEventKey, d, e)
	case *EventMouse:
		if w := d.ActiveWindow(); w != nil {
			if f := w.ProcessEvent(evt); f == EVENT_STOP {
				return EVENT_STOP
			}
		}
		return d.Emit(SignalEventMouse, d, e)
	case *EventResize:
		if w := d.ActiveWindow(); w != nil {
			if f := w.ProcessEvent(evt); f == EVENT_STOP {
				return EVENT_STOP
			}
		}
		return d.Emit(SignalEventResize, d, e)
	}
	if w := d.ActiveWindow(); w != nil {
		if f := w.ProcessEvent(evt); f == EVENT_STOP {
			return EVENT_STOP
		}
	}
	return d.Emit(SignalEvent, d, evt)
}

func (d *CDisplay) DrawScreen() EventFlag {
	d.Lock()
	defer d.Unlock()
	if d.screen == nil {
		d.LogError("display missing screen")
		return EVENT_PASS
	}
	var window Window
	if window = d.ActiveWindow(); window == nil {
		d.LogDebug("cannot draw the screen, display missing a window")
		return EVENT_PASS
	}
	w, h := d.screen.Size()
	canvas := NewCanvas(Point2I{0, 0}, Rectangle{w, h}, d.GetTheme())
	if f := window.Draw(canvas); f == EVENT_STOP {
		canvas.Render(d.screen)
		return EVENT_STOP
	}
	return EVENT_PASS
}

func (d *CDisplay) RequestDraw() {
	d.requests <- DrawRequest
}

func (d *CDisplay) RequestShow() {
	d.requests <- ShowRequest
}

func (d *CDisplay) RequestSync() {
	d.requests <- SyncRequest
}

func (d *CDisplay) RequestQuit() {
	d.requests <- QuitRequest
}

func (d *CDisplay) AsyncCall(fn DisplayCallbackFn) error {
	if !d.running {
		return fmt.Errorf("application not running")
	}
	d.queue <- fn
	return nil
}

func (d *CDisplay) AwaitCall(fn DisplayCallbackFn) error {
	if !d.running {
		return fmt.Errorf("application not running")
	}
	var err error
	done := make(chan bool)
	d.queue <- func(d Display) error {
		err = fn(d)
		done <- true
		return nil
	}
	<-done
	return err
}

func (d *CDisplay) PostEvent(evt Event) error {
	if !d.running {
		return fmt.Errorf("application not running")
	}
	d.events <- evt
	return nil
}

func (d *CDisplay) runPollEventThread() {
	for d.running {
		d.process <- d.screen.PollEvent()
	}
	d.done <- true
}

func (d *CDisplay) runProcessEventThread() {
	for d.running {
		if evt := <-d.process; evt != nil {
			if f := d.ProcessEvent(evt); f == EVENT_STOP {
				d.RequestDraw()
				d.RequestShow()
			}
		}
	}
}
func (d *CDisplay) runRequestThread() {
	if d.running {
		d.RequestDraw()
		d.RequestSync()
	}
	for d.running {
		switch <-d.requests {
		case DrawRequest:
			if d.screen != nil {
				d.DrawScreen()
			}
		case ShowRequest:
			if d.screen != nil {
				d.screen.Show()
			}
		case SyncRequest:
			if d.screen != nil {
				d.screen.Sync()
			}
		case QuitRequest:
			d.running = false
			d.process <- nil
			d.done <- true
		}
	}
}
func (d *CDisplay) Run() error {
	d.CaptureScreen(d.ttyPath)
	d.running = true
	go d.runPollEventThread()
	go d.runProcessEventThread()
	go d.runRequestThread()
	defer func() {
		if p := recover(); p != nil {
			d.ReleaseScreen()
			panic(p)
		}
	}()
	d.events <- NewEventResize(d.screen.Size())
	for {
		select {
		case fn := <-d.queue:
			if err := fn(d); err != nil {
				return err
			}
		case evt := <-d.events:
			d.screen.PostEvent(evt)
		case <-d.done:
			d.ReleaseScreen()
			return nil
		}
	}
	return nil
}

func (d *CDisplay) IsRunning() bool {
	return d.running
}
