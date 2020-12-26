// Copyright 2020 The TCell Authors
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

// Package extended contains an extended set of terminal descriptions.
// Applications desiring to have a better chance of Just Working by
// default should include this package.  This will significantly increase
// the size of the program.
package extended

import (
	// The following imports just register themselves --
	// these are the terminal types we aggregate in this package.
	_ "github.com/kckrinke/go-cdk/terminfo/a/aixterm"
	_ "github.com/kckrinke/go-cdk/terminfo/a/alacritty"
	_ "github.com/kckrinke/go-cdk/terminfo/a/ansi"
	_ "github.com/kckrinke/go-cdk/terminfo/b/beterm"
	_ "github.com/kckrinke/go-cdk/terminfo/c/cygwin"
	_ "github.com/kckrinke/go-cdk/terminfo/d/dtterm"
	_ "github.com/kckrinke/go-cdk/terminfo/e/emacs"
	_ "github.com/kckrinke/go-cdk/terminfo/g/gnome"
	_ "github.com/kckrinke/go-cdk/terminfo/h/hpterm"
	_ "github.com/kckrinke/go-cdk/terminfo/k/konsole"
	_ "github.com/kckrinke/go-cdk/terminfo/k/kterm"
	_ "github.com/kckrinke/go-cdk/terminfo/l/linux"
	_ "github.com/kckrinke/go-cdk/terminfo/p/pcansi"
	_ "github.com/kckrinke/go-cdk/terminfo/r/rxvt"
	_ "github.com/kckrinke/go-cdk/terminfo/s/screen"
	_ "github.com/kckrinke/go-cdk/terminfo/s/simpleterm"
	_ "github.com/kckrinke/go-cdk/terminfo/s/sun"
	_ "github.com/kckrinke/go-cdk/terminfo/t/termite"
	_ "github.com/kckrinke/go-cdk/terminfo/t/tmux"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt100"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt102"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt220"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt320"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt400"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt420"
	_ "github.com/kckrinke/go-cdk/terminfo/v/vt52"
	_ "github.com/kckrinke/go-cdk/terminfo/w/wy50"
	_ "github.com/kckrinke/go-cdk/terminfo/w/wy60"
	_ "github.com/kckrinke/go-cdk/terminfo/w/wy99_ansi"
	_ "github.com/kckrinke/go-cdk/terminfo/x/xfce"
	_ "github.com/kckrinke/go-cdk/terminfo/x/xterm"
	_ "github.com/kckrinke/go-cdk/terminfo/x/xterm_kitty"
	_ "github.com/kckrinke/go-cdk/terminfo/x/xterm_termite"
)
