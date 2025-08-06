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
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
)

// protocolStartupHandler  a new protocol startup executor.
type protocolStartupHandler struct {
	processID int32
	secretKey int32
}

// newProtocolStartupHandler returns a new protocol startup executor.
func newProtocolStartupHandler() *protocolStartupHandler {
	r, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		log.Error(err)
	}

	return &protocolStartupHandler{
		processID: int32(os.Getpid()),
		secretKey: int32(r.Int64()),
	}
}

// ParameterStatuses returns the parameter statuses.
func (server *server) ParameterStatuses(Conn) (protocol.Responses, error) {
	m := map[string]string{}
	m[protocol.ClientEncoding] = protocol.EncodingUTF8
	m[protocol.ServerEncoding] = protocol.EncodingUTF8

	return protocol.NewParameterStatusesWith(m)
}

// BackendKeyData returns the backend key data.
func (server *server) BackendKeyData(Conn) (protocol.Response, error) {
	return protocol.NewBackendKeyDataWith(server.processID, server.secretKey)
}
