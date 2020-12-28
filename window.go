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
	ITypeWindow ITypeTag = "window"
)

func init() {
	ITypesManager.AddType(ITypeWindow)
}

// Basic window interface
type Window interface {
	Object

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

func (w *CWindow) GetDisplay() Display {
	return w.display
}

func (w *CWindow) SetDisplay(d Display) {
	w.display = d
}

func (w *CWindow) Draw(canvas *Canvas) EventFlag {
	w.LogDebug("method not implemented")
	return EVENT_PASS
}

func (w *CWindow) ProcessEvent(evt Event) EventFlag {
	w.LogDebug("method not implemented")
	return EVENT_PASS
}
