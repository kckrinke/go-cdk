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

const (
	TypeObject        CTypeTag = "cdk-object"
	SignalDestroy     Signal   = "destroy"
	SignalSetProperty Signal   = "set-property"
	SignalObjectInit  Signal   = "object-init"
	PropertyDebug     Property = "debug"
)

func init() {
	_ = TypesManager.AddType(TypeObject)
}

// Basic object type
type Object interface {
	Signaling

	Init() (already bool)
	Destroy()
	GetTheme() Theme
	SetTheme(theme Theme)
	GetThemeRequest() (theme Theme)
	SetThemeRequest(theme Theme)
	RegisterProperty(name Property, write bool, def interface{}) error
	SetProperty(name Property, value interface{}) error
	GetProperty(name Property) interface{}
	GetPropertyAsBool(name Property, def bool) bool
	GetPropertyAsString(name Property, def string) string
	GetPropertyAsInt(name Property, def int) int
	GetPropertyAsFloat(name Property, def float64) float64
}

type CObject struct {
	CSignaling

	theme        Theme
	themeRequest *Theme
	properties   []*cObjectProperty
}

func (o *CObject) Init() (already bool) {
	if o.InitTypeItem(TypeObject) {
		return true
	}
	o.CSignaling.Init()
	o.theme = DefaultColorTheme
	o.themeRequest = nil
	o.properties = make([]*cObjectProperty, 0)
	o.RegisterProperty(PropertyDebug, true, false)
	return false
}

func (o *CObject) Destroy() {
	if f := o.Emit(SignalDestroy, o); f == EVENT_PASS {
		if err := o.DestroyObject(); err != nil {
			o.LogErr(err)
		}
	}
}

func (o *CObject) GetTheme() Theme {
	return o.theme
}

func (o *CObject) SetTheme(theme Theme) {
	o.theme = theme
}

func (o *CObject) GetThemeRequest() (theme Theme) {
	if o.themeRequest != nil {
		return *o.themeRequest
	}
	theme = o.GetTheme()
	return
}

func (o *CObject) SetThemeRequest(theme Theme) {
	o.themeRequest = &theme
}

func (o *CObject) RegisterProperty(name Property, write bool, def interface{}) error {
	existing := o.getProperty(name)
	if existing != nil {
		return fmt.Errorf("property exists: %v", name)
	}
	o.properties = append(
		o.properties,
		newProperty(name, write, def),
	)
	return nil
}

func (o *CObject) getProperty(name Property) *cObjectProperty {
	for _, prop := range o.properties {
		if prop.Name() == name {
			return prop
		}
	}
	return nil
}

// set the value for a named property
func (o *CObject) SetProperty(name Property, value interface{}) error {
	if prop := o.getProperty(name); prop != nil {
		if prop.ReadOnly() {
			return fmt.Errorf("cannot set read-only property: %v", name)
		}
		if f := o.Emit(SignalSetProperty, o, name, value); f == EVENT_PASS {
			if err := prop.Set(value); err != nil {
				return err
			}
		}
	}
	return nil
}

// return the named property value
func (o *CObject) GetProperty(name Property) interface{} {
	if prop := o.getProperty(name); prop != nil {
		return prop.Value()
	}
	return nil
}

// return the named property value as a string
func (o *CObject) GetPropertyAsBool(name Property, def bool) bool {
	v := o.GetProperty(name)
	if v, ok := v.(bool); ok {
		return v
	}
	return def
}

// return the named property value as a string
func (o *CObject) GetPropertyAsString(name Property, def string) string {
	v := o.GetProperty(name)
	if v, ok := v.(string); ok {
		return v
	}
	return def
}

// return the named property value as an integer
func (o *CObject) GetPropertyAsInt(name Property, def int) int {
	v := o.GetProperty(name)
	if v, ok := v.(int); ok {
		return v
	}
	return def
}

// return the named property value as a float
func (o *CObject) GetPropertyAsFloat(name Property, def float64) float64 {
	v := o.GetProperty(name)
	if v, ok := v.(float64); ok {
		return v
	}
	return def
}
