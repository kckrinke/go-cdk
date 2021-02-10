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
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/kckrinke/go-cdk/utils"
)

const (
	LevelError string = "error"
	LevelWarn  string = "warn"
	LevelInfo  string = "info"
	LevelDebug string = "debug"
	LevelTrace string = "trace"
)

var LogLevels = []string{
	LevelError,
	LevelWarn,
	LevelInfo,
	LevelDebug,
	LevelTrace,
}

const (
	FormatPretty string = "pretty"
	FormatText   string = "text"
	FormatJson   string = "json"
)

const (
	OutputStderr string = "stderr"
	OutputStdout string = "stdout"
	OutputFile   string = "file"
)

var (
	cdkLogger       = log.New()
	cdkLogFH        *os.File
	cdkLogFullPaths = false
	DefaultLogPath  = os.TempDir() + string(os.PathSeparator) + "cdk.log"
)

const (
	DefaultLogTimestampFormat = "2006-01-02T15:04:05.000"
)

func ReloadLogging() error {
	envy.Reload()
	disableTimestamp := true
	if v := envy.Get("GO_CDK_LOG_TIMESTAMPS", "false"); v == "true" {
		disableTimestamp = false
	}
	timestampFormat := DefaultLogTimestampFormat
	if v := envy.Get("GO_CDK_LOG_TIMESTAMP_FORMAT", ""); v != "" {
		timestampFormat = v
	}
	switch envy.Get("GO_CDK_LOG_FULL_PATHS", "false") {
	case "true":
		cdkLogFullPaths = true
	default:
		cdkLogFullPaths = false
	}
	switch envy.Get("GO_CDK_LOG_FORMAT", "pretty") {
	case FormatJson:
		cdkLogger.SetFormatter(&log.JSONFormatter{
			TimestampFormat:  timestampFormat,
			DisableTimestamp: disableTimestamp,
		})
	case FormatText:
		cdkLogger.SetFormatter(&log.TextFormatter{
			TimestampFormat:  timestampFormat,
			DisableTimestamp: disableTimestamp,
			DisableSorting:   true,
			DisableColors:    true,
			FullTimestamp:    true,
		})
	case FormatPretty:
		fallthrough
	default:
		cdkLogger.SetFormatter(&prefixed.TextFormatter{
			DisableTimestamp: disableTimestamp,
			TimestampFormat:  timestampFormat,
			ForceFormatting:  true,
			FullTimestamp:    true,
			DisableSorting:   true,
			DisableColors:    true,
		})
	}
	switch envy.Get("GO_CDK_LOG_LEVEL", LevelError) {
	case LevelTrace:
		cdkLogger.SetLevel(log.TraceLevel)
	case LevelDebug:
		cdkLogger.SetLevel(log.DebugLevel)
	case LevelInfo:
		cdkLogger.SetLevel(log.InfoLevel)
	case LevelWarn:
		cdkLogger.SetLevel(log.WarnLevel)
	case LevelError:
		fallthrough
	default:
		cdkLogger.SetLevel(log.ErrorLevel)
	}
	switch envy.Get("GO_CDK_LOG_OUTPUT", OutputFile) {
	case OutputStdout:
		cdkLogger.SetOutput(os.Stdout)
	case OutputStderr:
		cdkLogger.SetOutput(os.Stderr)
	case OutputFile:
		fallthrough
	default:
		_ = StopLogging()
		if logfile := envy.Get("GO_CDK_LOG_FILE", DefaultLogPath); !utils.IsEmpty(logfile) && logfile != "/dev/null" {
			logFH, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}
			cdkLogFH = logFH
			cdkLogger.SetOutput(cdkLogFH)
		} else {
			cdkLogger.SetOutput(ioutil.Discard)
		}
	}
	return nil
}

func StopLogging() error {
	if cdkLogFH != nil {
		if err := cdkLogFH.Close(); err != nil {
			return err
		}
		cdkLogFH = nil
	}
	return nil
}

func getLogPrefix(depth int) string {
	depth += 1
	if function, file, line, ok := runtime.Caller(depth); ok {
		fullName := runtime.FuncForPC(function).Name()
		funcName := fullName
		if i := strings.LastIndex(fullName, "."); i > -1 {
			funcName = fullName[i+1:]
		}
		funcName = utils.PadLeft(funcName, " ", 12)
		packName := fullName
		if i := strings.Index(fullName, "."); i > -1 {
			packName = fullName[:i+1]
		}
		path := file
		if !cdkLogFullPaths {
			if i := strings.Index(path, packName); i > -1 {
				path = file[i:]
			}
		}
		return fmt.Sprintf("%s:%d	%s", path, line, funcName)
	}
	return "(missing caller metadata)"
}

func TraceF(format string, argv ...interface{}) { TraceDF(1, format, argv...) }
func TraceDF(depth int, format string, argv ...interface{}) {
	cdkLogger.Tracef(utils.NLSprintf("%s	%s", getLogPrefix(depth+1), format), argv...)
}

func DebugF(format string, argv ...interface{}) { DebugDF(1, format, argv...) }
func DebugDF(depth int, format string, argv ...interface{}) {
	cdkLogger.Debugf(utils.NLSprintf("%s	%s", getLogPrefix(depth+1), format), argv...)
}

func InfoF(format string, argv ...interface{}) { InfoDF(1, format, argv...) }
func InfoDF(depth int, format string, argv ...interface{}) {
	cdkLogger.Infof(utils.NLSprintf("%s	%s", getLogPrefix(depth+1), format), argv...)
}

func WarnF(format string, argv ...interface{}) { WarnDF(1, format, argv...) }
func WarnDF(depth int, format string, argv ...interface{}) {
	cdkLogger.Warnf(utils.NLSprintf("%s	%s", getLogPrefix(depth+1), format), argv...)
}

func Error(err error)                           { ErrorDF(1, err.Error()) }
func ErrorF(format string, argv ...interface{}) { ErrorDF(1, format, argv...) }
func ErrorDF(depth int, format string, argv ...interface{}) {
	cdkLogger.Errorf(utils.NLSprintf("%s	%s", getLogPrefix(depth+1), format), argv...)
}

func Fatal(err error)                           { FatalDF(1, err.Error()) }
func FatalF(format string, argv ...interface{}) { FatalDF(1, format, argv...) }
func FatalDF(depth int, format string, argv ...interface{}) {
	if dm := GetDisplayManager(); dm != nil {
		dm.ReleaseDisplay()
	}
	message := fmt.Sprintf(utils.NLSprintf("%s\t%s", getLogPrefix(depth+1), format), argv...)
	cdkLogger.Fatalf(message)
}

func Panic(err error)                           { PanicDF(1, err.Error()) }
func PanicF(format string, argv ...interface{}) { PanicDF(1, format, argv...) }
func PanicDF(depth int, format string, argv ...interface{}) {
	if dm := GetDisplayManager(); dm != nil {
		dm.ReleaseDisplay()
	}
	message := fmt.Sprintf(utils.NLSprintf("%s\t%s", getLogPrefix(depth+1), format), argv...)
	cdkLogger.Errorf(message)
	_ = StopLogging()
	panic(message)
}

func Exit(code int) {
	InfoDF(1, "exiting with code: %d", code)
	cdkLogger.Exit(code)
}
