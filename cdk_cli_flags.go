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
	"github.com/urfave/cli/v2"
)

var (
	cdkProfileFlag = &cli.StringFlag{
		Name:        "cdk-profile",
		EnvVars:     []string{"GO_CDK_PROFILE"},
		Value:       "",
		Usage:       "profile one of: none, block, cpu, goroutine, mem, mutex, thread or trace",
		DefaultText: "none",
	}
	cdkProfilePathFlag = &cli.StringFlag{
		Name:        "cdk-profile-path",
		EnvVars:     []string{"GO_CDK_PROFILE_PATH"},
		Value:       "",
		Usage:       "specify the directory path to store the profile data",
		DefaultText: DefaultGoProfilePath,
	}
	cdkLogFileFlag = &cli.StringFlag{
		Name:        "cdk-log-file",
		EnvVars:     []string{"GO_CDK_LOG_FILE"},
		Value:       "",
		Usage:       "path to log file",
		DefaultText: DefaultLogPath,
	}
	cdkLogLevel = &cli.StringFlag{
		Name:        "cdk-log-level",
		EnvVars:     []string{"GO_CDK_LOG_LEVEL"},
		Value:       "error",
		Usage:       "highest level of verbosity",
		DefaultText: "error",
	}
	cdkLogFormatFlag = &cli.StringFlag{
		Name:        "cdk-log-format",
		EnvVars:     []string{"GO_CDK_LOG_FORMAT"},
		Value:       "pretty",
		Usage:       "json, text or pretty",
		DefaultText: "pretty",
	}
	cdkLogTimestampsFlag = &cli.BoolFlag{
		Name:        "cdk-log-timestamps",
		EnvVars:     []string{"GO_CDK_LOG_TIMESTAMPS"},
		Value:       false,
		Usage:       "enable timestamps",
		DefaultText: "false",
	}
	cdkLogTimestampFormatFlag = &cli.StringFlag{
		Name:        "cdk-log-timestamp-format",
		EnvVars:     []string{"GO_CDK_LOG_TIMESTAMP_FORMAT"},
		Value:       DefaultLogTimestampFormat,
		Usage:       "timestamp format",
		DefaultText: DefaultLogTimestampFormat,
	}
	cdkLogFullPathsFlag = &cli.BoolFlag{
		Name:        "cdk-log-full-paths",
		EnvVars:     []string{"GO_CDK_LOG_FULL_PATHS"},
		Value:       false,
		Usage:       "log the full paths of source files",
		DefaultText: "false",
	}
	cdkLogOutputFlag = &cli.StringFlag{
		Name:        "cdk-log-output",
		EnvVars:     []string{"GO_CDK_LOG_OUTPUT"},
		Value:       "file",
		Usage:       "logging output type: stdout, stderr or file",
		DefaultText: "file",
	}
	cdkLogLevelsFlag = &cli.BoolFlag{
		Name:  "cdk-log-levels",
		Value: false,
		Usage: "list the levels of logging verbosity",
	}
)
