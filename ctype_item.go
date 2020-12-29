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

	IsValid() bool
	String() string

	GetIType() ITypeTag
	GetName() string
	SetName(name string)

	ObjectID() int
	DestroyObject() error
	ObjectName() string

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
	typeTag ITypeTag
	name    string
	valid   bool

	sync.Mutex
}

func (o *CTypeItem) SetIType(tag ITypeTag) {
	if o.typeTag == ITypeNIL {
		o.typeTag = tag
	}
}

func (o *CTypeItem) Init() (already bool) {
	o.Lock()
	defer o.Unlock()
	if o.valid {
		return true
	}
	var err error
	o.id, err = ITypesManager.AddTypeItem(o.typeTag, o)
	if err != nil {
		Fataldf(1, "AddTypeItem(%v) failed: %v", o.typeTag, err)
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

func (o *CTypeItem) GetIType() ITypeTag {
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

func (o *CTypeItem) DestroyObject() error {
	o.valid = false
	o.id = -1
	return ITypesManager.RemoveTypeItem(o.typeTag, o)
}

// returns the unique string ITypeObject identity for this object
func (o *CTypeItem) ObjectName() string {
	if len(o.name) > 0 {
		return fmt.Sprintf("%v-%v-%v", o.typeTag, o.ObjectID(), o.name)
	}
	return fmt.Sprintf("%v-%v", o.typeTag, o.ObjectID())
}

func (o *CTypeItem) LogTag() string {
	if len(o.name) > 0 {
		return fmt.Sprintf("[%v:%v:%v]", o.typeTag, o.ObjectID(), o.name)
	}
	return fmt.Sprintf("[%v:%v]", o.typeTag, o.ObjectID())
}

func (o *CTypeItem) LogTrace(format string, argv ...interface{}) {
	Tracedf(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogDebug(format string, argv ...interface{}) {
	Debugdf(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogInfo(format string, argv ...interface{}) {
	Infodf(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogWarn(format string, argv ...interface{}) {
	Warndf(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogError(format string, argv ...interface{}) {
	Errordf(1, fmt.Sprintf("%s %s", o.LogTag(), format), argv...)
}

func (o *CTypeItem) LogErr(err error) {
	Errordf(1, err.Error())
}
