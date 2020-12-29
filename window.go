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
	ITypeWindow      ITypeTag = "window"
	SignalDraw       Signal   = "draw"
	SignalSetTitle   Signal   = "set-title"
	SignalSetDisplay Signal   = "set-display"
)

func init() {
	ITypesManager.AddType(ITypeWindow)
}

// Basic window interface
type Window interface {
	Object

	GetTitle() string
	SetTitle(title string)

	GetDisplay() Display
	SetDisplay(d Display)

	Draw(canvas *Canvas) EventFlag
	ProcessEvent(evt Event) EventFlag
}

// Basic window type
type CWindow struct {
	CObject

	title string

	display Display
}

func (w *CWindow) Init() bool {
	w.SetIType(ITypeWindow)
	if w.CObject.Init() {
		return true
	}
	ITypesManager.AddTypeItem(ITypeWindow, w)
	return false
}

func (w *CWindow) SetTitle(title string) {
	if f := w.Emit(SignalSetTitle, w, title); f == EVENT_PASS {
		w.title = title
	}
}

func (w *CWindow) GetTitle() string {
	return w.title
}

func (w *CWindow) GetDisplay() Display {
	return w.display
}

func (w *CWindow) SetDisplay(d Display) {
	if f := w.Emit(SignalSetDisplay, w, d); f == EVENT_PASS {
		w.display = d
	}
}

func (w *CWindow) Draw(canvas *Canvas) EventFlag {
	return w.Emit(SignalDraw, w, canvas)
}

func (w *CWindow) ProcessEvent(evt Event) EventFlag {
	return w.Emit(SignalEvent, w, evt)
}
