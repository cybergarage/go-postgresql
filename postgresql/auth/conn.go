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
	"net"

	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// Conn represents a connection.
type Conn interface {
	net.Conn

	// StartupMessage return the startup message.
	StartupMessage() (*message.Startup, bool)
	// IsTLSConnection return true if the connection is enabled TLS.
	IsTLSConnection() bool
	// TLSConnectionState returns the TLS connection state.
	TLSConnectionState() (*tls.ConnectionState, bool)

	// MessageReader returns a message reader.
	MessageReader() *message.MessageReader
	// ResponseMessage returns a response message.
	ResponseMessage(resMsg message.Response) error
	// ResponseMessages returns response messages.
	ResponseMessages(resMsgs message.Responses) error
	// ResponseError returns a response error.
	ResponseError(err error) error
}
