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
	ejson "encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TestLoggingInit(t *testing.T) {
	Convey("Logging initialization checks", t, func() {
		// check for system event
		// test output methods: stdout, stderr, file, filepath, /dev/null
		logged, _, err := DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			ReloadLogging()
			ErrorF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(cdkLogger.Formatter, ShouldHaveSameTypeAs, &prefixed.TextFormatter{})
		So(logged, ShouldStartWith, "ERROR")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stderr")
			ReloadLogging()
			ErrorF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(cdkLogger.Formatter, ShouldHaveSameTypeAs, &prefixed.TextFormatter{})
		So(logged, ShouldStartWith, "ERROR")
		So(logged, ShouldEndWith, "testing\n")
	})
}

func TestLoggingTimestamps(t *testing.T) {
	Convey("Logging timestamp checks", t, func() {
		logged, _, err := DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_TIMESTAMPS", "true")
			os.Setenv("GO_CDK_LOG_TIMESTAMP_FORMAT", "2006-01-02")
			ReloadLogging()
			ErrorF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		datestamp := time.Now().Format("2006-01-02")
		So(logged, ShouldStartWith, "["+datestamp+"]")
		So(logged, ShouldEndWith, "testing\n")
		os.Setenv("GO_CDK_LOG_TIMESTAMPS", "")
		os.Setenv("GO_CDK_LOG_TIMESTAMP_FORMAT", "")
	})
}

func TestLoggingFormatter(t *testing.T) {
	Convey("Logging json formatter checks", t, func() {
		// test formatter settings
		logged, _, err := DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "json")
			ReloadLogging()
			ErrorF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(cdkLogger.Formatter, ShouldHaveSameTypeAs, &log.JSONFormatter{})
		decoded := make(map[string]interface{})
		err = ejson.Unmarshal([]byte(logged), &decoded)
		So(err, ShouldBeNil)
		So(decoded, ShouldNotBeEmpty)
		So(decoded["level"], ShouldHaveSameTypeAs, "")
		So(decoded["level"].(string), ShouldEqual, "error")
		So(decoded["msg"], ShouldHaveSameTypeAs, "")
		So(decoded["msg"].(string), ShouldEndWith, "testing")
	})
	Convey("Logging text formatter checks", t, func() {
		logged, _, err := DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "text")
			ReloadLogging()
			ErrorF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(cdkLogger.Formatter, ShouldHaveSameTypeAs, &log.TextFormatter{})
		So(logged, ShouldStartWith, "level=error")
		So(logged, ShouldEndWith, "testing\"\n")
	})
}

func TestLoggingLevel(t *testing.T) {
	Convey("Logging level checks", t, func() {
		logged, _, err := DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
			os.Setenv("GO_CDK_LOG_LEVEL", "trace")
			ReloadLogging()
			TraceF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, "TRACE")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
			os.Setenv("GO_CDK_LOG_LEVEL", "debug")
			ReloadLogging()
			TraceF("testing")
			DebugF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, "DEBUG")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
			os.Setenv("GO_CDK_LOG_LEVEL", "info")
			ReloadLogging()
			TraceF("testing")
			DebugF("testing")
			InfoF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, " INFO")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
			os.Setenv("GO_CDK_LOG_LEVEL", "warn")
			ReloadLogging()
			TraceF("testing")
			DebugF("testing")
			InfoF("testing")
			WarnF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, " WARN")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
			os.Setenv("GO_CDK_LOG_LEVEL", "error")
			ReloadLogging()
			TraceF("testing")
			DebugF("testing")
			InfoF("testing")
			WarnF("testing")
			ErrorF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, "ERROR")
		So(logged, ShouldEndWith, "testing\n")
		// fatal
		var fatal bool = false
		logged, _, err = DoWithFakeIO(func() error {
			os.Setenv("GO_CDK_LOG_OUTPUT", "stdout")
			os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
			os.Setenv("GO_CDK_LOG_LEVEL", "error")
			ReloadLogging()
			cdkLogger.ExitFunc = func(int) { fatal = true }
			FatalF("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(fatal, ShouldEqual, true)
		So(logged, ShouldStartWith, "FATAL")
		So(logged, ShouldEndWith, "testing\n")
		cdkLogger.ExitFunc = nil
		prefix := getLogPrefix(99)
		So(prefix, ShouldEqual, "(missing caller metadata)")
	})
}

func TestLoggingToFiles(t *testing.T) {
	Convey("Logging file checks", t, func() {
		So(cdkLogFH, ShouldBeNil)
		if _, err := os.Stat(DefaultLogPath); err == nil {
			os.Remove(DefaultLogPath)
		}
		os.Setenv("GO_CDK_LOG_OUTPUT", "file")
		os.Setenv("GO_CDK_LOG_FORMAT", "pretty")
		os.Setenv("GO_CDK_LOG_LEVEL", "error")
		os.Setenv("GO_CDK_LOG_FILE", DefaultLogPath)
		ReloadLogging()
		So(cdkLogFH, ShouldNotBeNil)
		ErrorF("testing")
		found_file := false
		if _, err := os.Stat(DefaultLogPath); err == nil {
			found_file = true
		}
		So(found_file, ShouldEqual, true)
		logged, err := ioutil.ReadFile(DefaultLogPath)
		So(err, ShouldBeNil)
		So(string(logged), ShouldEndWith, "testing\n")
		cdkLogFH.Close()
		os.Remove(DefaultLogPath)
		os.Setenv("GO_CDK_LOG_FILE", "/dev/null")
		err = ReloadLogging()
		So(err, ShouldBeNil)
		ErrorF("testing")
		found_file = false
		if _, err := os.Stat(DefaultLogPath); err == nil {
			found_file = true
		}
		So(found_file, ShouldEqual, false)
		tmp_log := os.TempDir() + string(os.PathSeparator) + "cdk.not.log"
		os.Remove(tmp_log)
		os.Setenv("GO_CDK_LOG_FILE", tmp_log)
		ReloadLogging()
		ErrorF("testing")
		found_file = false
		if _, err := os.Stat(tmp_log); err == nil {
			found_file = true
		}
		So(found_file, ShouldEqual, true)
		logged, err = ioutil.ReadFile(tmp_log)
		So(err, ShouldBeNil)
		So(string(logged), ShouldEndWith, "testing\n")
		So(cdkLogFH, ShouldNotBeNil)
		StopLogging()
		So(cdkLogFH, ShouldBeNil)
		os.Chmod(tmp_log, 0000)
		err = ReloadLogging()
		So(err, ShouldNotBeNil)
		os.Chmod(tmp_log, 0660)
		os.Remove(tmp_log)
	})
}
