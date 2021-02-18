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

const TypeObject CTypeTag = "cdk-object"

func init() {
	_ = TypesManager.AddType(TypeObject, func() interface{} { return nil })
}

// This is the base type for all complex CDK object types. The Object type
// provides a means of installing properties, getting and setting property
// values
type Object interface {
	MetaData

	Init() (already bool)
	InitWithProperties(properties map[Property]string) (already bool, err error)
	Destroy()
	GetName() (name string)
	SetName(name string)
	GetTheme() (theme Theme)
	SetTheme(theme Theme)
	GetThemeRequest() (theme Theme)
	SetThemeRequest(theme Theme)
}

type CObject struct {
	CMetaData
}

func (o *CObject) Init() (already bool) {
	if o.InitTypeItem(TypeObject, o) {
		return true
	}
	o.CSignaling.Init()
	o.properties = make([]*CProperty, 0)
	_ = o.InstallProperty(PropertyDebug, BoolProperty, true, false)
	_ = o.InstallProperty(PropertyName, StringProperty, true, nil)
	_ = o.InstallProperty(PropertyTheme, ThemeProperty, true, DefaultColorTheme)
	_ = o.InstallProperty(PropertyThemeRequest, ThemeProperty, true, DefaultColorTheme)
	o.Connect(
		SignalSetProperty,
		Signal(fmt.Sprintf("%v.set-property--name", o.ObjectName())),
		func(data []interface{}, argv ...interface{}) EventFlag {
			if len(argv) == 3 {
				if key, ok := argv[1].(Property); ok && key == PropertyName {
					if name, ok := argv[2].(string); ok {
						o.CTypeItem.SetName(name)
					}
				}
			}
			return EVENT_PASS
		})
	return false
}

func (o *CObject) InitWithProperties(properties map[Property]string) (already bool, err error) {
	if o.Init() {
		return true, nil
	}
	if err = o.SetProperties(properties); err != nil {
		return false, err
	}
	return false, nil
}

func (o *CObject) Destroy() {
	if f := o.Emit(SignalDestroy, o); f == EVENT_PASS {
		if err := o.DestroyObject(); err != nil {
			o.LogErr(err)
		}
	}
}

func (o *CObject) GetName() (name string) {
	var err error
	if name, err = o.GetStringProperty(PropertyName); err != nil {
		return ""
	}
	return
}

func (o *CObject) SetName(name string) {
	if err := o.SetStringProperty(PropertyName, name); err != nil {
		o.LogErr(err)
	}
}

func (o *CObject) GetTheme() (theme Theme) {
	var err error
	if theme, err = o.GetThemeProperty(PropertyTheme); err != nil {
		o.LogErr(err)
	}
	return
}

func (o *CObject) SetTheme(theme Theme) {
	if err := o.SetThemeProperty(PropertyTheme, theme); err != nil {
		o.LogErr(err)
	}
}

func (o *CObject) GetThemeRequest() (theme Theme) {
	var err error
	if theme, err = o.GetThemeProperty(PropertyThemeRequest); err != nil {
		o.LogErr(err)
		theme = o.GetTheme()
	}
	return
}

func (o *CObject) SetThemeRequest(theme Theme) {
	if err := o.SetThemeProperty(PropertyThemeRequest, theme); err != nil {
		o.LogErr(err)
	}
}

// emitted when the object instance is destroyed
const SignalDestroy Signal = "destroy"

// request that the object be rendered with additional features useful to
// debugging custom Widget development
const PropertyDebug Property = "debug"

// property wrapper around the CTypeItem name field
const PropertyName Property = "name"

const PropertyTheme        Property = "theme"

const PropertyThemeRequest Property = "theme-request"
