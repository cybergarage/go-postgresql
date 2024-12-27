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

package protocol

import (
	"crypto/tls"

	"github.com/cybergarage/go-postgresql/postgresql/net"
)

// ConnID represents a connection ID.
type ConnID = net.ConnID

// MessageConn represents a message connection.
type MessageConn interface {
	// MessageReader returns a message reader.
	MessageReader() *MessageReader
	// ResponseMessage sends a response
	ResponseMessage(resMsg Response) error
	// ResponseMessages sends response messages.
	ResponseMessages(resMsgs Responses) error
	// ResponseError sends an error response.
	ResponseError(err error) error
	// SkipMessage skips a
	SkipMessage() error
	// ReadyForMessage sends a ready for
	ReadyForMessage(status TransactionStatus) error
}

// TLSConn represents a TLS connection.
type TLSConn interface {
	// IsTLSConnection return true if the connection is enabled TLS.
	IsTLSConnection() bool
	// TLSConn returns a TLS connection.
	TLSConn() *tls.Conn
}

// Conn represents a connection.
type Conn interface {
	net.Conn
	MessageConn
	TLSConn
}
