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
	"crypto/rand"
	"math"
	"math/big"
	"os"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// BaseProtocolExecutor represents a base frontend message executor.
type BaseProtocolExecutor struct {
	processID int32
	secretKey int32
}

// NewBaseProtocolExecutor returns a base frontend message executor.
func NewBaseProtocolExecutor() *BaseProtocolExecutor {
	r, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		log.Error(err)
	}
	return &BaseProtocolExecutor{
		processID: int32(os.Getpid()),
		secretKey: int32(r.Int64()),
	}
}

// Authenticate authenticates the connection with the startup message.
func (executor *BaseProtocolExecutor) Authenticate(*Conn, *message.Startup) (message.Response, error) {
	return message.NewAuthenticationOk()
}

// ParameterStatus returns the parameter status.
func (executor *BaseProtocolExecutor) ParameterStatus(*Conn) (message.Response, error) {
	m := map[string]string{}
	m[message.ClientEncoding] = message.EncodingUTF8
	m[message.ServerEncoding] = message.EncodingUTF8
	// FIXME : Get the time zone name from the system
	// m[message.TimeZone] = time.Now().Location().String()
	return message.NewParameterStatusWith(m)
}

// BackendKeyData returns the backend key data.
func (executor *BaseProtocolExecutor) BackendKeyData(*Conn) (message.Response, error) {
	return message.NewBackendKeyDataWith(executor.processID, executor.secretKey)
}

// Parse returns the parse response.
func (executor *BaseProtocolExecutor) Parse(*Conn, *message.Parse) (message.Response, error) {
	return message.NewParseComplete(), nil
}

// Bind returns the bind response.
func (executor *BaseProtocolExecutor) Bind(*Conn, *message.Parse, *message.Bind) (message.Response, error) {
	return message.NewBindComplete(), nil
}
