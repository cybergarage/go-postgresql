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
	"github.com/cybergarage/go-postgresql/postgresql/query"
	"github.com/cybergarage/go-tracing/tracer"
)

// SQLExecutor represents a SQL executor.
type SQLExecutor = query.SQLExecutor

// ProtocolExecutor represents a protocol message executor.
type ProtocolExecutor interface {
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
}

// Server represents a PostgreSQL protocol server.
type Server interface {
	tracer.Tracer
	Config
	AuthManager
	ProtocolExecutor
	// SetSQLExecutor sets a SQL executor.
	SetSQLExecutor(SQLExecutor)
	// SetMessageHandler sets a message handler.
	SetTracer(tracer.Tracer)
	// Start starts the server.
	Start() error
	// Stop stops the server.
	Stop() error
	// Restart restarts the server.
	Restart() error
}
