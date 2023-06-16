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

package postgresql

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// BaseExecutor represents a base frontend message executor.
type BaseExecutor struct {
}

// NewBaseExecutor returns a base frontend message executor.
func NewBaseExecutor() *BaseExecutor {
	return &BaseExecutor{}
}

// Authenticate authenticates the connection with the startup message.
func (executor *BaseExecutor) Authenticate(*Conn, *message.Startup) bool {
	return true
}

// ParameterStatus returns the parameter status.
func (executor *BaseExecutor) ParameterStatus(*Conn) map[string]string {
	m := map[string]string{}
	m[message.ClientEncoding] = message.EncodingUTF8
	m[message.ServerEncoding] = message.EncodingUTF8
	return m
}
