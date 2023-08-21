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
	// ParameterStatus returns the parameter status.
	ParameterStatus(*Conn) (message.Response, error)
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
	// CreateIndex handles a CREATE INDEX query.
	CreateIndex(*Conn, *query.CreateIndex) (message.Responses, error)
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

// BulkExecutor defines a executor interface for bulk operations.
type BulkExecutor interface {
	// Copy handles a COPY query.
	Copy(*Conn, *query.Copy, *CopyStream) (message.Responses, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
}

// Executor represents a frontend message executor.
type Executor interface {
	Authenticator
	StartupHandler
	QueryExecutor
	// SetAuthenticator sets a user authenticator.
	SetAuthenticator(Authenticator)
	// SetStartupHandler sets a user startup handler.
	SetStartupHandler(StartupHandler)
	// SetQueryExecutor sets a user query executor.
	SetQueryExecutor(QueryExecutor)
}
