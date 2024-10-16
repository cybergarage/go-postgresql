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
	"crypto/tls"
	"net"
	"time"

	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// connOption represents a connection option.
type connOption = func(*conn)

// conn represents a connection of PostgreSQL binary protocol.
type conn struct {
	net.Conn
	isClosed  bool
	msgReader *protocol.MessageReader
	db        string
	ts        time.Time
	uuid      uuid.UUID
	id        ConnID
	tracer.Context
	PreparedStatementMap
	PreparedPortalMap
	tlsState   *tls.ConnectionState
	startupMsg *protocol.Startup
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(netconn net.Conn, opts ...connOption) *conn {
	conn := &conn{
		Conn:                 netconn,
		isClosed:             false,
		msgReader:            protocol.NewMessageReaderWith(netconn),
		db:                   "",
		ts:                   time.Now(),
		uuid:                 uuid.New(),
		id:                   0,
		Context:              nil,
		PreparedStatementMap: NewPreparedStatementMap(),
		PreparedPortalMap:    NewPreparedPortalMap(),
		tlsState:             nil,
		startupMsg:           nil,
	}
	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// WithconnDatabase sets a database name.
func WithconnDatabase(name string) func(*conn) {
	return func(conn *conn) {
		conn.db = name
	}
}

// WithconnTracer sets a tracer context.
func WithconnTracer(t tracer.Context) func(*conn) {
	return func(conn *conn) {
		conn.Context = t
	}
}

// WithconnStartupMessage sets a startup protocol.
func WithconnStartupMessage(msg *protocol.Startup) func(*conn) {
	return func(conn *conn) {
		conn.startupMsg = msg
	}
}

// WithTLSConnectionState sets a TLS connection state.
func WithTLSConnectionState(s *tls.ConnectionState) func(*conn) {
	return func(conn *conn) {
		conn.tlsState = s
	}
}

// Close closes the connection.
func (conn *conn) Close() error {
	if conn.isClosed {
		return nil
	}
	if err := conn.Conn.Close(); err != nil {
		return err
	}
	conn.isClosed = true
	return nil
}

// SetDatabase sets the database name.
func (conn *conn) SetDatabase(db string) {
	conn.db = db
}

// Database returns the database name.
func (conn *conn) Database() string {
	return conn.db
}

// Timestamp returns the creation time of the connection.
func (conn *conn) Timestamp() time.Time {
	return conn.ts
}

// UUID returns the UUID of the connection.
func (conn *conn) UUID() uuid.UUID {
	return conn.uuid
}

// ID returns the ID of the connection.
func (conn *conn) ID() ConnID {
	return conn.id
}

// SetSpanContext sets the tracer span context of the connection.
func (conn *conn) SetSpanContext(ctx tracer.Context) {
	conn.Context = ctx
}

// SpanContext returns the tracer span context of the connection.
func (conn *conn) SpanContext() tracer.Context {
	return conn.Context
}

// SetStartupMessage sets a startup protocol.
func (conn *conn) SetStartupMessage(msg *protocol.Startup) {
	conn.startupMsg = msg
}

// StartupMessage return the startup protocol.
func (conn *conn) StartupMessage() (*protocol.Startup, bool) {
	return conn.startupMsg, conn.startupMsg != nil
}

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *conn) IsTLSConnection() bool {
	return conn.tlsState != nil
}

// TLSConnectionState returns the TLS connection state.
func (conn *conn) TLSConnectionState() (*tls.ConnectionState, bool) {
	return conn.tlsState, conn.tlsState != nil
}

// MessageReader returns a message reader.
func (conn *conn) MessageReader() *protocol.MessageReader {
	return conn.msgReader
}

// ResponseMessage sends a response protocol.
func (conn *conn) ResponseMessage(resMsg protocol.Response) error {
	if resMsg == nil {
		return nil
	}
	resBytes, err := resMsg.Bytes()
	if err != nil {
		return err
	}
	if _, err := conn.Conn.Write(resBytes); err != nil {
		return err
	}
	return nil
}

// ResponseMessages sends response messages.
func (conn *conn) ResponseMessages(resMsgs protocol.Responses) error {
	if len(resMsgs) == 0 {
		return nil
	}
	for _, resMsg := range resMsgs {
		err := conn.ResponseMessage(resMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResponseError sends an error response.
func (conn *conn) ResponseError(err error) error {
	if err == nil {
		return nil
	}
	errMsg, err := protocol.NewErrorResponseWith(err)
	if err != nil {
		return err
	}
	errBytes, err := errMsg.Bytes()
	if err != nil {
		return err
	}
	_, err = conn.Conn.Write(errBytes)
	return err
}

// SkipMessage skips a protocol.
func (conn *conn) SkipMessage() error {
	msg, err := protocol.NewMessageWithReader(conn.MessageReader())
	if err != nil {
		return err
	}
	_, err = msg.ReadMessageData()
	if err != nil {
		return err
	}
	return nil
}

// ReadyForMessage sends a ready for protocol.
func (conn *conn) ReadyForMessage(status protocol.TransactionStatus) error {
	readyMsg, err := protocol.NewReadyForQueryWith(status)
	if err != nil {
		return err
	}
	err = conn.ResponseMessage(readyMsg)
	if err != nil {
		return err
	}
	return nil
}
