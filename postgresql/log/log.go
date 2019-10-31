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

const (
	// TRACE outputs all logging messages.
	TRACE = (1 << 5)
	// INFO outputs information or higher level logging messages.
	INFO = (1 << 4)
	// WARN outputs warning or higher level logging messages.
	WARN = (1 << 3)
	// ERROR outputs error or higher level logging messages.
	ERROR = (1 << 2)
	// FATAL outputs only fatal level logging messages.
	FATAL = (1 << 1)
)

// Level represens a logging level
type Level int

// Trace sends the specified trace message to the current outputs.
func Trace(format string, args ...interface{}) int {
	return output(TRACE, format, args...)
}

// Info sends the information trace message to the current outputs.
func Info(format string, args ...interface{}) int {
	return output(INFO, format, args...)
}

// Warn sends the warning trace message to the current outputs.
func Warn(format string, args ...interface{}) int {
	return output(WARN, format, args...)
}

// Error sends the specified error message to the current outputs.
func Error(format string, args ...interface{}) int {
	return output(ERROR, format, args...)
}

// Fatal sends the specified fatal message to the current outputs.
func Fatal(format string, args ...interface{}) int {
	return output(FATAL, format, args...)
}

// Output sends the specified level message to the current outputs.
func Output(outputLevel Level, format string, args ...interface{}) int {
	return output(outputLevel, format, args...)
}
