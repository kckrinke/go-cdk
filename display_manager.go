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
	"time"
)

// display is really a wrapper around Display
// and Simulation screens

// basically a wrapper around Display()
// manages one or more windows backed by viewports
// viewports manage the allocation of space
// drawables within viewports render the space

var (
	DisplayCallQueueCapacity = 16
	cdkDisplayManager        DisplayManager
)

const (
	TypeDisplayManager    CTypeTag = "cdk-display-manager"
	SignalDisplayInit     Signal   = "display-init"
	SignalDisplayCaptured Signal   = "display-captured"
	SignalInterrupt       Signal   = "sigint"
	SignalEvent           Signal   = "event"
	SignalEventError      Signal   = "event-error"
	SignalEventKey        Signal   = "event-key"
	SignalEventMouse      Signal   = "event-mouse"
	SignalEventResize     Signal   = "event-resize"
)

func init() {
	_ = TypesManager.AddType(TypeDisplayManager)
}

type DisplayCallbackFn = func(d DisplayManager) error

type DisplayManager interface {
	Object

	GetTitle() string
	SetTitle(title string)

	GetTtyPath() string
	SetTtyPath(ttyPath string)

	Display() Display
	DisplayCaptured() bool
	CaptureDisplay(ttyPath string)
	ReleaseDisplay()
	IsMonochrome() bool
	Colors() (numberOfColors int)

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
type CDisplayManager struct {
	CObject

	title string

	captureCtrlC bool

	active  int
	windows []Window
	wCanvas []Canvas

	app      *CApp
	ttyPath  string
	display  Display
	captured bool

	running  bool
	waiting  bool
	done     chan bool
	queue    chan DisplayCallbackFn
	events   chan Event
	process  chan Event
	requests chan ScreenStateReq
}

func NewDisplayManager(title string, ttyPath string) *CDisplayManager {
	d := new(CDisplayManager)
	d.title = title
	d.ttyPath = ttyPath
	d.Init()
	return d
}

func GetDisplayManager() (dm DisplayManager) {
	dm = cdkDisplayManager
	return
}

func GetCurrentTheme() (theme Theme) {
	theme = DefaultColorTheme
	if cdkDisplayManager != nil {
		theme = cdkDisplayManager.GetTheme()
	}
	return
}

func SetCurrentTheme(theme Theme) {
	if cdkDisplayManager != nil {
		cdkDisplayManager.SetTheme(theme)
	}
}

// Initialization
func (d *CDisplayManager) Init() (already bool) {
	check := TypesManager.GetTypeItems(TypeDisplayManager)
	if len(check) > 0 {
		FatalF("only one display permitted at a time")
	}
	if d.InitTypeItem(TypeDisplayManager) {
		return true
	}
	d.CObject.Init()

	d.captured = false
	d.running = false
	d.waiting = true
	d.done = make(chan bool)
	d.queue = make(chan DisplayCallbackFn, DisplayCallQueueCapacity)
	d.events = make(chan Event, DisplayCallQueueCapacity)
	d.process = make(chan Event, DisplayCallQueueCapacity)
	d.requests = make(chan ScreenStateReq, DisplayCallQueueCapacity)

	d.windows = []Window{}
	d.active = -1
	d.SetTheme(DefaultColorTheme)

	cdkDisplayManager = d
	d.Emit(SignalDisplayInit, d)
	return false
}

func (d *CDisplayManager) Destroy() {
	if d.display != nil {
		d.display.Close()
	}
	close(d.done)
	close(d.queue)
	close(d.process)
	close(d.requests)
	d.CObject.Destroy()
}

func (d *CDisplayManager) GetTitle() string {
	return d.title
}

func (d *CDisplayManager) SetTitle(title string) {
	d.title = title
}

func (d *CDisplayManager) GetTtyPath() string {
	return d.ttyPath
}

func (d *CDisplayManager) SetTtyPath(ttyPath string) {
	d.ttyPath = ttyPath
}

func (d *CDisplayManager) Display() Display {
	return d.display
}

func (d *CDisplayManager) DisplayCaptured() bool {
	return d.display != nil && d.captured
}

func (d *CDisplayManager) CaptureDisplay(ttyPath string) {
	d.Lock()
	defer d.Unlock()
	var err error
	if ttyPath == OffscreenDisplayTtyPath {
		if d.display, err = MakeOffscreenDisplay(""); err != nil {
			FatalF("error getting offscreen display: %v", err)
		}
	} else {
		if d.display, err = NewDisplay(); err != nil {
			FatalF("error getting new display: %v", err)
		}
		if err = d.display.Init(); err != nil {
			FatalF("error initializing new display: %v", err)
		}
	}
	defStyle := StyleDefault.
		Background(ColorReset).
		Foreground(ColorReset)
	d.display.SetStyle(defStyle)
	d.display.EnableMouse()
	d.display.EnablePaste()
	d.display.Clear()
	d.captured = true
	d.Emit(SignalDisplayCaptured, d)
}

func (d *CDisplayManager) ReleaseDisplay() {
	d.Lock()
	defer d.Unlock()
	if d.display != nil {
		d.display.Close()
		d.display = nil
	}
	d.captured = false
}

func (d *CDisplayManager) IsMonochrome() bool {
	return d.Colors() == 0
}

func (d *CDisplayManager) Colors() (numberOfColors int) {
	numberOfColors = 0
	if d.display != nil {
		numberOfColors = d.display.Colors()
	}
	return
}

func (d *CDisplayManager) CaptureCtrlC() {
	d.Lock()
	defer d.Unlock()
	d.captureCtrlC = true
}

func (d *CDisplayManager) ReleaseCtrlC() {
	d.Lock()
	defer d.Unlock()
	d.captureCtrlC = false
}

func (d *CDisplayManager) DefaultTheme() Theme {
	if d.display != nil {
		if d.display.Colors() <= 0 {
			return DefaultMonoTheme
		}
	}
	return DefaultColorTheme
}

func (d *CDisplayManager) ActiveWindow() Window {
	if len(d.windows) > d.active && d.active >= 0 {
		return d.windows[d.active]
	}
	if len(d.windows) == 0 {
		return nil
	}
	d.active = 0
	return d.windows[0]
}

func (d *CDisplayManager) windowIndex(w Window) (index int) {
	var cw Window
	for index, cw = range d.windows {
		if w.ObjectID() == cw.ObjectID() {
			return
		}
	}
	index = -1
	return
}

func (d *CDisplayManager) SetActiveWindow(w Window) {
	d.Lock()
	id := -1
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

func (d *CDisplayManager) AddWindow(w Window) int {
	d.Lock()
	defer d.Unlock()
	id := -1
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
	size := MakeRectangle(0, 0)
	if d.display != nil {
		size = MakeRectangle(d.display.Size())
	}
	d.wCanvas = append(d.wCanvas, NewCanvas(Point2I{}, size, d.GetTheme().Content.Normal))
	w.SetDisplayManager(d)
	d.windows = append(d.windows, w)
	return len(d.windows) - 1
}

func (d *CDisplayManager) GetWindows() []Window {
	return d.windows
}

func (d *CDisplayManager) App() *CApp {
	return d.app
}

func (d *CDisplayManager) ProcessEvent(evt Event) EventFlag {
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
		if aw := d.ActiveWindow(); aw != nil {
			wid := d.windowIndex(aw)
			w, h := d.display.Size()
			alloc := MakeRectangle(w, h)
			if d.wCanvas[wid] != nil {
				d.wCanvas[wid].Resize(alloc, d.GetTheme().Content.Normal)
				if f := aw.ProcessEvent(evt); f == EVENT_STOP {
					return EVENT_STOP
				}
			} else {
				d.LogError("missing canvas for wid: %v", wid)
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

func (d *CDisplayManager) DrawScreen() EventFlag {
	d.Lock()
	defer d.Unlock()
	if !d.captured || d.display == nil {
		d.LogError("display not captured or otherwise missing")
		return EVENT_PASS
	}
	var window Window
	if window = d.ActiveWindow(); window == nil {
		d.LogDebug("cannot draw the display, display missing a window")
		return EVENT_PASS
	}
	wid := d.windowIndex(window)
	if len(d.wCanvas) < wid {
		d.LogError("missing canvas for window: %v", window.ObjectName())
		return EVENT_PASS
	}
	if f := window.Draw(d.wCanvas[wid]); f == EVENT_STOP {
		if err := d.wCanvas[wid].Render(d.display); err != nil {
			d.LogErr(err)
		}
		return EVENT_STOP
	}
	return EVENT_PASS
}

func (d *CDisplayManager) RequestDraw() {
	if d.running {
		d.requests <- DrawRequest
	} else {
		TraceF("application not running")
	}
}

func (d *CDisplayManager) RequestShow() {
	if d.running {
		d.requests <- ShowRequest
	} else {
		TraceF("application not running")
	}
}

func (d *CDisplayManager) RequestSync() {
	if d.running {
		d.requests <- SyncRequest
	} else {
		TraceF("application not running")
	}
}

func (d *CDisplayManager) RequestQuit() {
	if d.running {
		d.requests <- QuitRequest
	} else {
		TraceF("application not running")
	}
}

func (d *CDisplayManager) AsyncCall(fn DisplayCallbackFn) error {
	if !d.running {
		return fmt.Errorf("application not running")
	}
	d.queue <- fn
	return nil
}

func (d *CDisplayManager) AwaitCall(fn DisplayCallbackFn) error {
	if !d.running {
		return fmt.Errorf("application not running")
	}
	var err error
	done := make(chan bool)
	d.queue <- func(d DisplayManager) error {
		err = fn(d)
		done <- true
		return nil
	}
	<-done
	return err
}

func (d *CDisplayManager) PostEvent(evt Event) error {
	if !d.running {
		return fmt.Errorf("application not running")
	}
	d.events <- evt
	return nil
}

func (d *CDisplayManager) pollEventWorker() {
	for d.running {
		if d.display != nil {
			d.process <- d.display.PollEvent()
		}
	}
}

func (d *CDisplayManager) processEventWorker() {
	for d.running {
		if evt := <-d.process; evt != nil {
			if f := d.ProcessEvent(evt); f == EVENT_STOP {
				// TODO: ProcessEvent must ONLY flag stop when UI changes
				d.RequestDraw()
				d.RequestShow()
			}
		}
	}
}
func (d *CDisplayManager) screenRequestWorker() {
	if d.running {
		if err := d.app.InitUI(); err != nil {
			FatalDF(1, "%v", err)
		}

	}
	for d.running {
		switch <-d.requests {
		case DrawRequest:
			if d.display != nil && !d.waiting {
				d.DrawScreen()
			}
		case ShowRequest:
			if d.display != nil && !d.waiting {
				d.display.Show()
			}
		case SyncRequest:
			if d.display != nil && !d.waiting {
				d.display.Sync()
			}
		case QuitRequest:
			d.done <- true
		}
	}
}

func (d *CDisplayManager) Run() error {
	d.CaptureDisplay(d.ttyPath)
	d.running = true
	go d.pollEventWorker()
	go d.processEventWorker()
	go d.screenRequestWorker()
	defer func() {
		d.ReleaseDisplay()
		close(d.done)
		close(d.events)
		close(d.queue)
		if p := recover(); p != nil {
			panic(p)
		}
	}()
	AddTimeout(time.Millisecond*51, func() EventFlag {
		if d.display != nil {
			d.waiting = false
			if err := d.display.PostEvent(NewEventResize(d.display.Size())); err != nil {
				Error(err)
			}
		}
		return EVENT_STOP
	})
	d.RequestDraw()
	d.RequestSync()
	for d.running {
		select {
		case fn, ok := <-d.queue:
			if !ok {
				d.running = false
				break
			}
			if err := fn(d); err != nil {
				return err
			}
		case evt, ok := <-d.events:
			if !ok {
				d.running = false
				break
			}
			if d.display != nil {
				if err := d.display.PostEvent(evt); err != nil {
					Error(err)
				}
			} else {
				d.LogTrace("missing display, dropping event: %v", evt)
			}
		case <-d.done:
			d.running = false
			break
		}
	}
	return nil
}

func (d *CDisplayManager) IsRunning() bool {
	return d.running
}
