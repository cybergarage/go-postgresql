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
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Authenticator represents a frontend message authenticator.
type Authenticator interface {
	// Authenticate handles the Start-up message and returns an Authentication or ErrorResponse message.
	Authenticate(*Conn, *message.Startup) (message.Response, error)
}

// StatusExecutor represents a backend status message executor.
type StatusExecutor interface {
	// ParameterStatus returns the parameter status.
	ParameterStatus(*Conn) (message.Response, error)
	// BackendKeyData returns the backend key data.
	BackendKeyData(*Conn) (message.Response, error)
}

// ParseExecutor represents a backend parse message executor.
type ParseExecutor interface {
	// Parse returns the parse response.
	Parse(*Conn, *message.Parse) (message.Response, error)
	// Bind returns the bind response.
	Bind(*Conn, *message.Bind) (message.Response, error)
}

// ProtocolExecutor represents a backend protocol message executor.
type ProtocolExecutor interface {
	Authenticator
	StatusExecutor
	ParseExecutor
}

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	CreateDatabase(*Conn, *query.CreateDatabase) (message.Response, error)
	CreateTable(*Conn, *query.CreateTable) (message.Response, error)
	CreateIndex(*Conn, *query.CreateIndex) (message.Response, error)
	DropDatabase(*Conn, *query.DropDatabase) (message.Response, error)
	DropTable(*Conn, *query.DropTable) (message.Response, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	Insert(*Conn, *query.Insert) (message.Response, error)
	Select(*Conn, *query.Select) (message.Response, error)
	Update(*Conn, *query.Update) (message.Response, error)
	Delete(*Conn, *query.Delete) (message.Response, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
}

// Executor represents a frontend message executor.
type Executor interface {
	ProtocolExecutor
}
