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
	IsProperty(name Property) bool
	RegisterProperty(name Property, kind PropertyType, write bool, def interface{}) error
	GetBoolProperty(name Property) (value bool, err error)
	SetBoolProperty(name Property, value bool) error
	GetStringProperty(name Property) (value string, err error)
	SetStringProperty(name Property, value string) error
	GetIntProperty(name Property) (value int, err error)
	SetIntProperty(name Property, value int) error
	GetFloatProperty(name Property) (value float64, err error)
	SetFloatProperty(name Property, value float64) error
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
	_ = o.RegisterProperty(PropertyDebug, BoolProperty, true, false)
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

func (o *CObject) IsProperty(name Property) bool {
	if prop := o.getProperty(name); prop != nil {
		return true
	}
	return false
}

func (o *CObject) RegisterProperty(name Property, kind PropertyType, write bool, def interface{}) error {
	existing := o.getProperty(name)
	if existing != nil {
		return fmt.Errorf("property exists: %v", name)
	}
	o.properties = append(
		o.properties,
		newProperty(name, kind, write, def),
	)
	return nil
}

func (o *CObject) GetBoolProperty(name Property) (value bool, err error) {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == BoolProperty {
			if v, ok := prop.Value().(bool); ok {
				return v, nil
			}
			if v, ok := prop.Default().(bool); ok {
				return v, nil
			}
		}
		return false, fmt.Errorf("%v.(%v) property is not a bool", name, prop.Type())
	}
	return false, fmt.Errorf("property not found: %v", name)
}

func (o *CObject) SetBoolProperty(name Property, value bool) error {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == BoolProperty {
			return o.setProperty(name, value)
		}
		return fmt.Errorf("%v.(%v) property is not a bool", name, prop.Type())
	}
	return fmt.Errorf("property not found: %v", name)
}

func (o *CObject) GetStringProperty(name Property) (value string, err error) {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == StringProperty {
			if v, ok := prop.Value().(string); ok {
				return v, nil
			}
			if v, ok := prop.Default().(string); ok {
				return v, nil
			}
		}
		return "", fmt.Errorf("%v.(%v) property is not a string", name, prop.Type())
	}
	return "", fmt.Errorf("property not found: %v", name)
}

func (o *CObject) SetStringProperty(name Property, value string) error {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == StringProperty {
			return o.setProperty(name, value)
		}
		return fmt.Errorf("%v.(%v) property is not a string", name, prop.Type())
	}
	return fmt.Errorf("property not found: %v", name)
}

func (o *CObject) GetIntProperty(name Property) (value int, err error) {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == IntProperty {
			if v, ok := prop.Value().(int); ok {
				return v, nil
			}
			if v, ok := prop.Default().(int); ok {
				return v, nil
			}
		}
		return 0, fmt.Errorf("%v.(%v) property is not an int", name, prop.Type())
	}
	return 0, fmt.Errorf("property not found: %v", name)
}

func (o *CObject) SetIntProperty(name Property, value int) error {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == IntProperty {
			return o.setProperty(name, value)
		}
		return fmt.Errorf("%v.(%v) property is not an int", name, prop.Type())
	}
	return fmt.Errorf("property not found: %v", name)
}

func (o *CObject) GetFloatProperty(name Property) (value float64, err error) {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == FloatProperty {
			if v, ok := prop.Value().(float64); ok {
				return v, nil
			}
			if v, ok := prop.Default().(float64); ok {
				return v, nil
			}
		}
		return 0.0, fmt.Errorf("%v.(%v) property is not a float", name, prop.Type())
	}
	return 0.0, fmt.Errorf("property not found: %v", name)
}

func (o *CObject) SetFloatProperty(name Property, value float64) error {
	if prop := o.getProperty(name); prop != nil {
		if prop.Type() == FloatProperty {
			return o.setProperty(name, value)
		}
		return fmt.Errorf("%v.(%v) property is not a float64", name, prop.Type())
	}
	return fmt.Errorf("property not found: %v", name)
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
func (o *CObject) setProperty(name Property, value interface{}) error {
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
