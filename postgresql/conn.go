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
	"crypto/tls"

	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-sqlparser/sql"
)

// ConnID represents a connection ID.
type ConnID = sql.ConnID

type PreparedConn interface {
	// PreparedStatement returns a prepared statement.
	PreparedStatement(name string) (*PreparedStatement, error)
	// SetPreparedStatement sets a prepared statement.
	SetPreparedStatement(msg *protocol.Parse) error
	// RemovePreparedStatement removes a prepared statement.
	RemovePreparedStatement(name string) error
	// PreparedPortal returns a prepared query statement.
	PreparedPortal(name string) (*PreparedPortal, error)
	// SetPreparedPortal sets a prepared query statement.
	SetPreparedPortal(name string, query *PreparedPortal) error
	// RemovePreparedPortal removes a prepared query statement.
	RemovePreparedPortal(name string) error
}

type MessageConn interface {
	// SetStartupMessage sets a startup protocol.
	SetStartupMessage(msg *protocol.Startup)
	// StartupMessage return the startup protocol.
	StartupMessage() (*protocol.Startup, bool)
	// MessageReader returns a message reader.
	MessageReader() *protocol.MessageReader
	// ResponseMessage sends a response protocol.
	ResponseMessage(resMsg protocol.Response) error
	// ResponseMessages sends response messages.
	ResponseMessages(resMsgs protocol.Responses) error
	// ResponseError sends an error response.
	ResponseError(err error) error
	// SkipMessage skips a protocol.
	SkipMessage() error
	// ReadyForMessage sends a ready for protocol.
	ReadyForMessage(status protocol.TransactionStatus) error
}

type TLSConn interface {
	// IsTLSConnection return true if the connection is enabled TLS.
	IsTLSConnection() bool
	// TLSConnectionState returns the TLS connection state.
	TLSConnectionState() (*tls.ConnectionState, bool)
}

// Conn represents a connection.
type Conn interface {
	sql.Conn
	PreparedConn
	MessageConn
	TLSConn
}
