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
	"bufio"
	"net"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-tracing/tracer"
)

// ConnOption represents a connection option.
type ConnOption = func(*Conn)

// Conn represents a connection of PostgreSQL binary protocol.
type Conn struct {
	conn net.Conn
	*message.MessageReader
	db string
	ts time.Time
	tracer.Context
	PreparedStatementMap
	PreparedPortalMap
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(netConn net.Conn, opts ...ConnOption) *Conn {
	conn := &Conn{
		conn:                 netConn,
		MessageReader:        message.NewMessageReaderWith(bufio.NewReader(netConn)),
		db:                   "",
		ts:                   time.Now(),
		Context:              nil,
		PreparedStatementMap: NewPreparedStatementMap(),
		PreparedPortalMap:    NewPreparedPortalMap(),
	}
	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// WithConnDatabase sets the database name.
func WithConnDatabase(name string) func(*Conn) {
	return func(conn *Conn) {
		conn.db = name
	}
}

// WithConnDatabase sets the database name.
func WithConnTracer(t tracer.Context) func(*Conn) {
	return func(conn *Conn) {
		conn.Context = t
	}
}

// SetDatabase sets the database name.
func (conn *Conn) SetDatabase(db string) {
	conn.db = db
}

// Database returns the database name.
func (conn *Conn) Database() string {
	return conn.db
}

// Conn returns the raw connection.
func (conn *Conn) Conn() net.Conn {
	return conn.conn
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// SetSpanContext sets the tracer span context of the connection.
func (conn *Conn) SetSpanContext(ctx tracer.Context) {
	conn.Context = ctx
}

// SpanContext returns the tracer span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
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
	if _, err := conn.conn.Write(resBytes); err != nil {
		return err
	}
	return nil
}

// ResponseMessages sends response messages.
func (conn *Conn) ResponseMessages(resMsgs message.Responses) error {
	if resMsgs == nil {
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
func (conn *Conn) ResponseError(err error) {
	errMsg, err := message.NewErrorResponseWith(err)
	if err != nil {
		log.Error(err)
		return
	}
	errBytes, err := errMsg.Bytes()
	if err != nil {
		log.Error(err)
		return
	}
	if _, err := conn.conn.Write(errBytes); err != nil {
		log.Error(err)
	}
}
