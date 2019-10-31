// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"time"
)

type loggerOutpter func(file string, level Level, msg string) (int, error)

// Logger represents a logger interface.
type Logger struct {
	file     string
	level    Level
	outputer loggerOutpter
}

const (
	filePerm                 = 0644
	loggerLevelUnknownString = "UNKNOWN"
	loggerStdout             = "stdout"
)

var levelStrings = map[Level]string{
	TRACE: "TRACE",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

func getLevelString(Level Level) string {
	logString, hasString := levelStrings[Level]
	if !hasString {
		return loggerLevelUnknownString
	}
	return logString
}

// SetLevel sets a logging level into all outputs.
func (logger *Logger) SetLevel(level Level) {
	logger.level = level
}

// Level returns the current logging level for all outputs.
func (logger *Logger) Level() Level {
	return logger.level
}

func output(outputLevel Level, msgFormat string, msgArgs ...interface{}) int {
	if sharedLogger == nil {
		return 0
	}

	Level := sharedLogger.Level()
	if (Level < outputLevel) || (Level <= FATAL) || (TRACE < Level) {
		return 0
	}

	t := time.Now()
	logDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	headerString := fmt.Sprintf("[%s]", getLevelString(outputLevel))
	logMsg := fmt.Sprintf("%s %s %s", logDate, headerString, fmt.Sprintf(msgFormat, msgArgs...))
	logMsgLen := len(logMsg)

	if 0 < logMsgLen {
		logMsgLen, _ = sharedLogger.outputer(sharedLogger.file, Level, logMsg)
	}

	return logMsgLen
}
