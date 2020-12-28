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
	ITypesManager = NewTypeRegistry()
)

type TypeRegistry interface {
	AddType(tag ITypeTag) error
	AddTypeItem(tag ITypeTag, item interface{}) (id int, err error)
	GetTypeItems(tag ITypeTag) []interface{}
	RemoveTypeItem(tag ITypeTag, item interface{}) error
}

type CTypeRegistry struct {
	register map[ITypeTag]*CType
	tracking []interface{}

	sync.Mutex
}

func NewTypeRegistry() TypeRegistry {
	r := &CTypeRegistry{}
	r.register = make(map[ITypeTag]*CType)
	return r
}

func (r *CTypeRegistry) AddType(tag ITypeTag) error {
	r.Lock()
	defer r.Unlock()
	if tag == ITypeNIL {
		return fmt.Errorf("cannot add NIL IType")
	}
	if _, ok := r.register[tag]; ok {
		return fmt.Errorf("existing CType: %v", tag)
	}
	r.register[tag] = NewCType(tag)
	return nil
}

func (r *CTypeRegistry) AddTypeItem(tag ITypeTag, item interface{}) (id int, err error) {
	r.Lock()
	defer r.Unlock()
	if tag == ITypeNIL {
		id, err = -1, fmt.Errorf("cannot add to NIL IType")
		return
	}
	if _, ok := r.register[tag]; !ok {
		id, err = -1, fmt.Errorf("unknown CType: %v", tag)
		return
	}
	_ = r.register[tag].Add(item)
	var ri interface{}
	for id, ri = range r.tracking {
		if ri == item {
			break
		}
	}
	if id == -1 {
		r.tracking = append(r.tracking, item)
		id = len(r.tracking) - 1
	}
	return
}

func (r *CTypeRegistry) GetTypeItems(tag ITypeTag) []interface{} {
	r.Lock()
	defer r.Unlock()
	if tag != ITypeNIL {
		if t, ok := r.register[tag]; ok {
			return t.Items()
		}
	}
	return nil
}

func (r *CTypeRegistry) RemoveTypeItem(tag ITypeTag, item interface{}) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.register[tag]; ok {
		if err := r.register[tag].Remove(item); err != nil {
			return err
		}
		var idx int = -1
		var itm interface{}
		for idx, itm = range r.tracking {
			if itm == item {
				break
			}
		}
		if idx > -1 {
			r.tracking[idx] = nil
			return nil
		}
		return fmt.Errorf("tracking CTypeItem not found")
	}
	return fmt.Errorf("unknown CType: %v", tag)
}
