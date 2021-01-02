// +build cdkdemo

package main

import (
	"os"
	"time"

	"github.com/kckrinke/go-cdk"
)

type CdkDemoWindow struct {
	cdk.CWindow
}

func (w *CdkDemoWindow) Init() (already bool) {
	if w.CWindow.Init() {
		return true
	}
	return false
}

func (w *CdkDemoWindow) Draw(canvas *cdk.Canvas) cdk.EventFlag {
	w.LogInfo("Draw: %s", canvas)
	theme := w.GetDisplay().DefaultTheme()
	size := canvas.GetSize()
	canvas.Box(cdk.Point2I{0, 0}, size, true, true, theme)
	content := "<b><u>H</u><span foreground=\"gold\">ello</span> <i>W</i><span foreground=\"brown\">orld</span></b>\n"
	content += "<span foreground=\"grey\" background=\"cyan\">(press CTRL+c to exit)</span>\n"
	content += "<span foreground=\"silver\" background=\"darkblue\">"+time.Now().Format("2006-01-02 15:04:05")+"</span>"
	textPoint := cdk.Point2I{
		X: size.W / 2 / 2,
		Y: size.H/2 - 1,
	}
	textSize := cdk.Rectangle{
		W: size.W / 2,
		H: size.H / 2,
	}
	canvas.DrawText(textPoint, textSize, cdk.JUSTIFY_CENTER, false, cdk.WRAP_WORD, cdk.DefaultColorTheme.Normal, true, content)
	return cdk.EVENT_STOP
}

func (w *CdkDemoWindow) ProcessEvent(evt cdk.Event) cdk.EventFlag {
	w.LogInfo("ProcessEvent: %v", evt)
	return cdk.EVENT_STOP
}

func main() {
	app := cdk.NewApp(
		"cdk-demo",
		"An example of a formal CDK Application",
		"0.0.1",
		"demo",
		"CDK Demo",
		"/dev/tty",
		func(d cdk.DisplayManager) error {
			cdk.DebugF("cdk-demo initFn hit")
			d.CaptureCtrlC()
			w := &CdkDemoWindow{}
			w.Init()
			d.SetActiveWindow(w)
			cdk.AddTimeout(time.Second, func() cdk.EventFlag {
				d.RequestDraw()
				d.RequestShow()
				return cdk.EVENT_PASS // keep looping every second
			})
			return nil
		},
	)
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
