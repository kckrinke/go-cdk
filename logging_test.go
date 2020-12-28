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

	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/kckrinke/go-cdk/utils"
)

func TestLoggingInit(t *testing.T) {
	Convey("Logging initialization checks", t, func() {
		// check for system event
		So(_cdk_system_events[SYSTEM_EVENT_BOOT]["cdk.logging.boot"], ShouldNotBeNil)
		So(_cdk_system_events[SYSTEM_EVENT_SHUTDOWN]["cdk.logging.shutdown"], ShouldNotBeNil)
		// test output methods: stdout, stderr, file, filepath, /dev/null
		logged, _, err := DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			ReloadLogging()
			Errorf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(_cdk_logger.Formatter, ShouldHaveSameTypeAs, &prefixed.TextFormatter{})
		So(logged, ShouldStartWith, "ERROR")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stderr")
			ReloadLogging()
			Errorf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(_cdk_logger.Formatter, ShouldHaveSameTypeAs, &prefixed.TextFormatter{})
		So(logged, ShouldStartWith, "ERROR")
		So(logged, ShouldEndWith, "testing\n")
	})
}

func TestLoggingTimestamps(t *testing.T) {
	Convey("Logging timestamp checks", t, func() {
		logged, _, err := DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_TIMESTAMPS", "true")
			envy.Set("GO_CDK_LOG_TIMESTAMP_FORMAT", "2006-01-02")
			ReloadLogging()
			Errorf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		datestamp := time.Now().Format("2006-01-02")
		So(logged, ShouldStartWith, "["+datestamp+"]")
		So(logged, ShouldEndWith, "testing\n")
		envy.Set("GO_CDK_LOG_TIMESTAMPS", "")
		envy.Set("GO_CDK_LOG_TIMESTAMP_FORMAT", "")
	})
}

func TestLoggingFormatter(t *testing.T) {
	Convey("Logging json formatter checks", t, func() {
		// test formatter settings
		logged, _, err := DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "json")
			ReloadLogging()
			Errorf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(_cdk_logger.Formatter, ShouldHaveSameTypeAs, &log.JSONFormatter{})
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
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "text")
			ReloadLogging()
			Errorf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(_cdk_logger.Formatter, ShouldHaveSameTypeAs, &log.TextFormatter{})
		So(logged, ShouldStartWith, "level=error")
		So(logged, ShouldEndWith, "testing\"\n")
	})
}

func TestLoggingLevel(t *testing.T) {
	Convey("Logging level checks", t, func() {
		logged, _, err := DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "trace")
			ReloadLogging()
			Tracef("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, "TRACE")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "debug")
			ReloadLogging()
			Tracef("testing")
			Debugf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, "DEBUG")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "info")
			ReloadLogging()
			Tracef("testing")
			Debugf("testing")
			Infof("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, " INFO")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "warn")
			ReloadLogging()
			Tracef("testing")
			Debugf("testing")
			Infof("testing")
			Warnf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, " WARN")
		So(logged, ShouldEndWith, "testing\n")
		logged, _, err = DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "error")
			ReloadLogging()
			Tracef("testing")
			Debugf("testing")
			Infof("testing")
			Warnf("testing")
			Errorf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(logged, ShouldStartWith, "ERROR")
		So(logged, ShouldEndWith, "testing\n")
		// fatal
		var fatal bool = false
		logged, _, err = DoWithFakeIO(func() error {
			envy.Set("GO_CDK_LOG_OUTPUT", "stdout")
			envy.Set("GO_CDK_LOG_FORMAT", "pretty")
			envy.Set("GO_CDK_LOG_LEVEL", "error")
			ReloadLogging()
			_cdk_logger.ExitFunc = func(int) { fatal = true }
			Fatalf("testing")
			return nil
		})
		So(err, ShouldBeNil)
		So(fatal, ShouldEqual, true)
		So(logged, ShouldStartWith, "FATAL")
		So(logged, ShouldEndWith, "testing\n")
		_cdk_logger.ExitFunc = nil
		prefix := get_log_prefix(99)
		So(prefix, ShouldEqual, "(missing caller metadata)")
	})
}

func TestLoggingToFiles(t *testing.T) {
	Convey("Logging file checks", t, func() {
		So(_cdk_logfh, ShouldBeNil)
		if _, err := os.Stat(DEFAULT_LOG_PATH); err == nil {
			os.Remove(DEFAULT_LOG_PATH)
		}
		envy.Set("GO_CDK_LOG_OUTPUT", "file")
		envy.Set("GO_CDK_LOG_FORMAT", "pretty")
		envy.Set("GO_CDK_LOG_LEVEL", "error")
		ReloadLogging()
		Errorf("testing")
		So(_cdk_logfh, ShouldNotBeNil)
		found_file := false
		if _, err := os.Stat(DEFAULT_LOG_PATH); err == nil {
			found_file = true
		}
		So(found_file, ShouldEqual, true)
		logged, err := ioutil.ReadFile(DEFAULT_LOG_PATH)
		So(err, ShouldBeNil)
		So(string(logged), ShouldEndWith, "testing\n")
		_cdk_logfh.Close()
		os.Remove(DEFAULT_LOG_PATH)
		envy.Set("GO_CDK_LOG_FILE", "/dev/null")
		err = ReloadLogging()
		So(err, ShouldBeNil)
		Errorf("testing")
		found_file = false
		if _, err := os.Stat(DEFAULT_LOG_PATH); err == nil {
			found_file = true
		}
		So(found_file, ShouldEqual, false)
		tmp_log := os.TempDir() + string(os.PathSeparator) + "cdk.not.log"
		os.Remove(tmp_log)
		envy.Set("GO_CDK_LOG_FILE", tmp_log)
		ReloadLogging()
		Errorf("testing")
		found_file = false
		if _, err := os.Stat(tmp_log); err == nil {
			found_file = true
		}
		So(found_file, ShouldEqual, true)
		logged, err = ioutil.ReadFile(tmp_log)
		So(err, ShouldBeNil)
		So(string(logged), ShouldEndWith, "testing\n")
		So(_cdk_logfh, ShouldNotBeNil)
		HandleSystemEvent(SYSTEM_EVENT_SHUTDOWN)
		So(_cdk_logfh, ShouldBeNil)
		os.Chmod(tmp_log, 0000)
		err = ReloadLogging()
		So(err, ShouldNotBeNil)
		os.Chmod(tmp_log, 0660)
		os.Remove(tmp_log)
	})
}
