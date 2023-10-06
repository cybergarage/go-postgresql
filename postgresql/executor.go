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
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// Authenticator represents a frontend message authenticator.
type Authenticator interface {
	// Authenticate handles the Start-up message and returns an Authentication or ErrorResponse message.
	Authenticate(*Conn, *message.Startup) (message.Response, error)
}

// StatusHandler represents a backend status message handler.
type StatusHandler interface {
	// ParameterStatuses returns the parameter statuses.
	ParameterStatuses(*Conn) (message.Responses, error)
	// BackendKeyData returns the backend key data.
	BackendKeyData(*Conn) (message.Response, error)
}

// StartupHandler represents a backend protocol message handler.
type StartupHandler interface {
	StatusHandler
}

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase handles a CREATE DATABASE query.
	CreateDatabase(*Conn, *query.CreateDatabase) (message.Responses, error)
	// CreateTable handles a CREATE TABLE query.
	CreateTable(*Conn, *query.CreateTable) (message.Responses, error)
	// AlterDatabase handles a ALTER DATABASE query.
	AlterDatabase(*Conn, *query.AlterDatabase) (message.Responses, error)
	// AlterTable handles a ALTER TABLE query.
	AlterTable(*Conn, *query.AlterTable) (message.Responses, error)
	// DropDatabase handles a DROP DATABASE query.
	DropDatabase(*Conn, *query.DropDatabase) (message.Responses, error)
	// DropIndex handles a DROP INDEX query.
	DropTable(*Conn, *query.DropTable) (message.Responses, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert handles a INSERT query.
	Insert(*Conn, *query.Insert) (message.Responses, error)
	// Select handles a SELECT query.
	Select(*Conn, *query.Select) (message.Responses, error)
	// Update handles a UPDATE query.
	Update(*Conn, *query.Update) (message.Responses, error)
	// Delete handles a DELETE query.
	Delete(*Conn, *query.Delete) (message.Responses, error)
}

// DMOExtraExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExtraExecutor interface {
	// Vacuum handles a VACUUM query.
	Vacuum(*Conn, *query.Vacuum) (message.Responses, error)
	// Truncate handles a TRUNCATE query.
	Truncate(*Conn, *query.Truncate) (message.Responses, error)
}

// TCLExecutor defines a executor interface for TCL (Transaction Control Language).
type TCLExecutor interface {
	// Begin handles a BEGIN query.
	Begin(*Conn, *query.Begin) (message.Responses, error)
	// Commit handles a COMMIT query.
	Commit(*Conn, *query.Commit) (message.Responses, error)
	// Rollback handles a ROLLBACK query.
	Rollback(*Conn, *query.Rollback) (message.Responses, error)
}

// ExtendedQueryExecutor defines a executor interface for extended query operations.
type ExtendedQueryExecutor interface {
	// Query handles a query message.
	Query(*Conn, *message.Query) (message.Responses, error)
	// Prepare handles a parse message.
	Parse(*Conn, *message.Parse) (message.Responses, error)
	// Bind handles a bind message.
	Bind(*Conn, *message.Bind) (message.Responses, error)
	// Describe handles a describe message.
	Describe(*Conn, *message.Describe) (message.Responses, error)
	// Execute handles a execute message.
	Execute(*Conn, *message.Execute) (message.Responses, error)
	// Close handles a close message.
	Close(*Conn, *message.Close) (message.Responses, error)
	// Sync handles a sync message.
	Sync(*Conn, *message.Sync) (message.Responses, error)
	// Flush handles a flush message.
	Flush(*Conn, *message.Flush) (message.Responses, error)
}

// BulkExecutor defines a executor interface for bulk operations.
type BulkExecutor interface {
	// Copy handles a COPY query.
	Copy(*Conn, *query.Copy) (message.Responses, error)
	// CopyData handles a COPY data message.
	CopyData(*Conn, *query.Copy, *CopyStream) (message.Responses, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
}

// SystemQueryExecutor represents a system query message executor.
type SystemQueryExecutor interface {
	// SystemSelect handles a SELECT query for system tables.
	SystemSelect(*Conn, *query.Select) (message.Responses, error)
}

// QueryExtraExecutor represents a user query message executor.
type QueryExtraExecutor interface {
	DMOExtraExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(*Conn, string, error) (message.Responses, error)
}

// Executor represents a frontend message executor.
type Executor interface { // nolint: interfacebloat
	Authenticator
	StartupHandler
	QueryExecutor
	QueryExtraExecutor
	TCLExecutor
	ExtendedQueryExecutor
	SystemQueryExecutor
	BulkExecutor
	ErrorHandler
	// SetAuthenticator sets a user authenticator.
	SetAuthenticator(Authenticator)
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
