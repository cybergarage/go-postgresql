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

package server

import (
	"github.com/cybergarage/go-postgresql/examples/go-postgresqld/server/store"
	"github.com/cybergarage/go-postgresql/postgresql"
)

// Server represents a test server.
type Server struct {
	*postgresql.Server
	Store
}

// NewServerWithStore returns a test server instance with the specified store.
func NewServerWithStore(store Store) *Server {
	server := &Server{
		Server: postgresql.NewServer(),
		Store:  store,
	}
	server.SetAuthenticator(server)
	server.SetQueryExecutor(server)
	server.SetBulkExecutor(server)
	server.SetErrorHandler(server)
	server.SetTransactionExecutor(server)
	server.SetSystemQueryExecutor(server)
	return server
}

// NewServer returns a test server instance.
func NewServer() *Server {
	return NewServerWithStore(store.NewMemStore())
}
