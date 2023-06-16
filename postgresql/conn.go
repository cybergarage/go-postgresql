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
	"net"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
)

// Conn represents a connection of PostgreSQL binary protocol.
type Conn struct {
	Database string
	ts       time.Time
	tracer.Context
	conn net.Conn
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(c net.Conn, t tracer.Context) *Conn {
	conn := &Conn{
		Database: "",
		ts:       time.Now(),
		Context:  t,
		conn:     c,
	}
	return conn
}

// Conn returns the raw connection.
func (conn *Conn) Conn() net.Conn {
	return conn.conn
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// SpanContext returns the tracer span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
}
