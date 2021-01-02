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
	TypeObject        CTypeTag = "cdk-object"
	SignalDestroy     Signal   = "destroy"
	SignalSetProperty Signal   = "set-property"
	SignalObjectInit  Signal   = "object-init"
)

func init() {
	_ = TypesManager.AddType(TypeObject)
}

// Basic object type
type Object interface {
	Signaling

	Destroy()

	GetTheme() Theme
	SetTheme(theme Theme)

	SetAllocation(size Rectangle)
	GetAllocation() Rectangle

	SetProperty(name string, value interface{})
	GetProperty(name string) interface{}
	GetPropertyAsBool(name string, def bool) bool
	GetPropertyAsString(name string, def string) string
	GetPropertyAsInt(name string, def int) int
	GetPropertyAsFloat(name string, def float64) float64
}

type CObject struct {
	CSignaling

	theme      Theme
	allocation Rectangle
	properties map[string]interface{}
}

func (o *CObject) Init() (already bool) {
	if o.InitTypeItem(TypeObject) {
		return true
	}
	o.CSignaling.Init()
	o.theme = DefaultColorTheme
	o.properties = make(map[string]interface{})
	o.Emit(SignalObjectInit, o)
	return false
}

func (o *CObject) Destroy() {
	if f := o.Emit(SignalDestroy, o); f == EVENT_PASS {
		if err := o.DestroyObject(); err != nil {
			o.LogErr(err)
		}
	}
}

func (o *CObject) SetAllocation(size Rectangle) {
	o.allocation.SetArea(size.W, size.H)
	o.allocation.Floor(0, 0)
}

func (w *CObject) GetAllocation() Rectangle {
	return w.allocation
}

func (o *CObject) GetTheme() Theme {
	return o.theme
}

func (o *CObject) SetTheme(theme Theme) {
	o.theme = theme
}

// set the value for a named property
func (o *CObject) SetProperty(name string, value interface{}) {
	if f := o.Emit(SignalSetProperty, o, name, value); f == EVENT_PASS {
		o.properties[name] = value
	}
}

// return the named property value
func (o *CObject) GetProperty(name string) interface{} {
	if v, ok := o.properties[name]; ok {
		return v
	}
	return nil
}

// return the named property value as a string
func (o *CObject) GetPropertyAsBool(name string, def bool) bool {
	v := o.GetProperty(name)
	if v, ok := v.(bool); ok {
		return v
	}
	return def
}

// return the named property value as a string
func (o *CObject) GetPropertyAsString(name string, def string) string {
	v := o.GetProperty(name)
	if v, ok := v.(string); ok {
		return v
	}
	return def
}

// return the named property value as an integer
func (o *CObject) GetPropertyAsInt(name string, def int) int {
	v := o.GetProperty(name)
	if v, ok := v.(int); ok {
		return v
	}
	return def
}

// return the named property value as a float
func (o *CObject) GetPropertyAsFloat(name string, def float64) float64 {
	v := o.GetProperty(name)
	if v, ok := v.(float64); ok {
		return v
	}
	return def
}
