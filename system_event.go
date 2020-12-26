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
)

type SystemEventType string
type SystemEventFn = func() error

const (
	SYSTEM_EVENT_BOOT     SystemEventType = "boot"
	SYSTEM_EVENT_INIT     SystemEventType = "init"
	SYSTEM_EVENT_SHUTDOWN SystemEventType = "shutdown"
)

var _cdk_system_events = make(map[SystemEventType]map[string]SystemEventFn)

func AddSystemEventHandler(event SystemEventType, tag string, fn SystemEventFn) error {
	if _, ok := _cdk_system_events[event]; !ok {
		_cdk_system_events[event] = make(map[string]SystemEventFn)
	}
	if v, ok := _cdk_system_events[event][tag]; ok {
		return fmt.Errorf("existing %v event %v handler found: %T", event, tag, v)
	}
	_cdk_system_events[event][tag] = fn
	return nil
}

func DelSystemEventHandler(event SystemEventType, tag string) error {
	if _, ok := _cdk_system_events[event]; ok {
		if _, ok := _cdk_system_events[event][tag]; ok {
			delete(_cdk_system_events[event], tag)
			return nil
		}
		return fmt.Errorf("system %v event, tag not found: %v", event, tag)
	}
	return fmt.Errorf("system event not found: %v", event)
}

func HandleSystemEvent(event SystemEventType) error {
	if _, ok := _cdk_system_events[event]; ok {
		for tag, fn := range _cdk_system_events[event] {
			if err := fn(); err != nil {
				return fmt.Errorf("system event %v error: [%v] %v", event, tag, err)
			}
		}
	}
	return nil
}
