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
	"errors"

	"github.com/cybergarage/go-authenticator/auth"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
)

// server represents a PostgreSQL protocol server.
type server struct {
	protocol.Server
	*protocolStartupHandler
	*protocolQueryHandler

	sqlExecutor         SQLExecutor
	queryExecutor       QueryExecutor
	systemQueryExecutor SystemQueryExecutor
	exQueryExecutor     ExQueryExecutor
	bulkQueryExecutor   BulkQueryExecutor
	errorHandler        ErrorHandler
	authManager         auth.Manager
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		Server:                 protocol.NewServer(),
		protocolStartupHandler: newProtocolStartupHandler(),
		protocolQueryHandler:   newProtocolQueryHandler(),
		sqlExecutor:            nil,
		queryExecutor:          NewDefaultQueryExecutor(),
		exQueryExecutor:        nil,
		bulkQueryExecutor:      NewNullBulkExecutor(),
		errorHandler:           NewNullErrorHandler(),
		systemQueryExecutor:    NewNullSystemQueryExecutor(),
		authManager:            auth.NewManager(),
	}

	server.exQueryExecutor = NewDefaultExQueryExecutorWith(
		server.queryExecutor,
	)

	server.SetProductName(PackageName)
	server.SetProductVersion(Version)
	server.SetMessageHandler(server)

	server.SetQueryExecutor(server.queryExecutor)
	server.SetSystemQueryExecutor(server.systemQueryExecutor)
	server.SetBulkQueryExecutor(server.bulkQueryExecutor)
	server.SetErrorHandler(server.errorHandler)

	return server
}

// SetSQLExecutor sets a SQL server.
func (server *server) SetSQLExecutor(sqlExeutor SQLExecutor) {
	server.sqlExecutor = sqlExeutor

	executors := []any{
		server.queryExecutor,
		server.exQueryExecutor,
		server.systemQueryExecutor,
		server.bulkQueryExecutor,
		server.errorHandler,
	}
	for _, executor := range executors {
		if executor == nil {
			continue
		}

		if _, ok := executor.(Server); ok {
			continue
		}

		if setter, ok := executor.(SQLExecutorSetter); ok {
			setter.SetSQLExecutor(sqlExeutor)
		}
	}
}

// SetQueryExecutor sets a user query server.
func (server *server) SetQueryExecutor(qe QueryExecutor) {
	server.queryExecutor = qe
}

// SetExQueryExecutor sets a user query extra server.
func (server *server) SetExQueryExecutor(qe ExQueryExecutor) {
	server.exQueryExecutor = qe
}

// SetBulkQueryExecutor sets a user bulk server.
func (server *server) SetBulkQueryExecutor(be BulkQueryExecutor) {
	server.bulkQueryExecutor = be
}

// SetErrorHandler sets a user error handler.
func (server *server) SetErrorHandler(eh ErrorHandler) {
	server.errorHandler = eh
}

// SetSystemQueryExecutor sets a system query server.
func (server *server) SetSystemQueryExecutor(sq SystemQueryExecutor) {
	server.systemQueryExecutor = sq
}

// SQLExecutor returns a SQL executor.
func (server *server) SQLExecutor() SQLExecutor {
	return server.sqlExecutor
}

// QueryExecutor returns a user query executor.
func (server *server) QueryExecutor() QueryExecutor {
	return server.queryExecutor
}

// ExQueryExecutor returns a user extended query executor.
func (server *server) ExQueryExecutor() ExQueryExecutor {
	return server.exQueryExecutor
}

// SystemQueryExecutor returns a system query executor.
func (server *server) SystemQueryExecutor() SystemQueryExecutor {
	return server.systemQueryExecutor
}

// BulkQueryExecutor returns a user bulk executor.
func (server *server) BulkQueryExecutor() BulkQueryExecutor {
	return server.bulkQueryExecutor
}

// ErrorHandler returns a user error handler.
func (server *server) ErrorHandler() ErrorHandler {
	return server.errorHandler
}

// Start starts the server.
func (server *server) Start() error {
	type starter interface {
		Start() error
	}

	starters := []starter{
		server.Server,
	}
	for _, s := range starters {
		if err := s.Start(); err != nil {
			return server.Stop()
		}
	}

	return nil
}

// Stop stops the server.
func (server *server) Stop() error {
	type stopper interface {
		Stop() error
	}

	stoppers := []stopper{
		server.Server,
	}

	var err error

	for _, s := range stoppers {
		if e := s.Stop(); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// Restart restarts the server.
func (server *server) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	return server.Start()
}
