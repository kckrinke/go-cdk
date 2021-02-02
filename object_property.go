package cdk

import (
	"fmt"
)

type Property string

func (p Property) String() string {
	return string(p)
}

type cObjectProperty struct {
	name  Property
	kind  PropertyType
	write bool
	def   interface{}
	value interface{}
}

func newProperty(name Property, kind PropertyType, write bool, def interface{}) (property *cObjectProperty) {
	property = new(cObjectProperty)
	property.name = name
	property.kind = kind
	property.write = write
	property.def = def
	property.value = def
	return
}

func (p *cObjectProperty) Name() Property {
	return p.name
}

func (p *cObjectProperty) Type() PropertyType {
	return p.kind
}

func (p *cObjectProperty) ReadOnly() bool {
	return !p.write
}

func (p *cObjectProperty) Set(value interface{}) error {
	if p.write {
		p.value = value
		return nil
	}
	return fmt.Errorf("error setting read-only property: %v", p.name)
}

func (p *cObjectProperty) Default() (def interface{}) {
	def = p.def
	return
}

func (p *cObjectProperty) Value() (value interface{}) {
	if p.value == nil {
		value = p.def
	} else {
		value = p.value
	}
	return
}
