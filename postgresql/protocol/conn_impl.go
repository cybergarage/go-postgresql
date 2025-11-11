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
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// connOption represents a connection option.
type connOption = func(*conn)

// conn represents a connection of PostgreSQL binary protocol.
type conn struct {
	net.Conn

	isClosed      bool
	msgReader     *MessageReader
	db            string
	schemas       []string
	user          string
	ts            time.Time
	uuid          uuid.UUID
	id            ConnID
	tracerContext tracer.Context
	tlsConn       *tls.Conn
	txMutex       sync.Mutex
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(netConn net.Conn, opts ...connOption) *conn {
	conn := &conn{
		Conn:          netConn,
		isClosed:      false,
		msgReader:     NewMessageReaderWith(WithMessageReadeConn(netConn)),
		db:            "",
		schemas:       []string{},
		user:          "",
		ts:            time.Now(),
		uuid:          uuid.New(),
		id:            0,
		tracerContext: nil,
		tlsConn:       nil,
		txMutex:       sync.Mutex{},
	}
	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// WithConnDatabase sets a database name.
func WithConnDatabase(name string) func(*conn) {
	return func(conn *conn) {
		conn.db = name
	}
}

// WithConnTracer sets a tracer context.
func WithConnTracer(t tracer.Context) func(*conn) {
	return func(conn *conn) {
		conn.tracerContext = t
	}
}

// WithConnTLSConn sets a TLS connection.
func WithConnTLSConn(s *tls.Conn) func(*conn) {
	return func(conn *conn) {
		conn.tlsConn = s
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

// SetSchemas sets the schema names.
func (conn *conn) SetSchemas(schemas ...string) {
	conn.schemas = schemas
}

// Schemas returns the schema names.
func (conn *conn) Schemas() []string {
	return conn.schemas
}

// SetUser sets the user name.
func (conn *conn) SetUser(user string) {
	conn.user = user
}

// User returns the user name.
func (conn *conn) User() string {
	return conn.user
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

// Context returns the context of the connection.
func (conn *conn) Context() context.Context {
	return context.Background()
}

// SetSpanContext sets the tracer span context of the connection.
func (conn *conn) SetSpanContext(ctx tracer.Context) {
	conn.tracerContext = ctx
}

// SpanContext returns the tracer span context of the connection.
func (conn *conn) SpanContext() tracer.Context {
	return conn.tracerContext
}

// Span returns the current top tracer span on the tracer span stack.
func (conn *conn) Span() tracer.Span {
	return conn.tracerContext.Span()
}

// StartSpan starts a new child tracer span and pushes it onto the tracer span stack.
func (conn *conn) StartSpan(name string) bool {
	return conn.tracerContext.StartSpan(name)
}

// FinishSpan ends the current top tracer span and pops it from the tracer span stack.
func (conn *conn) FinishSpan() bool {
	return conn.tracerContext.FinishSpan()
}

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *conn) IsTLSConnection() bool {
	return conn.tlsConn != nil
}

// TLSConnectionState returns the TLS connection state.
func (conn *conn) TLSConn() *tls.Conn {
	return conn.tlsConn
}

// MessageReader returns a message reader.
func (conn *conn) MessageReader() *MessageReader {
	return conn.msgReader
}

// LockTransaction locks the transaction.
func (conn *conn) LockTransaction() error {
	if !conn.txMutex.TryLock() {
		return ErrTransactionBlocked
	}
	return nil
}

// UnlockTransaction unlocks the transaction.
func (conn *conn) UnlockTransaction() error {
	conn.txMutex.Unlock()
	return nil
}

// TransactionStatus returns the transaction status.
func (conn *conn) TransactionStatus() TransactionStatus {
	// Check if the transaction is locked
	if conn.txMutex.TryLock() {
		conn.txMutex.Unlock()
		return TransactionIdle
	}
	return TransactionBlock
}

// ResponseMessage sends a response.
func (conn *conn) ResponseMessage(resMsg Response) error {
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
func (conn *conn) ResponseMessages(resMsgs Responses) error {
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
	errMsg, err := NewErrorResponseWith(err)
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

// SkipMessage skips a.
func (conn *conn) SkipMessage() error {
	msg, err := NewMessageWithReader(conn.MessageReader())
	if err != nil {
		return err
	}
	_, err = msg.ReadMessageData()
	if err != nil {
		return err
	}
	return nil
}

// ReadyForMessage sends a ready for.
func (conn *conn) ReadyForMessage() error {
	readyMsg, err := NewReadyForQueryWith(conn.TransactionStatus())
	if err != nil {
		return err
	}
	err = conn.ResponseMessage(readyMsg)
	if err != nil {
		return err
	}
	return nil
}
