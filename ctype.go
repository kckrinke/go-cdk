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

// Imaginary Type Tag
//
// Used to denote a concrete type identity
type ITypeTag string

// Stringer interface implementation
func (tag ITypeTag) String() string {
	return string(tag)
}

// Base ITYPE tags
const (
	InvalidITypeID int      = -1
	ITypeNIL       ITypeTag = ""
)

type CType struct {
	tag   ITypeTag
	items []interface{}

	sync.Mutex
}

func NewCType(tag ITypeTag) *CType {
	return &CType{
		tag:   tag,
		items: make([]interface{}, 0),
	}
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

func (t *CType) Remove(item interface{}) error {
	var idx int
	var itm interface{}
	for idx, itm = range t.items {
		if itm == item {
			break
		}
	}
	if idx >= len(t.items) {
		return fmt.Errorf("item not found")
	}
	t.items = append(
		t.items[:idx],
		t.items[idx+1:],
	)
	return nil
}
