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

const (
	TypeSignaling       CTypeTag = "cdk-signaling"
	SignalSignalingInit Signal   = "signaling-init"
)

type Signaling interface {
	TypeItem

	Connect(signal, handle Signal, c SignalListenerFn, data ...interface{})
	Disconnect(signal, handle Signal) error
	Emit(signal Signal, argv ...interface{}) EventFlag
	StopSignal(signal Signal)
	IsSignalStopped(signal Signal) bool
	PassSignal(signal Signal)
	IsSignalPassed(signal Signal) bool
	ResumeSignal(signal Signal)
}

func init() {
	TypesManager.AddType(TypeSignaling)
}

type CSignaling struct {
	CTypeItem

	stopped   []Signal
	passed    []Signal
	listeners map[Signal][]*CSignalListener
}

func (o *CSignaling) Init() (already bool) {
	if o.InitTypeItem(TypeSignaling) {
		return true
	}
	o.CTypeItem.Init()
	o.stopped = []Signal{}
	o.passed = []Signal{}
	o.listeners = make(map[Signal][]*CSignalListener)
	o.Emit(SignalSignalingInit)
	return false
}

// Connect callback to signal, identified by handle
func (o *CSignaling) Connect(signal, handle Signal, c SignalListenerFn, data ...interface{}) {
	if _, ok := o.listeners[signal]; !ok {
		o.listeners[signal] = make([]*CSignalListener, 0)
	}
	o.listeners[signal] = append(
		o.listeners[signal],
		&CSignalListener{
			handle,
			c,
			data,
		},
	)
}

// Disconnect callback from signal identified by handle
func (o *CSignaling) Disconnect(signal, handle Signal) error {
	id := -1
	for i, s := range o.listeners[signal] {
		if s.n == handle {
			id = i
			break
		}
	}
	if id == -1 {
		return fmt.Errorf("unknown signal handle: %v", handle)
	}
	o.LogDebug("disconnecting(%v) from signal(%v)", handle, signal)
	o.listeners[signal] = append(
		o.listeners[signal][:id],
		o.listeners[signal][id+1:]...,
	)
	return nil
}

// Emit a signal event to all connected listener callbacks
func (o *CSignaling) Emit(signal Signal, argv ...interface{}) EventFlag {
	if o.IsSignalStopped(signal) {
		return EVENT_STOP
	}
	if o.IsSignalPassed(signal) {
		return EVENT_PASS
	}
	if listeners, ok := o.listeners[signal]; ok {
		for _, s := range listeners {
			r := s.c(s.d, argv...)
			if r == EVENT_STOP {
				o.LogTrace("emit(%v) stopped by listener(%v)", signal, s.n)
				return EVENT_STOP
			}
		}
	}
	return EVENT_PASS
}

// Disable propagation of the given signal
func (o *CSignaling) StopSignal(signal Signal) {
	if !o.IsSignalStopped(signal) {
		o.LogDebug("stopping signal(%v)", signal)
		o.stopped = append(o.stopped, signal)
	}
}

func (o *CSignaling) IsSignalStopped(signal Signal) bool {
	return o.getSignalStopIndex(signal) >= 0
}

func (o *CSignaling) getSignalStopIndex(signal Signal) int {
	for idx, stop := range o.stopped {
		if signal == stop {
			return idx
		}
	}
	return -1
}

func (o *CSignaling) PassSignal(signal Signal) {
	if !o.IsSignalPassed(signal) {
		o.LogDebug("passing signal(%v)", signal)
		o.passed = append(o.passed, signal)
	}
}

func (o *CSignaling) IsSignalPassed(signal Signal) bool {
	return o.getSignalPassIndex(signal) >= 0
}

func (o *CSignaling) getSignalPassIndex(signal Signal) int {
	for idx, stop := range o.passed {
		if signal == stop {
			return idx
		}
	}
	return -1
}

// Enable propagation of the given signal
func (o *CSignaling) ResumeSignal(signal Signal) {
	id := o.getSignalStopIndex(signal)
	if id >= 0 {
		o.LogDebug("resuming signal(%v) from being stopped", signal)
		o.stopped = append(
			o.stopped[:id],
			o.stopped[id+1:]...,
		)
		return
	}
	id = o.getSignalPassIndex(signal)
	if id >= 0 {
		o.LogDebug("resuming signal(%v) from being passed", signal)
		o.passed = append(
			o.passed[:id],
			o.passed[id+1:]...,
		)
		return
	}
	if _, ok := o.listeners[signal]; ok {
		o.LogWarn("signal(%v) already resumed", signal)
	} else {
		o.LogError("cannot resume unknown signal: %v", signal)
	}
}
