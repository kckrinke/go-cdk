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
	_ = TypesManager.AddType(TypeDisplayManager, func() interface{} { return &CDisplayManager{} })
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
	ActiveCanvas() Canvas
	SetActiveWindow(w Window)
	AddWindow(w Window)
	RemoveWindow(wid int)
	AddWindowOverlay(pid int, overlay Window, region Region)
	RemoveWindowOverlay(pid int, oid int)
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
	windows map[int]*cWindowCanvas
	overlay map[int][]*cWindowCanvas

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

type cWindowCanvas struct {
	window Window
	canvas Canvas
}

func newWindowCanvas(w Window, origin Point2I, size Rectangle, style Style) *cWindowCanvas {
	wc := new(cWindowCanvas)
	wc.canvas = NewCanvas(origin, size, style)
	wc.window = w
	return wc
}

func NewDisplayManager(title string, ttyPath string) *CDisplayManager {
	d := new(CDisplayManager)
	d.title = title
	d.ttyPath = ttyPath
	d.Init()
	return d
}

func (d *CDisplayManager) Init() (already bool) {
	check := TypesManager.GetTypeItems(TypeDisplayManager)
	if len(check) > 0 {
		FatalDF(1, "only one display manager permitted at a time")
	}
	if d.InitTypeItem(TypeDisplayManager, d) {
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

	d.windows = make(map[int]*cWindowCanvas)
	d.overlay = make(map[int][]*cWindowCanvas)
	d.active = -1
	d.SetTheme(DefaultColorTheme)

	cdkDisplayManager = d
	d.Emit(SignalDisplayInit, d)
	return false
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
	// defStyle := DefaultColorStyle.
	// 	Background(ColorReset).
	// 	Foreground(ColorReset)
	d.display.SetStyle(DefaultColorStyle)
	d.display.EnableMouse()
	d.display.EnablePaste()
	d.display.Clear()
	d.captured = true
	d.Emit(SignalDisplayCaptured, d)
}

func (d *CDisplayManager) ReleaseDisplay() {
	if d.captured {
		// d.Lock()
		// defer d.Unlock()
		if d.display != nil {
			d.display.Close()
			d.display = nil
		}
		d.captured = false
	}
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
	if wc, ok := d.windows[d.active]; ok {
		return wc.window
	}
	d.LogWarn("active window not found: %v", d.active)
	return nil
}

func (d *CDisplayManager) ActiveCanvas() Canvas {
	if wc, ok := d.windows[d.active]; ok {
		return wc.canvas
	}
	d.LogWarn("active canvas not found: %v", d.active)
	return nil
}

func (d *CDisplayManager) SetActiveWindow(w Window) {
	// d.Lock()
	if _, ok := d.windows[w.ObjectID()]; !ok {
		// d.Unlock()
		d.AddWindow(w)
		// d.Lock()
	}
	d.active = w.ObjectID()
	// d.Unlock()
}

func (d *CDisplayManager) AddWindow(w Window) {
	// d.Lock()
	// defer d.Unlock()
	if _, ok := d.windows[w.ObjectID()]; ok {
		d.LogWarn("window already added to display: %v", w.ObjectName())
		return
	}
	w.SetDisplayManager(d)
	size := MakeRectangle(0, 0)
	if d.display != nil {
		size = MakeRectangle(d.display.Size())
	}
	d.windows[w.ObjectID()] = newWindowCanvas(w, Point2I{}, size, d.GetTheme().Content.Normal)
	d.overlay[w.ObjectID()] = make([]*cWindowCanvas, 0)
}

func (d *CDisplayManager) RemoveWindow(wid int) {
	if _, ok := d.windows[wid]; ok {
		delete(d.windows, wid)
	}
	if _, ok := d.overlay[wid]; ok {
		delete(d.overlay, wid)
	}
}

func (d *CDisplayManager) AddWindowOverlay(pid int, overlay Window, region Region) {
	if wc, ok := d.overlay[pid]; ok {
		d.overlay[pid] = append(wc, newWindowCanvas(overlay, region.Origin(), region.Size(), d.GetTheme().Content.Normal))
	}
}

func (d *CDisplayManager) RemoveWindowOverlay(pid int, oid int) {
	if wc, ok := d.overlay[pid]; ok {
		var revised []*cWindowCanvas
		for _, oc := range wc {
			if oc.window.ObjectID() != oid {
				revised = append(revised, oc)
			}
		}
		d.overlay[pid] = revised
	}
}

func (d *CDisplayManager) GetWindows() (windows []Window) {
	for _, wc := range d.windows {
		windows = append(windows, wc.window)
	}
	return
}

func (d *CDisplayManager) App() *CApp {
	return d.app
}

func (d *CDisplayManager) ProcessEvent(evt Event) EventFlag {
	if w := d.ActiveWindow(); w != nil {
		if overlays, ok := d.overlay[w.ObjectID()]; ok {
			if last := len(overlays) - 1; last > -1 {
				return overlays[last].window.ProcessEvent(evt)
			}
		}
	}
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
			if ac := d.ActiveCanvas(); ac != nil {
				alloc := MakeRectangle(d.display.Size())
				ac.Resize(alloc, d.GetTheme().Content.Normal)
				if f := aw.ProcessEvent(evt); f == EVENT_STOP {
					return EVENT_STOP
				}
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
	if aw := d.ActiveWindow(); aw != nil {
		if ac := d.ActiveCanvas(); ac != nil {
			if f := window.Draw(ac); f == EVENT_STOP {
				if overlays, ok := d.overlay[aw.ObjectID()]; ok {
					for _, overlay := range overlays {
						if of := overlay.window.Draw(overlay.canvas); of == EVENT_STOP {
							if err := ac.Composite(overlay.canvas); err != nil {
								d.LogErr(err)
							}
						}
					}
				}
				if err := ac.Render(d.display); err != nil {
					d.LogErr(err)
				}
				return EVENT_STOP
			}
		} else {
			d.LogError("missing canvas for active window: %v", aw.ObjectID())
		}
	} else {
		d.LogError("active window not found")
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
