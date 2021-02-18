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
	"strconv"
	"strings"
)

type Property string

func (p Property) String() string {
	return string(p)
}

type CProperty struct {
	name      Property
	kind      PropertyType
	write     bool
	buildable bool
	def       interface{}
	value     interface{}
}

func NewProperty(name Property, kind PropertyType, write bool, buildable bool, def interface{}) (property *CProperty) {
	property = new(CProperty)
	property.name = name
	property.kind = kind
	property.write = write
	property.buildable = buildable
	property.def = def
	property.value = def
	return
}

func (p *CProperty) Name() Property {
	return p.name
}

func (p *CProperty) Type() PropertyType {
	return p.kind
}

func (p *CProperty) ReadOnly() bool {
	return !p.write
}

func (p *CProperty) Buildable() bool {
	return p.buildable
}

func (p *CProperty) Set(value interface{}) error {
	if p.write {
		p.value = value
		return nil
	}
	return fmt.Errorf("error setting read-only property: %v", p.name)
}

func (p *CProperty) SetFromString(value string) error {
	switch p.Type() {
	case BoolProperty:
		switch strings.ToLower(value) {
		case "true", "t", "1":
			return p.Set(true)
		}
		return p.Set(false)
	case StringProperty:
		return p.Set(value)
	case IntProperty:
		if strings.HasSuffix(value, "px") {
			value = strings.Replace(value, "px", "", -1)
		}
		if strings.HasSuffix(value, "%") {
			value = strings.Replace(value, "%", "", -1)
		}
		if v, err := strconv.Atoi(value); err != nil {
			return err
		} else {
			return p.Set(v)
		}
	case FloatProperty:
		if v, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			return p.Set(v)
		}
	case ColorProperty:
		if c, ok := ParseColor(value); ok {
			return p.Set(c)
		} else {
			return fmt.Errorf("invalid color value: %v", value)
		}
	case StyleProperty:
		if c, err := ParseStyle(value); err != nil {
			return err
		} else {
			return p.Set(c)
		}
	case ThemeProperty:
		return fmt.Errorf("theme property not supported by builder features")
	case PointProperty:
		if v, ok := ParsePoint2I(value); ok {
			return p.Set(v)
		} else {
			return fmt.Errorf("invalid point value: %v", value)
		}
	case RectangleProperty:
		if v, ok := ParseRectangle(value); ok {
			return p.Set(v)
		} else {
			return fmt.Errorf("invalid rectangle value: %v", value)
		}
	case RegionProperty:
		if v, ok := ParseRegion(value); ok {
			return p.Set(v)
		} else {
			return fmt.Errorf("invalid region value: %v", value)
		}
	case StructProperty:
		return fmt.Errorf("struct property not supported by builder features")
	}
	return fmt.Errorf("error")
}

func (p *CProperty) Default() (def interface{}) {
	def = p.def
	return
}

func (p *CProperty) Value() (value interface{}) {
	if p.value == nil {
		value = p.def
	} else {
		value = p.value
	}
	return
}
