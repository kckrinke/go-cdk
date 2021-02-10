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

const (
	TypeNil CTypeTag = ""
)

type Type interface {
	New() interface{}
	Items() []interface{}
	Add(item interface{}) (id int)
	Remove(item TypeItem) error
}

type CType struct {
	tag   TypeTag
	items []interface{}
	new   func() interface{}

	sync.Mutex
}

func NewType(tag TypeTag, constructor func() interface{}) Type {
	return &CType{
		tag:   tag,
		items: make([]interface{}, 0),
		new:   constructor,
	}
}

func (t *CType) New() interface{} {
	if t.new != nil {
		return t.new()
	}
	return nil
}

func (t *CType) Items() []interface{} {
	t.Lock()
	defer t.Unlock()
	return t.items
}

func (t *CType) Add(item interface{}) (id int) {
	t.Lock()
	defer t.Unlock()
	t.items = append(t.items, item)
	id = len(t.items) - 1
	return
}

func (t *CType) Remove(item TypeItem) error {
	t.Lock()
	defer t.Unlock()
	var idx int
	var itm interface{}
	for idx, itm = range t.items {
		if tt, ok := itm.(TypeItem); ok {
			if tt.ObjectID() == item.ObjectID() {
				break
			}
		}
	}
	count := len(t.items)
	if count > 0 && idx >= count {
		return fmt.Errorf("item not found")
	} else if count > 1 {
		t.items = append(
			t.items[:idx],
			t.items[idx+1:]...,
		)
	} else {
		t.items = []interface{}{}
	}
	return nil
}
