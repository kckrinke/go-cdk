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

// TODO: use UUIDs instead of Timer ID numbers

import (
	"time"
)

var cdkTimeouts *timers = &timers{}

type timers struct {
	timers []*timer
}

func (t *timers) Add(n *timer) (id int) {
	t.timers = append(t.timers, n)
	id = len(t.timers) - 1
	return
}

func (t *timers) Valid(id int) bool {
	if id < len(t.timers) {
		if t.timers[id] != nil {
			return true
		}
	}
	return false
}

func (t *timers) Get(id int) *timer {
	if t.Valid(id) {
		return t.timers[id]
	}
	return nil
}

func (t *timers) Remove(id int) {
	// retain the timer IDs
	if t.Valid(id) {
		t.timers[id] = nil
	}
}

func (t *timers) Stop(id int) bool {
	if t.Valid(id) {
		t.timers[id].timer.Stop()
		t.Remove(id)
		return true
	}
	return false
}

type timer struct {
	id    int
	d     time.Duration
	fn    TimerCallbackFn
	timer *time.Timer
}

func (t *timer) handler() {
	if f := t.fn(); f == EVENT_STOP {
		cdkTimeouts.Remove(t.id)
	} else {
		t.timer.Stop()
		t.timer = time.AfterFunc(t.d, t.handler)
	}
}

type TimerCallbackFn = func() EventFlag

func AddTimeout(d time.Duration, fn TimerCallbackFn) (id int) {
	t := &timer{
		fn: fn,
	}
	t.timer = time.AfterFunc(d, t.handler)
	t.id = cdkTimeouts.Add(t)
	id = t.id
	return
}

func CancelTimeout(id int) bool {
	return cdkTimeouts.Stop(id)
}
