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

// BaseQueryExecutor represents a base query message executor.
type BaseQueryExecutor struct {
}

// NewBaseQueryExecutor returns a base frontend message executor.
func NewBaseQueryExecutor() *BaseQueryExecutor {
	return &BaseQueryExecutor{}
}

// CreateDatabase handles a CREATE DATABASE query.
func (executor *BaseQueryExecutor) CreateDatabase(*Conn, *query.CreateDatabase) ([]message.Response, error) {
	return nil, nil
}

// CreateTable handles a CREATE TABLE query.
func (executor *BaseQueryExecutor) CreateTable(*Conn, *query.CreateTable) ([]message.Response, error) {
	return nil, nil
}

// CreateIndex handles a CREATE INDEX query.
func (executor *BaseQueryExecutor) CreateIndex(*Conn, *query.CreateIndex) ([]message.Response, error) {
	return nil, nil
}

// DropDatabase handles a DROP DATABASE query.
func (executor *BaseQueryExecutor) DropDatabase(*Conn, *query.DropDatabase) ([]message.Response, error) {
	return nil, nil
}

// DropIndex handles a DROP INDEX query.
func (executor *BaseQueryExecutor) DropTable(*Conn, *query.DropTable) ([]message.Response, error) {
	return nil, nil
}

// Insert handles a INSERT query.
func (executor *BaseQueryExecutor) Insert(*Conn, *query.Insert) ([]message.Response, error) {
	return nil, nil
}

// Select handles a SELECT query.
func (executor *BaseQueryExecutor) Select(*Conn, *query.Select) ([]message.Response, error) {
	return nil, nil
}

// Update handles a UPDATE query.
func (executor *BaseQueryExecutor) Update(*Conn, *query.Update) ([]message.Response, error) {
	return nil, nil
}

// Delete handles a DELETE query.
func (executor *BaseQueryExecutor) Delete(*Conn, *query.Delete) ([]message.Response, error) {
	return nil, nil
}
