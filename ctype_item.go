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

type TypeItem interface {
	Init() (already bool)
	InitTypeItem(tag TypeTag) (already bool)
	IsValid() bool
	String() string
	GetTypeTag() TypeTag
	GetName() string
	SetName(name string)
	ObjectID() int
	ObjectName() string
	DestroyObject() error
	LogTag() string
	LogTrace(format string, argv ...interface{})
	LogDebug(format string, argv ...interface{})
	LogInfo(format string, argv ...interface{})
	LogWarn(format string, argv ...interface{})
	LogError(format string, argv ...interface{})
	LogErr(err error)

	sync.Locker
}

type CTypeItem struct {
	id      int
	typeTag CTypeTag
	name    string
	valid   bool

	sync.Mutex
}

func NewTypeItem(tag CTypeTag, name string) TypeItem {
	return &CTypeItem{
		id:      -1,
		typeTag: tag,
		name:    name,
		valid:   false,
	}
}

func (o *CTypeItem) InitTypeItem(tag TypeTag) (already bool) {
	o.Lock()
	defer o.Unlock()
	already = o.valid
	if !already && o.typeTag == TypeNil {
		o.typeTag = tag.Tag()
	}
	return
}

func (o *CTypeItem) Init() (already bool) {
	if o.valid {
		return true
	}
	if o.typeTag == TypeNil {
		FatalDF(1, "invalid object type: nil")
	}
	var err error
	o.id, err = TypesManager.AddTypeItem(o.typeTag, o)
	if err != nil {
		FatalDF(1, "failed to add self to \"%v\" type: %v", o.typeTag, err)
	}
	o.valid = true
	return false
}

func (o *CTypeItem) IsValid() bool {
	return o.valid
}

func (o *CTypeItem) String() string {
	return o.ObjectName()
}

func (o *CTypeItem) GetTypeTag() TypeTag {
	return o.typeTag
}

func (o *CTypeItem) GetName() string {
	return o.name
}

func (o *CTypeItem) SetName(name string) {
	o.Lock()
	defer o.Unlock()
	o.name = name
}

func (o *CTypeItem) ObjectID() int {
	return o.id
}

func (o *CTypeItem) ObjectName() string {
	if len(o.name) > 0 {
		return fmt.Sprintf("%v-%v-%v", o.typeTag, o.ObjectID(), o.name)
	}
	return fmt.Sprintf("%v-%v", o.typeTag, o.ObjectID())
}

func (o *CTypeItem) DestroyObject() error {
	err := TypesManager.RemoveTypeItem(o.typeTag, o)
	o.valid = false
	o.id = -1
	return err
}

func (o *CTypeItem) LogTag() string {
	if len(o.name) > 0 {
		return fmt.Sprintf("[%v.%v.%v]", o.typeTag, o.ObjectID(), o.name)
	}
	return fmt.Sprintf("[%v.%v]", o.typeTag, o.ObjectID())
}

func (o *CTypeItem) LogTrace(format string, argv ...interface{}) {
	TraceDF(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogDebug(format string, argv ...interface{}) {
	DebugDF(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogInfo(format string, argv ...interface{}) {
	InfoDF(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogWarn(format string, argv ...interface{}) {
	WarnDF(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogError(format string, argv ...interface{}) {
	ErrorDF(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogErr(err error) {
	ErrorDF(1, err.Error())
}
