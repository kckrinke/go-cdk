// Copyright 2021 The CDK Authors
// Copyright 2018 The TCell Authors
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

// The Curses Development Kit is the low-level compliment to the higher-level
// Curses Tool Kit.
//
// Introduction
//
// CDK is based primarily upon the TCell codebase, however, it is a hard-fork
// with many API-breaking changes. Some of the more salient changes are as
// follows:
//
//  - `Screen` is now a `Display`
//  - `Theme` is a collection of purpose-intended `Style`s
//  - Developers are not intended to use the `Display` directly
//  - a `Canvas` type is used to composite the rendering process
//  - a `Window` type is used to render to a `Display` and handle events
//  - `EventFlag`s are used to manage the propagation of events
//  - robust canvas rendering enables reliable wrapping, justification and
//    alignment of textual content
//  - `Tango` markup format to augment the text rendering with applied style
//    attributes (with similar features as Pango for GDK/GTK)
//
// CDK Demo Walkthrough
//
// All CDK applications require some form of `Window` implementation in order
// to function. One can use the build-in `cdk.CWindow` type to construct a basic
// `Window` object and tie into it's signals in order to render to the canvas
// and handle events. For this example however, a more formal approach is taken.
//
// 		type CdkDemoWindow struct {
// 			cdk.CWindow
// 		}
//
// Starting out with a very slim definition of our custom `Window`, all that's
// necessary is the embed the concrete `cdk.CWindow` type and proceed with
// overriding various methods.
//
// 		func (w *CdkDemoWindow) Init() (already bool) {
// 			if w.CWindow.Init() {
// 				return true
// 			}
// 			return false
// 		}
//
// The `Init` method is not necessary to overload, however, this is a good spot
// to do any UI initializations or the like. For this demo, the boilerplate
// minimum is given as an example to start from.
//
// 		func (w *CdkDemoWindow) Draw(canvas cdk.Canvas) cdk.EventFlag {
// 			... things ...
// 		}
//
// The next method implemented is the `Draw` method. This method receives a
// pre-configured Canvas object, set to the available size of the `Display`.
// Within this method, the application needs to process as little "work" as
// possible and focus primarily on just rendering the content upon the canvas.
// Let's walk through each line of the `Draw` method.
//
// 		w.LogInfo("Draw: %s", canvas)
//
// This line uses the built-in logging facilities to log an "info" message that
// the `Draw` method was invoked and let's us know some sort of human-readable
// description of the canvas (resembles JSON text but isn't really JSON). The
// advantage to using these built-in logging facilities is that the log entry
// itself will be prefixed with some extra metadata identifying the particular
// object instance with a bit of text in the format of "typeName-ID" where
// typeName is the object's CDK-type and the ID is an integer (marking the
// unique instance).
//
//  	theme := w.GetDisplayManager().DefaultTheme()
//
// Within CDK, there is a concept of `Theme`, which really is just a collection
// of useful purpose-driven `Style`s. One can set the default theme for the
// running CDK system, however the stock state is either a monochrome base theme
// or a colorized variant. Some of the rendering functions require `Style`s or
// `Theme`s passed as arguments and so we're getting that here for later use.
//
// 		size := canvas.Size()
//
// Simply getting a `Rectangle` primitive with it's `W`idth and `H`eight values
// set according to the size of the canvas' internal buffer.
//
// `Rectangle` is a CDK primitive which has just two fields: `W` and `H`. Most
// places where spacial bounds are necessary, these primitives are used (such as
// the concept of a box `size` for painting upon a canvas).
//
// 		canvas.Box(cdk.Point2I{}, size, true, true, theme)
//
// This is the first actual draw command. In this case, the `Box` method is
// configured to draw a box on the screen, starting at a position of 0,0 (top
// left corner), taking up the full volume of the canvas, with a border (first
// boolean `true` argument), ensuring to fill the entire area with the filler
// rune and style within a given theme, which is the last argument to the `Box`
// method. On a color-supporting terminal, this will paint a navy-blue box over
// the entire terminal screen.
//
// 		content := "..."
//      content += "..."
//      content += "..."
//
// These few lines of code are merely concatenating a string of `Tango` markup
// that includes use of `<b></b>`, `<u></u>`, `<i></i>`, and `<span></span>`
// tags. All colors have fallback versions and are typically safe even for
// monochrome terminal sessions.
//
// 		textPoint := cdk.MakePoint2I(size.W/2/2, size.H/2-1)
//
// This sets up a variable holding a `Point2I` instance configured for 1/4 of
// the width into the screen (from the left) and halfway minus one of the height
// into the screen (from the top).
//
// `Point2I` is a CDK primitive which has just two fields: `X` and `Y`. Most
// places where coordinates are necessary, these primitives are used (such as
// the concept of an `origin` point for painting upon a canvas).
//
// 		textSize := cdk.MakeRectangle(size.W/2, size.H/2)
//
// This sets up a variable holding a `Rectangle` configured to be half the size
// of the canvas itself.
//
// 		canvas.DrawText(textPoint, textSize, cdk.JUSTIFY_CENTER, false, cdk.WRAP_WORD, cdk.DefaultColorTheme.Normal, true, content)
//
// This last command within the `Draw` method paints the textual-content
// prepared earlier onto the canvas provided, center-justified, wrapping on
// word boundaries, using the default `Normal` theme, specifying that the
// content is in fact to be parsed as `Tango` markup and finally the content
// itself.
//
// The result of this drawing process should be a navy-blue screen, with a
// border, and three lines of text centered neatly. The three lines of text
// should be marked up with bold, italics, underlines and colorization. The
// last line of text should be telling the current time and date at the time
// of rendering.
//
// 		func (w *CdkDemoWindow) ProcessEvent(evt cdk.Event) cdk.EventFlag {
// 			w.LogInfo("ProcessEvent: %v", evt)
// 			return cdk.EVENT_STOP
// 		}
//
// The `ProcessEvent` method is the main event handler. Whenever a new event is
// received by a CDK `Display`, it is passed on to the active `Window` and in
// this demonstration, all that's happening is a log entry is being made which
// mentions the event received.
//
// When implementing your own `ProcessEvent` method, if the `Display` should
// repaint the screen for example, one would make two calls to methods on the
// `DisplayManager`:
//
// 		w.GetDisplayManager().RequestDraw()
//      w.GetDisplayManager().RequestShow()
//
// CDK is a multi-threaded framework and the various `Request*()` methods on the
// `DisplayManager` are used to funnel requests along the right channels in
// order to first render the canvas (via `Draw` calls on the active `Window`)
// and then follow that up with the running `Display` painting itself from the
// canvas modified in the `Draw` process.
//
// The other complete list of request methods is as follows:
//
// 		RequestDraw()  // window draw rendering
//      RequestShow()  // display rendering (modified cells only)
//      RequestSync()  // full display synchronization (all cells updated)
//      RequestQuit()  // graceful shutdown of a CDK `Display`
//
// This concludes the `CdkDemoWindow` type implementation. Now on to using it!
//
// The Main Func
//
// The `main()` function within the `_demos/cdk-demo.go` sources is deceptively
// simple to implement.
//
// 	app := cdk.NewApp(
// 		"cdk-demo",
// 		"An example of a formal CDK Application",
// 		"0.0.1",
// 		"demo",
// 		"CDK Demo",
// 		"/dev/tty",
//		func(d cdk.DisplayManager) error {
// 			cdk.DebugF("cdk-demo initFn hit")
//			d.CaptureCtrlC()
// 			w := &CdkDemoWindow{}
// 			w.Init()
// 			d.SetActiveWindow(w)
// 			cdk.AddTimeout(time.Second, func() cdk.EventFlag {
// 				d.RequestDraw()
// 				d.RequestShow()
// 				return cdk.EVENT_PASS // keep looping every second
// 			})
// 			return nil
// 		},
// 	}
//
// The bulk of the code is constructing a new CDK `App` instance. This object
// is a wrapper around the `github.com/urfave/cli/v2` CLI package, providing
// a tidy interface to managing CLI arguments, help documentation and so on.
// In this example, the `App` is configured with a bunch of metadata for:
// the program's name "cdk-demo", a simply usage summary, the current version
// number, an internally-used tag, a title for the main window and the display
// is to use the `/dev/tty` (default) terminal device.
//
// Beyond the metadata, the final argument is an initialization function. This
// function receives a fully instantiated and running `Display` instance and it
// is expected that the application instantiates it's `Window` and sets it as
// the active window for the given `Display`.
//
// In addition to that is one final call to `AddTimeout`. This call will trigger
// the given `func() cdk.EventFlag` once, after a second. Because the `func()`
// implemented here in this demonstration returns the `cdk.EVENT_PASS` flag
// it will be continually called once per second. For this demonstration, this
// implementation simply requests a draw and show cycle which will cause the
// screen to be repainted with the current date and time every second the demo
// application is running.
//
// 	if err := app.Run(os.Args); err != nil {
// 		panic(err)
// 	}
//
// The final bit of code in this CDK demonstration simply passes the arguments
// provided by the end-user on the command-line in a call to the `App`'s `Run()`
// method. This will cause the `DisplayManager`, `Display` and other systems to
// instantiate and begin processing events and render cycles.
//
package cdk
