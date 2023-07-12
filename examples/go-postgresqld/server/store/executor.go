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

package store

import (
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// CreateDatabase handles a CREATE DATABASE query.
func (store *MemStore) CreateDatabase(*postgresql.Conn, *query.CreateDatabase) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("CREATE DATABASE")
}

// CreateTable handles a CREATE TABLE query.
func (store *MemStore) CreateTable(*postgresql.Conn, *query.CreateTable) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("CREATE TABLE")
}

// CreateIndex handles a CREATE INDEX query.
func (store *MemStore) CreateIndex(*postgresql.Conn, *query.CreateIndex) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("CREATE INDEX")
}

// DropDatabase handles a DROP DATABASE query.
func (store *MemStore) DropDatabase(*postgresql.Conn, *query.DropDatabase) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("DROP DATABASE")
}

// DropIndex handles a DROP INDEX query.
func (store *MemStore) DropTable(*postgresql.Conn, *query.DropTable) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("DROP TABLE")
}

// Insert handles a INSERT query.
func (store *MemStore) Insert(*postgresql.Conn, *query.Insert) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("INSERT")
}

// Select handles a SELECT query.
func (store *MemStore) Select(*postgresql.Conn, *query.Select) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("SELECT")
}

// Update handles a UPDATE query.
func (store *MemStore) Update(*postgresql.Conn, *query.Update) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("UPDATE")
}

// Delete handles a DELETE query.
func (store *MemStore) Delete(*postgresql.Conn, *query.Delete) ([]message.Response, error) {
	return nil, postgresql.NewErrNotImplemented("DELETE")
}
