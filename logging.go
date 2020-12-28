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
	LEVEL_ERROR string = "error"
	LEVEL_WARN  string = "warn"
	LEVEL_INFO  string = "info"
	LEVEL_DEBUG string = "debug"
	LEVEL_TRACE string = "trace"
)

var LOG_LEVELS = []string{
	LEVEL_ERROR,
	LEVEL_WARN,
	LEVEL_INFO,
	LEVEL_DEBUG,
	LEVEL_TRACE,
}

const (
	FORMAT_PRETTY string = "pretty"
	FORMAT_TEXT   string = "text"
	FORMAT_JSON   string = "json"
)

const (
	OUTPUT_STDERR string = "stderr"
	OUTPUT_STDOUT string = "stdout"
	OUTPUT_FILE   string = "file"
)

var (
	_cdk_logger        *log.Logger = log.New()
	_cdk_logfh         *os.File
	_cdk_log_fullpaths bool = false
	DEFAULT_LOG_PATH        = os.TempDir() + string(os.PathSeparator) + "cdk.log"
)

func ReloadLogging() error {
	disable_timestamp := true
	if v := envy.Get("GO_CDK_LOG_TIMESTAMPS", "false"); v == "true" {
		disable_timestamp = false
	}
	timestamp_format := "2006-01-02T15:04:05.000"
	if v := envy.Get("GO_CDK_LOG_TIMESTAMP_FORMAT", ""); v != "" {
		timestamp_format = v
	}
	switch envy.Get("GO_CDK_LOG_FULLPATHS", "false") {
	case "true":
		_cdk_log_fullpaths = true
	default:
		_cdk_log_fullpaths = false
	}
	switch envy.Get("GO_CDK_LOG_FORMAT", "pretty") {
	case FORMAT_JSON:
		_cdk_logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat:  timestamp_format,
			DisableTimestamp: disable_timestamp,
		})
	case FORMAT_TEXT:
		_cdk_logger.SetFormatter(&log.TextFormatter{
			TimestampFormat:  timestamp_format,
			DisableTimestamp: disable_timestamp,
			DisableSorting:   true,
			DisableColors:    true,
			FullTimestamp:    true,
		})
	case FORMAT_PRETTY:
		fallthrough
	default:
		_cdk_logger.SetFormatter(&prefixed.TextFormatter{
			DisableTimestamp: disable_timestamp,
			TimestampFormat:  timestamp_format,
			ForceFormatting:  true,
			FullTimestamp:    true,
			DisableSorting:   true,
			DisableColors:    true,
		})
	}
	switch envy.Get("GO_CDK_LOG_LEVEL", LEVEL_ERROR) {
	case LEVEL_TRACE:
		_cdk_logger.SetLevel(log.TraceLevel)
	case LEVEL_DEBUG:
		_cdk_logger.SetLevel(log.DebugLevel)
	case LEVEL_INFO:
		_cdk_logger.SetLevel(log.InfoLevel)
	case LEVEL_WARN:
		_cdk_logger.SetLevel(log.WarnLevel)
	case LEVEL_ERROR:
		fallthrough
	default:
		_cdk_logger.SetLevel(log.ErrorLevel)
	}
	switch envy.Get("GO_CDK_LOG_OUTPUT", OUTPUT_FILE) {
	case OUTPUT_STDOUT:
		_cdk_logger.SetOutput(os.Stdout)
	case OUTPUT_STDERR:
		_cdk_logger.SetOutput(os.Stderr)
	case OUTPUT_FILE:
		fallthrough
	default:
		StopLogging()
		if logfile := envy.Get("GO_CDK_LOG_FILE", DEFAULT_LOG_PATH); !utils.IsEmpty(logfile) && logfile != "/dev/null" {
			logfh, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}
			_cdk_logfh = logfh
			_cdk_logger.SetOutput(_cdk_logfh)
		} else {
			_cdk_logger.SetOutput(ioutil.Discard)
		}
	}
	return nil
}

func StopLogging() error {
	if _cdk_logfh != nil {
		_cdk_logfh.Close()
		_cdk_logfh = nil
	}
	return nil
}

func get_log_prefix(depth int) string {
	depth += 1
	if function, file, line, ok := runtime.Caller(depth); ok {
		full_name := runtime.FuncForPC(function).Name()
		func_name := full_name
		if i := strings.LastIndex(full_name, "."); i > -1 {
			func_name = full_name[i+1:]
		}
		func_name = utils.PadLeft(func_name, " ", 12)
		pack_name := full_name
		if i := strings.Index(full_name, "."); i > -1 {
			pack_name = full_name[:i+1]
		}
		path := file
		if !_cdk_log_fullpaths {
			if i := strings.Index(path, pack_name); i > -1 {
				path = file[i:]
			}
		}
		return fmt.Sprintf("%s:%d	%s", path, line, func_name)
	}
	return "(missing caller metadata)"
}

func Tracef(format string, argv ...interface{}) { Tracedf(1, format, argv...) }
func Tracedf(depth int, format string, argv ...interface{}) {
	_cdk_logger.Tracef(utils.NLSprintf("%s	%s", get_log_prefix(depth+1), format), argv...)
}

func Debugf(format string, argv ...interface{}) { Debugdf(1, format, argv...) }
func Debugdf(depth int, format string, argv ...interface{}) {
	_cdk_logger.Debugf(utils.NLSprintf("%s	%s", get_log_prefix(depth+1), format), argv...)
}

func Infof(format string, argv ...interface{}) { Infodf(1, format, argv...) }
func Infodf(depth int, format string, argv ...interface{}) {
	_cdk_logger.Infof(utils.NLSprintf("%s	%s", get_log_prefix(depth+1), format), argv...)
}

func Warnf(format string, argv ...interface{}) { Warndf(1, format, argv...) }
func Warndf(depth int, format string, argv ...interface{}) {
	_cdk_logger.Warnf(utils.NLSprintf("%s	%s", get_log_prefix(depth+1), format), argv...)
}

func Error(err error)                           { Errordf(1, err.Error()) }
func Errorf(format string, argv ...interface{}) { Errordf(1, format, argv...) }
func Errordf(depth int, format string, argv ...interface{}) {
	_cdk_logger.Errorf(utils.NLSprintf("%s	%s", get_log_prefix(depth+1), format), argv...)
}

func Fatal(err error)                           { Fataldf(1, err.Error()) }
func Fatalf(format string, argv ...interface{}) { Fataldf(1, format, argv...) }
func Fataldf(depth int, format string, argv ...interface{}) {
	_cdk_logger.Fatalf(utils.NLSprintf("%s	%s", get_log_prefix(depth+1), format), argv...)
}

func Exit(code int) {
	_cdk_logger.Exit(code)
}
