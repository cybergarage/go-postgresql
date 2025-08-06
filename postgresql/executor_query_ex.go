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
	sql "github.com/cybergarage/go-sqlparser/sql/query"
)

// defaultExQueryExecutor represents a DMOExtraExecutor instance.
type defaultExQueryExecutor struct {
	QueryExecutor
}

// NewDefaultExQueryExecutorWith returns a defaultDMOExtraExecutor instance with the given DMOExecutor.
func NewDefaultExQueryExecutorWith(executor QueryExecutor) *defaultExQueryExecutor {
	return &defaultExQueryExecutor{
		QueryExecutor: executor,
	}
}

// SetSQLExecutor sets a SQL executor.
func (executor *defaultExQueryExecutor) SetSQLExecutor(se SQLExecutor) {
	if setter, ok := executor.QueryExecutor.(SQLExecutorSetter); ok {
		setter.SetSQLExecutor(se)
	}
}

// CreateIndex handles a CREATE INDEX query.
func (executor *defaultExQueryExecutor) CreateIndex(conn Conn, stmt query.CreateIndex) (protocol.Responses, error) {
	alterStmt, err := sql.NewAlterTableFrom(stmt)
	if err != nil {
		return nil, err
	}
	return executor.QueryExecutor.AlterTable(conn, alterStmt)
}

// DropIndex handles a DROP INDEX query.
func (executor *defaultExQueryExecutor) DropIndex(conn Conn, stmt query.DropIndex) (protocol.Responses, error) {
	alterStmt, err := sql.NewAlterTableFrom(stmt)
	if err != nil {
		return nil, err
	}
	return executor.QueryExecutor.AlterTable(conn, alterStmt)
}

// Vacuum handles a VACUUM query.
func (executor *defaultExQueryExecutor) Vacuum(conn Conn, stmt query.Vacuum) (protocol.Responses, error) {
	return protocol.NewCommandCompleteResponsesWith(("VACUUM"))
}

// Truncate handles a TRUNCATE query.
func (executor *defaultExQueryExecutor) Truncate(conn Conn, stmt query.Truncate) (protocol.Responses, error) {
	for _, table := range stmt.Tables() {
		stmt := sql.NewDeleteWith(table, sql.NewCondition())
		_, err := executor.QueryExecutor.Delete(conn, stmt)
		if err != nil {
			return nil, err
		}
	}
	return protocol.NewCommandCompleteResponsesWith(("TRUNCATE"))
}
