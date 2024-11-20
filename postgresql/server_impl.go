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
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		Server:                 protocol.NewServer(),
		protocolStartupHandler: newProtocolStartupHandler(),
		protocolQueryHandler:   newProtocolQueryHandler(),
		sqlExecutor:            nil,
		queryExecutor:          NewSQLQueryExecutor(),
		exQueryExecutor:        nil,
		bulkQueryExecutor:      NewBaseBulkExecutor(),
		errorHandler:           NewBaseErrorHandler(),
		systemQueryExecutor:    NewSystemQueryExecutor(),
	}
	server.exQueryExecutor = newExtraQueryExecutorWith(server.queryExecutor)

	server.Server.SetProductName(PackageName)
	server.Server.SetProductVersion(Version)
	server.Server.SetMessageHandler(server)

	server.SetQueryExecutor(server.queryExecutor)
	server.SetSystemQueryExecutor(server.systemQueryExecutor)
	server.SetBulkQueryExecutor(server.bulkQueryExecutor)
	server.SetErrorHandler(server.errorHandler)

	return server
}

// SetSQLExecutor sets a SQL server.
func (server *server) SetSQLExecutor(se SQLExecutor) {
	server.sqlExecutor = se
	if server.queryExecutor != nil {
		if executor, ok := server.queryExecutor.(SQLQueryExecutor); ok {
			executor.SetSQLExecutor(se)
		}
	}
	if server.systemQueryExecutor != nil {
		if executor, ok := server.systemQueryExecutor.(SQLSystemQueryExecutor); ok {
			executor.SetSQLExecutor(se)
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
