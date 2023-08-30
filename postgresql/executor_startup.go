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

// ParameterStatuses returns the parameter statuses.
func (executor *BaseProtocolExecutor) ParameterStatuses(*Conn) (message.Responses, error) {
	m := map[string]string{}
	m[message.ClientEncoding] = message.EncodingUTF8
	m[message.ServerEncoding] = message.EncodingUTF8
	return message.NewParameterStatusesWith(m)
}

// BackendKeyData returns the backend key data.
func (executor *BaseProtocolExecutor) BackendKeyData(*Conn) (message.Response, error) {
	return message.NewBackendKeyDataWith(executor.processID, executor.secretKey)
}
