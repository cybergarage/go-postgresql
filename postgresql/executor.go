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
)

// ProtocolStartupHander represents a start-up message handler.
type ProtocolStartupHander interface {
	AuthManager
	protocol.StartupHandler
}

// ProtocolQueryHandler represents a query message handler.
type ProtocolQueryHandler interface {
	protocol.QueryHandler
}

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase handles a CREATE DATABASE query.
	CreateDatabase(Conn, query.CreateDatabase) (protocol.Responses, error)
	// CreateTable handles a CREATE TABLE query.
	CreateTable(Conn, query.CreateTable) (protocol.Responses, error)
	// AlterDatabase handles a ALTER DATABASE query.
	AlterDatabase(Conn, query.AlterDatabase) (protocol.Responses, error)
	// AlterTable handles a ALTER TABLE query.
	AlterTable(Conn, query.AlterTable) (protocol.Responses, error)
	// DropDatabase handles a DROP DATABASE query.
	DropDatabase(Conn, query.DropDatabase) (protocol.Responses, error)
	// DropIndex handles a DROP INDEX query.
	DropTable(Conn, query.DropTable) (protocol.Responses, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert handles a INSERT query.
	Insert(Conn, query.Insert) (protocol.Responses, error)
	// Select handles a SELECT query.
	Select(Conn, query.Select) (protocol.Responses, error)
	// Update handles a UPDATE query.
	Update(Conn, query.Update) (protocol.Responses, error)
	// Delete handles a DELETE query.
	Delete(Conn, query.Delete) (protocol.Responses, error)
}

// DMOExExecutor defines a executor interface for extended DMO (Data Manipulation Operations).
type DMOExExecutor interface {
	// Vacuum handles a VACUUM query.
	Vacuum(Conn, query.Vacuum) (protocol.Responses, error)
	// Truncate handles a TRUNCATE query.
	Truncate(Conn, query.Truncate) (protocol.Responses, error)
}

// TCOExecutor defines a executor interface for TCL (Transaction Control Operations).
type TCOExecutor interface {
	// Begin handles a BEGIN query.
	Begin(Conn, query.Begin) (protocol.Responses, error)
	// Commit handles a COMMIT query.
	Commit(Conn, query.Commit) (protocol.Responses, error)
	// Rollback handles a ROLLBACK query.
	Rollback(Conn, query.Rollback) (protocol.Responses, error)
}

// BulkQueryExecutor defines a executor interface for bulk operations.
type BulkQueryExecutor interface {
	// Copy handles a COPY query.
	Copy(Conn, query.Copy) (protocol.Responses, error)
	// CopyData handles a COPY data protocol.
	CopyData(Conn, query.Copy, *CopyStream) (protocol.Responses, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	TCOExecutor
	DDOExecutor
	DMOExecutor
}

// SystemDMOExecutor represents a system DMO message executor.
type SystemDMOExecutor interface {
	// Select handles a SELECT query.
	SystemSelect(Conn, query.Select) (protocol.Responses, error)
}

// SystemQueryExecutor represents a system query message executor.
type SystemQueryExecutor interface {
	SystemDMOExecutor
}

// ExQueryExecutor represents a user extended query message executor.
type ExQueryExecutor interface {
	DMOExExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(Conn, string, error) (protocol.Responses, error)
}
