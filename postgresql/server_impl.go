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
	*BaseExecutor
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		Server:       protocol.NewServer(),
		BaseExecutor: NewBaseExecutor(),
	}
	server.Server.SetProductName(PackageName)
	server.Server.SetProductVersion(Version)
	server.Server.SetMessageHandler(server)
	server.SetQueryExecutor(server)
	server.SetBulkQueryExecutor(server)
	server.SetErrorHandler(server)
	server.SetSystemQueryExecutor(server)

	return server
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
