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
	"net"
	"time"

	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// ConnOption represents a connection option.
type ConnOption = func(*Conn)

// Conn represents a connection of PostgreSQL binary protocol.
type Conn struct {
	net.Conn
	isClosed  bool
	msgReader *message.MessageReader
	db        string
	ts        time.Time
	uuid      uuid.UUID
	tracer.Context
	PreparedStatementMap
	PreparedPortalMap
	tlsState   *tls.ConnectionState
	startupMsg *message.Startup
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(netConn net.Conn, opts ...ConnOption) *Conn {
	conn := &Conn{
		Conn:                 netConn,
		isClosed:             false,
		msgReader:            message.NewMessageReaderWith(netConn),
		db:                   "",
		ts:                   time.Now(),
		uuid:                 uuid.New(),
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

// WithConnDatabase sets a database name.
func WithConnDatabase(name string) func(*Conn) {
	return func(conn *Conn) {
		conn.db = name
	}
}

// WithConnTracer sets a tracer context.
func WithConnTracer(t tracer.Context) func(*Conn) {
	return func(conn *Conn) {
		conn.Context = t
	}
}

// WithConnStartupMessage sets a startup message.
func WithConnStartupMessage(msg *message.Startup) func(*Conn) {
	return func(conn *Conn) {
		conn.startupMsg = msg
	}
}

// WithTLSConnectionState sets a TLS connection state.
func WithTLSConnectionState(s *tls.ConnectionState) func(*Conn) {
	return func(conn *Conn) {
		conn.tlsState = s
	}
}

// Close closes the connection.
func (conn *Conn) Close() error {
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
func (conn *Conn) SetDatabase(db string) {
	conn.db = db
}

// Database returns the database name.
func (conn *Conn) Database() string {
	return conn.db
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// UUID returns the UUID of the connection.
func (conn *Conn) UUID() uuid.UUID {
	return conn.uuid
}

// SetSpanContext sets the tracer span context of the connection.
func (conn *Conn) SetSpanContext(ctx tracer.Context) {
	conn.Context = ctx
}

// SpanContext returns the tracer span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
}

// SetStartupMessage sets a startup message.
func (conn *Conn) SetStartupMessage(msg *message.Startup) {
	conn.startupMsg = msg
}

// StartupMessage return the startup message.
func (conn *Conn) StartupMessage() (*message.Startup, bool) {
	return conn.startupMsg, conn.startupMsg != nil
}

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *Conn) IsTLSConnection() bool {
	return conn.tlsState != nil
}

// TLSConnectionState returns the TLS connection state.
func (conn *Conn) TLSConnectionState() (*tls.ConnectionState, bool) {
	return conn.tlsState, conn.tlsState != nil
}

// MessageReader returns a message reader.
func (conn *Conn) MessageReader() *message.MessageReader {
	return conn.msgReader
}

// ResponseMessage sends a response message.
func (conn *Conn) ResponseMessage(resMsg message.Response) error {
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
func (conn *Conn) ResponseMessages(resMsgs message.Responses) error {
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
func (conn *Conn) ResponseError(err error) error {
	if err == nil {
		return nil
	}
	errMsg, err := message.NewErrorResponseWith(err)
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

// SkipMessage skips a message.
func (conn *Conn) SkipMessage() error {
	msg, err := message.NewMessageWithReader(conn.MessageReader())
	if err != nil {
		return err
	}
	_, err = msg.ReadMessageData()
	if err != nil {
		return err
	}
	return nil
}

// ReadyForMessage sends a ready for message.
func (conn *Conn) ReadyForMessage(status message.TransactionStatus) error {
	readyMsg, err := message.NewReadyForQueryWith(status)
	if err != nil {
		return err
	}
	err = conn.ResponseMessage(readyMsg)
	if err != nil {
		return err
	}
	return nil
}
