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
	"net"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
)

// ConnOption represents a connection option.
type ConnOption = func(*Conn)

// Conn represents a connection of PostgreSQL binary protocol.
type Conn struct {
	db string
	ts time.Time
	tracer.Context
	conn net.Conn
	PreparedQueryMap
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(c net.Conn, opts ...ConnOption) *Conn {
	conn := &Conn{
		db:               "",
		ts:               time.Now(),
		Context:          nil,
		conn:             c,
		PreparedQueryMap: NewPreparedQueryMap(),
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
