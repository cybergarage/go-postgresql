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
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
	"github.com/cybergarage/go-tracing/tracer"
)

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	TCOExecutor
	DDOExecutor
	DMOExecutor
}

// ExQueryExecutor represents a user extended query message executor.
type ExQueryExecutor interface {
	DDOExExecutor
	DMOExExecutor
}

// SystemQueryExecutor represents a system query message executor.
type SystemQueryExecutor interface {
	SystemDMOExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(Conn, string, error) (protocol.Responses, error)
}

// SQLExecutor represents a SQL executor.
type SQLExecutor = query.SQLExecutor

// SQLExecutorSetter represents a SQL executor setter.
type SQLExecutorSetter interface {
	SetSQLExecutor(SQLExecutor)
}

// Server represents a PostgreSQL protocol server.
type Server interface {
	tracer.Tracer
	Config
	AuthManager

	// SetTracer sets a tracing tracer.
	SetTracer(tracer.Tracer)
	// SetSQLExecutor sets a SQL executor.
	SetSQLExecutor(SQLExecutor)
	// SetQueryExecutor sets a user query executor.
	SetQueryExecutor(QueryExecutor)
	// SetExQueryExecutor sets a user query executor.
	SetExQueryExecutor(ExQueryExecutor)
	// SetSystemQueryExecutor sets a system query executor.
	SetSystemQueryExecutor(SystemQueryExecutor)
	// SetBulkQueryExecutor sets a user bulk executor.
	SetBulkQueryExecutor(BulkQueryExecutor)
	// SetErrorHandler sets a user error handler.
	SetErrorHandler(ErrorHandler)

	// SQLExecutor returns a SQL executor.
	SQLExecutor() SQLExecutor
	// QueryExecutor returns a user query executor.
	QueryExecutor() QueryExecutor
	// ExQueryExecutor returns a user extended query executor.
	ExQueryExecutor() ExQueryExecutor
	// SystemQueryExecutor returns a system query executor.
	SystemQueryExecutor() SystemQueryExecutor
	// BulkQueryExecutor returns a user bulk executor.
	BulkQueryExecutor() BulkQueryExecutor
	// ErrorHandler returns a user error handler.
	ErrorHandler() ErrorHandler

	// Start starts the server.
	Start() error
	// Stop stops the server.
	Stop() error
	// Restart restarts the server.
	Restart() error
}
