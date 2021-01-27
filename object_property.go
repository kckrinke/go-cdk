package cdk

import (
	"fmt"
)

type cObjectProperty struct {
	name  string
	write bool
	def   interface{}
	value interface{}
}

func newProperty(name string, write bool, def interface{}) (property *cObjectProperty) {
	property = new(cObjectProperty)
	property.name = name
	property.write = write
	property.def = def
	property.value = def
	return
}

func (p *cObjectProperty) Name() string {
	return p.name
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
	value = p.value
	return
}
