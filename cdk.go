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

// cli, flags, timers, threading, logging

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/envy"
	"github.com/urfave/cli/v2"
)

type ScreenStateReq uint64

const (
	NullRequest ScreenStateReq = 1 << iota
	DrawRequest
	ShowRequest
	SyncRequest
	QuitRequest
)

type DisplayInitFn = func(d Display) error

type App interface {
	GetContext() *cli.Context
	Tag() string
	Title() string
	Name() string
	Display() Display
	CLI() *cli.App
	Version() string
	InitUI(c *cli.Context) error
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
	display *CDisplay
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
	app.display = NewDisplay(app.title, app.ttyPath)
	app.display.app = app
	app.cli = &cli.App{
		Name:    app.name,
		Usage:   app.usage,
		Version: app.version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cdk-log-file",
				EnvVars:     []string{"GO_CDK_LOG_FILE"},
				Value:       "",
				Usage:       "path to log file",
				DefaultText: DEFAULT_LOG_PATH,
			},
			&cli.StringFlag{
				Name:        "cdk-log-level",
				EnvVars:     []string{"GO_CDK_LOG_LEVEL"},
				Value:       "error",
				Usage:       "highest level of verbosity",
				DefaultText: "error",
			},
			&cli.StringFlag{
				Name:        "cdk-log-format",
				EnvVars:     []string{"GO_CDK_LOG_FORMAT"},
				Value:       "pretty",
				Usage:       "json, text or pretty",
				DefaultText: "pretty",
			},
			&cli.BoolFlag{
				Name:  "cdk-log-levels",
				Value: false,
				Usage: "list the levels of logging verbosity",
			},
		},
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

func (app *CApp) Display() Display {
	return app.display
}

func (app *CApp) CLI() *cli.App {
	return app.cli
}

func (app *CApp) Version() string {
	return app.version
}

func (app *CApp) InitUI(c *cli.Context) error {
	return app.initFn(app.Display())
}

func (app *CApp) AddFlag(f cli.Flag) {
	app.cli.Flags = append(app.cli.Flags, f)
}

func (app *CApp) AddCommand(c *cli.Command) {
	app.cli.Commands = append(app.cli.Commands, c)
}

func (app *CApp) Run(args []string) error {
	sort.Sort(cli.CommandsByName(app.cli.Commands))
	sort.Sort(CliFlagSorter(app.cli.Flags))
	return app.cli.Run(args)
}

func (app *CApp) MainActionFn(c *cli.Context) error {
	if c.Bool("ctk-log-levels") {
		for i := len(LOG_LEVELS) - 1; i >= 0; i-- {
			fmt.Printf("%s\n", LOG_LEVELS[i])
		}
		return nil
	}
	if v := c.String("cdk-log-file"); v != "" {
		envy.Set("GO_CDK_LOG_OUTPUT", "file")
		envy.Set("GO_CDK_LOG_FILE", v)
	}
	if v := c.String("cdk-log-level"); v != "" {
		envy.Set("GO_CDK_LOG_LEVEL", v)
	}
	ReloadLogging()
	defer StopLogging()
	app.context = c
	if err := app.InitUI(c); err != nil {
		return err
	}
	return app.Display().Run()
}
