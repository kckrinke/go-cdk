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

const (
	TypeOffscreenWindow       CTypeTag = "cdk-offscreen-window"
)

func init() {
	_ = TypesManager.AddType(TypeOffscreenWindow)
}

// Basic window interface
type OffscreenWindow interface {
	Object

	GetTitle() string
	SetTitle(title string)

	GetDisplayManager() DisplayManager
	SetDisplayManager(d DisplayManager)

	Draw(canvas *Canvas) EventFlag
	ProcessEvent(evt Event) EventFlag
}

// Basic window type
type COffscreenWindow struct {
	CObject

	title   string
	display OffscreenDisplay
}

func NewOffscreenWindow(title string) Window {
	d, err := MakeOffscreenDisplay(GetCharset())
	if err != nil {
		Fatal(err)
	}
	w := &COffscreenWindow{
		title:   title,
		display: d,
	}
	w.Init()
	return w
}

func (w *COffscreenWindow) Init() bool {
	if w.InitTypeItem(TypeWindow) {
		return true
	}
	w.CObject.Init()
	return false
}

func (w *COffscreenWindow) SetTitle(title string) {
	if f := w.Emit(SignalSetTitle, w, title); f == EVENT_PASS {
		w.title = title
	}
}

func (w *COffscreenWindow) GetTitle() string {
	return w.title
}

func (w *COffscreenWindow) GetDisplayManager() DisplayManager {
	// return w.display
	return nil
}

func (w *COffscreenWindow) SetDisplayManager(d DisplayManager) {
	if f := w.Emit(SignalSetDisplay, w, d); f == EVENT_PASS {
		// w.display = d
	}
}

func (w *COffscreenWindow) Draw(canvas *Canvas) EventFlag {
	return w.Emit(SignalDraw, w, canvas)
}

func (w *COffscreenWindow) ProcessEvent(evt Event) EventFlag {
	return w.Emit(SignalEvent, w, evt)
}
