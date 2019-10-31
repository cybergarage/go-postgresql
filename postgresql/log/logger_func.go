// Copyright (C) 2019 Satoshi Konno. All rights reserved.
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

var sharedLogger *Logger

// SetSharedLogger sets a shard global logger.
func SetSharedLogger(logger *Logger) {
	sharedLogger = logger
}

// GetSharedLogger returns a shard global logger.
func GetSharedLogger() *Logger {
	return sharedLogger
}

// SetStdoutDebugEnbled sets a trace stdout logger for debug.
func SetStdoutDebugEnbled(flag bool) {
	if flag {
		SetSharedLogger(NewStdoutLogger(TRACE))
	} else {
		SetSharedLogger(nil)
	}
}
