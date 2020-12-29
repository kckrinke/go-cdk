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
	"sync"
)

var (
	TypesManager = NewTypeRegistry()
)

type TypeRegistry interface {
	AddType(tag TypeTag) error
	HasType(tag TypeTag) bool
	GetType(tag TypeTag) (t Type, found bool)
	AddTypeItem(tag TypeTag, item TypeItem) (id int, err error)
	GetTypeItems(tag TypeTag) []TypeItem
	RemoveTypeItem(tag TypeTag, item TypeItem) error

	sync.Locker
}

type CTypeRegistry struct {
	register map[TypeTag]Type
	tracking CTypeItemList

	sync.Mutex
}

func NewTypeRegistry() TypeRegistry {
	r := &CTypeRegistry{}
	r.register = make(map[TypeTag]Type)
	r.tracking = make(CTypeItemList, 0)
	return r
}

func (r *CTypeRegistry) AddType(tag TypeTag) error {
	r.Lock()
	defer r.Unlock()
	if tag == TypeNil {
		return fmt.Errorf("cannot add nil type")
	}
	if _, ok := r.register[tag]; ok {
		return fmt.Errorf("type %v exists already", tag)
	}
	r.register[tag] = NewType(tag)
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

func (r *CTypeRegistry) AddTypeItem(tag TypeTag, item TypeItem) (id int, err error) {
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
	_ = r.register[tag].Add(item)
	r.tracking = append(r.tracking, item)
	id = len(r.tracking) - 1
	return
}

func (r *CTypeRegistry) GetTypeItems(tag TypeTag) []TypeItem {
	r.Lock()
	defer r.Unlock()
	if t, ok := r.register[tag]; ok {
		return t.Items()
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
	idx := r.tracking.Index(item)
	if idx == -1 {
		return fmt.Errorf("item not found")
	}
	r.tracking[idx] = nil
	if err := r.register[tag].Remove(item); err != nil {
		return err
	}
	return nil
}
