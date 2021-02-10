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
	"sort"
	"sync"
)

var (
	TypesManager = NewTypeRegistry()
)

type TypeRegistry interface {
	GetTypeTags() (tags []TypeTag)
	MakeType(tag TypeTag) (thing interface{}, err error)
	AddType(tag TypeTag, constructor func() interface{}) error
	HasType(tag TypeTag) (exists bool)
	GetType(tag TypeTag) (t Type, found bool)
	AddTypeItem(tag TypeTag, item interface{}) (id int, err error)
	GetTypeItems(tag TypeTag) []interface{}
	GetTypeItemByID(id int) interface{}
	GetTypeItemByName(name string) interface{}
	RemoveTypeItem(tag TypeTag, item TypeItem) error
}

type CTypeRegistry struct {
	register map[TypeTag]Type

	sync.Mutex
}

func NewTypeRegistry() TypeRegistry {
	r := &CTypeRegistry{}
	r.register = make(map[TypeTag]Type)
	return r
}

func (r *CTypeRegistry) GetTypeTags() (tags []TypeTag) {
	for tt, _ := range r.register {
		tags = append(tags, tt)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].String() < tags[j].String()
	})
	return
}

func (r *CTypeRegistry) MakeType(tag TypeTag) (thing interface{}, err error) {
	if t, ok := r.GetType(tag); ok {
		thing = t.New()
		if thing == nil {
			err = fmt.Errorf("type not buildable: %v", tag)
		}
	} else {
		err = fmt.Errorf("type not found: %v", tag)
	}
	return
}

func (r *CTypeRegistry) AddType(tag TypeTag, constructor func() interface{}) error {
	r.Lock()
	defer r.Unlock()
	if tag == TypeNil {
		return fmt.Errorf("cannot add nil type")
	}
	if _, ok := r.register[tag]; ok {
		return fmt.Errorf("type %v exists already", tag)
	}
	r.register[tag] = NewType(tag, constructor)
	return nil
}

func (r *CTypeRegistry) HasType(tag TypeTag) (exists bool) {
	_, exists = r.register[tag]
	return
}

func (r *CTypeRegistry) GetType(tag TypeTag) (t Type, found bool) {
	t, found = r.register[tag]
	return
}

func (r *CTypeRegistry) AddTypeItem(tag TypeTag, item interface{}) (id int, err error) {
	r.Lock()
	defer r.Unlock()
	if tag == TypeNil {
		id, err = -1, fmt.Errorf("cannot add to nil type")
		return
	}
	if _, ok := r.register[tag]; !ok {
		id, err = -1, fmt.Errorf("unknown type: %v", tag)
		return
	}
	r.register[tag].Add(item)
	id = r.GetNextID()
	return
}

func (r *CTypeRegistry) GetNextID() (id int) {
	id = 1 // first valid ID is 1
	for _, t := range r.register {
		for _, ti := range t.Items() {
			if tc, ok := ti.(TypeItem); ok {
				if id == tc.ObjectID() {
					id += 1
				}
			}
		}
	}
	return
}

func (r *CTypeRegistry) GetTypeItems(tag TypeTag) []interface{} {
	r.Lock()
	defer r.Unlock()
	if t, ok := r.register[tag]; ok {
		return t.Items()
	}
	return nil
}

func (r *CTypeRegistry) GetTypeItemByID(id int) interface{} {
	r.Lock()
	defer r.Unlock()
	for _, t := range r.register {
		for _, i := range t.Items() {
			if c, ok := i.(TypeItem); ok {
				if c.ObjectID() == id {
					return i
				}
			}
		}
	}
	return nil
}

func (r *CTypeRegistry) GetTypeItemByName(name string) interface{} {
	r.Lock()
	defer r.Unlock()
	for _, t := range r.register {
		for _, i := range t.Items() {
			if c, ok := i.(TypeItem); ok {
				if c.GetName() == name {
					return i
				}
			}
		}
	}
	return nil
}

func (r *CTypeRegistry) RemoveTypeItem(tag TypeTag, item TypeItem) error {
	r.Lock()
	defer r.Unlock()
	if item == nil || !item.IsValid() {
		return fmt.Errorf("item not valid")
	}
	if _, ok := r.register[tag]; !ok {
		return fmt.Errorf("unknown type: %v", tag)
	}
	if err := r.register[tag].Remove(item); err != nil {
		return err
	}
	return nil
}
