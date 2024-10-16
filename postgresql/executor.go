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
	"github.com/cybergarage/go-postgresql/postgresql/auth"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// StartupAuthHandler represents a start-up authenticatation handler.
type StartupAuthHandler interface {
	// Authenticate handles the Start-up message and returns an Authentication or ErrorResponse protocol.
	Authenticate(Conn) (protocol.Response, error)
	// AddAuthenticator adds a new authenticator.
	AddAuthenticator(auth.Authenticator)
	// ClearAuthenticators clears all authenticators.
	ClearAuthenticators()
}

// StartupAuthHandler represents a start-up message handler.
type StartupHandler interface {
	StartupAuthHandler
	// ParameterStatuses returns the parameter statuses.
	ParameterStatuses(Conn) (protocol.Responses, error)
	// BackendKeyData returns the backend key data.
	BackendKeyData(Conn) (protocol.Response, error)
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

// DMOExtraExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExtraExecutor interface {
	// Vacuum handles a VACUUM query.
	Vacuum(Conn, query.Vacuum) (protocol.Responses, error)
	// Truncate handles a TRUNCATE query.
	Truncate(Conn, query.Truncate) (protocol.Responses, error)
}

// TCLExecutor defines a executor interface for TCL (Transaction Control Language).
type TCLExecutor interface {
	// Begin handles a BEGIN query.
	Begin(Conn, query.Begin) (protocol.Responses, error)
	// Commit handles a COMMIT query.
	Commit(Conn, query.Commit) (protocol.Responses, error)
	// Rollback handles a ROLLBACK query.
	Rollback(Conn, query.Rollback) (protocol.Responses, error)
}

// ExtendedQueryExecutor defines a executor interface for extended query operations.
type ExtendedQueryExecutor interface {
	// Query handles a query protocol.
	Query(Conn, *protocol.Query) (protocol.Responses, error)
	// Prepare handles a parse protocol.
	Parse(Conn, *protocol.Parse) (protocol.Responses, error)
	// Bind handles a bind protocol.
	Bind(Conn, *protocol.Bind) (protocol.Responses, error)
	// Describe handles a describe protocol.
	Describe(Conn, *protocol.Describe) (protocol.Responses, error)
	// Execute handles a execute protocol.
	Execute(Conn, *protocol.Execute) (protocol.Responses, error)
	// Close handles a close protocol.
	Close(Conn, *protocol.Close) (protocol.Responses, error)
	// Sync handles a sync protocol.
	Sync(Conn, *protocol.Sync) (protocol.Responses, error)
	// Flush handles a flush protocol.
	Flush(Conn, *protocol.Flush) (protocol.Responses, error)
}

// BulkExecutor defines a executor interface for bulk operations.
type BulkExecutor interface {
	// Copy handles a COPY query.
	Copy(Conn, query.Copy) (protocol.Responses, error)
	// CopyData handles a COPY data protocol.
	CopyData(Conn, query.Copy, *CopyStream) (protocol.Responses, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
}

// SystemQueryExecutor represents a system query message executor.
type SystemQueryExecutor interface {
	// SystemSelect handles a SELECT query for system tables.
	SystemSelect(Conn, query.Select) (protocol.Responses, error)
}

// QueryExtraExecutor represents a user query message executor.
type QueryExtraExecutor interface {
	DMOExtraExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(Conn, string, error) (protocol.Responses, error)
}

// Executor represents a frontend message executor.
type Executor interface { // nolint: interfacebloat
	StartupHandler
	QueryExecutor
	QueryExtraExecutor
	TCLExecutor
	ExtendedQueryExecutor
	SystemQueryExecutor
	BulkExecutor
	ErrorHandler
	// SetStartupHandler sets a user startup handler.
	SetStartupHandler(StartupHandler)
	// SetQueryExecutor sets a user query executor.
	SetQueryExecutor(QueryExecutor)
	// SetQueryExecutor sets a user query executor.
	SetQueryExtraExecutor(QueryExtraExecutor)
	// SetTransactionExecutor sets a user transaction executor.
	SetTransactionExecutor(TCLExecutor)
	// SetSystemQueryExecutor sets a system query executor.
	SetSystemQueryExecutor(SystemQueryExecutor)
	// SetBulkExecutor sets a user bulk executor.
	SetBulkExecutor(BulkExecutor)
	// SetErrorHandler sets a user error handler.
	SetErrorHandler(ErrorHandler)
}
