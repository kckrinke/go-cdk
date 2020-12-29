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

// ITYPE - Imaginary Type System
//
// This is a simple system for tracking and maintaining arbitrary classes of
// `interface{}` based objects.

import (
	"fmt"
	"sync"
)

// Base ITYPE tags
const (
	InvalidITypeID int      = -1
	TypeNil        CTypeTag = ""
)

type Type interface {
	Items() []TypeItem
	Add(item TypeItem) (id int)
	Remove(item TypeItem) error
}

type CType struct {
	tag   TypeTag
	items []TypeItem

	sync.Mutex
}

func NewType(tag TypeTag) Type {
	return &CType{
		tag:   tag,
		items: make([]TypeItem, 0),
	}
}

func (t *CType) Items() []TypeItem {
	t.Lock()
	defer t.Unlock()
	return t.items
}

func (t *CType) Add(item TypeItem) (id int) {
	t.Lock()
	defer t.Unlock()
	t.items = append(t.items, item)
	id = len(t.items) - 1
	return
}

func (t *CType) Remove(item TypeItem) error {
	var idx int
	var itm TypeItem
	for idx, itm = range t.items {
		if itm.ObjectID() == item.ObjectID() {
			break
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
		t.items = []TypeItem{}
	}
	return nil
}
