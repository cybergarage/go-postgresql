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

package auth

import (
	"crypto/tls"

	"github.com/cybergarage/go-postgresql/postgresql/protocol"
)

// Conn represents a connection.
type Conn interface {
	// StartupMessage return the startup protocol.
	StartupMessage() (*protocol.Startup, bool)
	// IsTLSConnection return true if the connection is enabled TLS.
	IsTLSConnection() bool
	// TLSConn returns a TLS connection.
	TLSConn() *tls.Conn
	// MessageReader returns a message reader.
	MessageReader() *protocol.MessageReader
	// ResponseMessage returns a response protocol.
	ResponseMessage(resMsg protocol.Response) error
	// ResponseMessages returns response messages.
	ResponseMessages(resMsgs protocol.Responses) error
	// ResponseError returns a response error.
	ResponseError(err error) error
}
