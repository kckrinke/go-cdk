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
	"os"
	"sort"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/pkg/profile"
	"github.com/urfave/cli/v2"

	"github.com/kckrinke/go-cdk/utils"
)

type ScreenStateReq uint64

const (
	NullRequest ScreenStateReq = 1 << iota
	DrawRequest
	ShowRequest
	SyncRequest
	QuitRequest
)

var (
	DefaultGoProfilePath = os.TempDir() + string(os.PathSeparator) + "cdk.pprof"
)

type goProfileFn = func(p *profile.Profile)

type DisplayInitFn = func(d DisplayManager) error

type App interface {
	GetContext() *cli.Context
	Tag() string
	Title() string
	Name() string
	Usage() string
	DisplayManager() DisplayManager
	CLI() *cli.App
	Version() string
	InitUI() error
	AddFlag(f cli.Flag)
	AddCommand(c *cli.Command)
	Run(args []string) error
	MainActionFn(c *cli.Context) error
}

type CApp struct {
	name    string
	usage   string
	version string
	tag     string
	title   string
	ttyPath string
	display *CDisplayManager
	context *cli.Context
	cli     *cli.App
	initFn  DisplayInitFn
	valid   bool
}

func NewApp(name, usage, version, tag, title, ttyPath string, initFn DisplayInitFn) *CApp {
	app := &CApp{
		name:    name,
		usage:   usage,
		version: version,
		tag:     tag,
		title:   title,
		ttyPath: ttyPath,
		initFn:  initFn,
	}
	app.init()
	return app
}

func (app *CApp) init() {
	app.display = NewDisplayManager(app.title, app.ttyPath)
	app.display.app = app
	app.cli = &cli.App{
		Name:     app.name,
		Usage:    app.usage,
		Version:  app.version,
		Flags:    getCdkCliFlags(),
		Commands: []*cli.Command{},
		Action:   app.MainActionFn,
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{},
		Usage:   "display the version",
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h", "usage"},
		Usage:   "display command-line usage information",
	}
	app.valid = true
	return
}

func (app *CApp) Destroy() {
	if app.display != nil {
		app.display.ReleaseDisplay()
		app.display.Destroy()
	}
	app.display = nil
	app.context = nil
	app.cli = nil
	app.valid = false
}

func (app *CApp) GetContext() *cli.Context {
	return app.context
}

func (app *CApp) Tag() string {
	return app.tag
}

func (app *CApp) Title() string {
	return app.title
}

func (app *CApp) Name() string {
	return app.name
}

func (app *CApp) Usage() string {
	return app.usage
}

func (app *CApp) DisplayManager() DisplayManager {
	return app.display
}

func (app *CApp) CLI() *cli.App {
	return app.cli
}

func (app *CApp) Version() string {
	return app.version
}

func (app *CApp) InitUI() error {
	return app.initFn(app.DisplayManager())
}

func (app *CApp) AddFlag(flag cli.Flag) {
	app.cli.Flags = append(app.cli.Flags, flag)
}

func (app *CApp) AddFlags(flags []cli.Flag) {
	for _, f := range flags {
		app.AddFlag(f)
	}
}

func (app *CApp) AddCommand(command *cli.Command) {
	app.cli.Commands = append(app.cli.Commands, command)
}

func (app *CApp) AddCommands(commands []*cli.Command) {
	for _, c := range commands {
		app.AddCommand(c)
	}
}

func (app *CApp) Run(args []string) error {
	sort.Sort(cli.CommandsByName(app.cli.Commands))
	sort.Sort(FlagSorter(app.cli.Flags))
	return app.cli.Run(args)
}

func (app *CApp) MainActionFn(c *cli.Context) error {
	app.context = c
	if Build.LogLevel {
		if v := c.String("cdk-log-level"); !utils.IsEmpty(v) {
			envy.Set("GO_CDK_LOG_LEVEL", v)
		}
		if Build.LogLevels {
			if c.Bool("ctk-log-levels") {
				for i := len(LogLevels) - 1; i >= 0; i-- {
					fmt.Printf("%s\n", LogLevels[i])
				}
				return nil
			}
		}
	}
	if Build.LogFile {
		if v := c.String("cdk-log-file"); !utils.IsEmpty(v) {
			envy.Set("GO_CDK_LOG_OUTPUT", "file")
			envy.Set("GO_CDK_LOG_FILE", v)
		}
	}
	profilePath := DefaultGoProfilePath
	if Build.Profiling {
		if v := c.String("cdk-profile-path"); !utils.IsEmpty(v) {
			if !utils.IsDir(v) {
				if err := utils.MakeDir(v, 0770); err != nil {
					Fatal(err)
				}
			}
			envy.Set("GO_CDK_PROFILE_PATH", v)
			profilePath = v
		}
	}
	_ = ReloadLogging()
	defer func() { _ = StopLogging() }()
	if Build.Profiling {
		if v := c.String("cdk-profile"); !utils.IsEmpty(v) {
			v = strings.ToLower(v)
			var p goProfileFn
			envy.Set("GO_CDK_PROFILE", v)
			// none, block, cpu, goroutine, mem, mutex, thread or trace
			switch v {
			case "block":
				p = profile.BlockProfile
			case "cpu":
				p = profile.CPUProfile
			case "goroutine":
				p = profile.GoroutineProfile
			case "mem":
				p = profile.MemProfile
			case "mutex":
				p = profile.MutexProfile
			case "thread":
				p = profile.ThreadcreationProfile
			case "trace":
				p = profile.TraceProfile
			default:
				p = nil
			}
			if p != nil {
				DebugF("starting profile of \"%v\" to path: %v", v, profilePath)
				defer profile.Start(p, profile.ProfilePath(profilePath)).Stop()
			}
		}
	}
	return app.DisplayManager().Run()
}
