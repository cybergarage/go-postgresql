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
	"github.com/cybergarage/go-sqlparser/sql/query"
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

// BindHandler represents a backend parse message handler.
type BindHandler interface {
	// Parse returns the parse response.
	Parse(*Conn, *message.Parse) (message.Response, error)
	// Bind returns the bind response.
	Bind(*Conn, *message.Parse, *message.Bind) (message.Response, error)
}

// ProtocolHandler represents a backend protocol message handler.
type ProtocolHandler interface {
	StatusHandler
	BindHandler
}

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase handles a CREATE DATABASE query.
	CreateDatabase(*Conn, *query.CreateDatabase) ([]message.Response, error)
	// CreateTable handles a CREATE TABLE query.
	CreateTable(*Conn, *query.CreateTable) ([]message.Response, error)
	// CreateIndex handles a CREATE INDEX query.
	CreateIndex(*Conn, *query.CreateIndex) ([]message.Response, error)
	// DropDatabase handles a DROP DATABASE query.
	DropDatabase(*Conn, *query.DropDatabase) ([]message.Response, error)
	// DropIndex handles a DROP INDEX query.
	DropTable(*Conn, *query.DropTable) ([]message.Response, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert handles a INSERT query.
	Insert(*Conn, *query.Insert) ([]message.Response, error)
	// Select handles a SELECT query.
	Select(*Conn, *query.Select) ([]message.Response, error)
	// Update handles a UPDATE query.
	Update(*Conn, *query.Update) ([]message.Response, error)
	// Delete handles a DELETE query.
	Delete(*Conn, *query.Delete) ([]message.Response, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
}

// Executor represents a frontend message executor.
type Executor interface {
	Authenticator
	ProtocolHandler
	QueryExecutor
}
