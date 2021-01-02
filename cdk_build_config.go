package cdk

import (
	"github.com/urfave/cli/v2"
)

type Config struct {
	Profiling          bool
	LogFile            bool
	LogFormat          bool
	LogFullPaths       bool
	LogLevel           bool
	LogLevels          bool
	LogTimestamps      bool
	LogTimestampFormat bool
	LogOutput          bool
}

var Build = Config{
	LogFile:   true,
	LogLevel:  true,
	LogLevels: true,
}

func getCdkCliFlags() (flags []cli.Flag) {
	if Build.Profiling {
		flags = append(flags, cdkProfileFlag, cdkProfilePathFlag)
	}
	if Build.LogFile {
		flags = append(flags, cdkLogFileFlag)
	}
	if Build.LogFormat {
		flags = append(flags, cdkLogFormatFlag)
	}
	if Build.LogFullPaths {
		flags = append(flags, cdkLogFullPathsFlag)
	}
	if Build.LogLevel {
		flags = append(flags, cdkLogLevel)
	}
	if Build.LogLevels {
		flags = append(flags, cdkLogLevelsFlag)
	}
	if Build.LogTimestampFormat {
		flags = append(flags, cdkLogTimestampFormatFlag)
	}
	if Build.LogTimestamps {
		flags = append(flags, cdkLogTimestampsFlag)
	}
	if Build.LogOutput {
		flags = append(flags, cdkLogOutputFlag)
	}
	return
}
