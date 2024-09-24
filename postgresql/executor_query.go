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
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// BaseQueryExecutor represents a base query message executor.
type BaseQueryExecutor struct {
}

// NewBaseQueryExecutor returns a base frontend message executor.
func NewBaseQueryExecutor() *BaseQueryExecutor {
	return &BaseQueryExecutor{}
}

// CreateDatabase handles a CREATE DATABASE query.
func (executor *BaseQueryExecutor) CreateDatabase(*Conn, *query.CreateDatabase) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("CREATE DATABASE")
}

// CreateTable handles a CREATE TABLE query.
func (executor *BaseQueryExecutor) CreateTable(*Conn, *query.CreateTable) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("CREATE TABLE")
}

// AlterDatabase handles a ALTER DATABASE query.
func (executor *BaseQueryExecutor) AlterDatabase(*Conn, *query.AlterDatabase) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("ALTER DATABASE")
}

// AlterTable handles a ALTER TABLE query.
func (executor *BaseQueryExecutor) AlterTable(*Conn, *query.AlterTable) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("ALTER TABLE")
}

// DropDatabase handles a DROP DATABASE query.
func (executor *BaseQueryExecutor) DropDatabase(*Conn, *query.DropDatabase) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("DROP DATABASE")
}

// DropIndex handles a DROP INDEX query.
func (executor *BaseQueryExecutor) DropTable(*Conn, *query.DropTable) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("DROP TABLE")
}

// Insert handles a INSERT query.
func (executor *BaseQueryExecutor) Insert(*Conn, *query.Insert) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("INSERT")
}

// Select handles a SELECT query.
func (executor *BaseQueryExecutor) Select(*Conn, *query.Select) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("SELECT")
}

// Update handles a UPDATE query.
func (executor *BaseQueryExecutor) Update(*Conn, *query.Update) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("UPDATE")
}

// Delete handles a DELETE query.
func (executor *BaseQueryExecutor) Delete(*Conn, *query.Delete) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("DELETE")
}
