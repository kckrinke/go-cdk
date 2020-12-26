// +build ignore

// Copyright 2019 The TCell Authors
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

// beep makes a beep every second until you press ESC
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kckrinke/go-cdk"
)

func main() {
	cdk.SetEncodingFallback(cdk.EncodingFallbackASCII)
	s, e := cdk.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(cdk.StyleDefault)
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *cdk.EventKey:
				switch ev.Key() {
				case cdk.KeyEscape, cdk.KeyEnter, cdk.KeyCtrlC:
					close(quit)
					return
				case cdk.KeyCtrlL:
					s.Sync()
				}
			case *cdk.EventResize:
				s.Sync()
			}
		}
	}()
	beep(s, quit)
	s.Close()
}

func beep(s cdk.Screen, quit <-chan struct{}) {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-quit:
			return
		case <-t.C:
			s.Beep()
		}
	}
}
