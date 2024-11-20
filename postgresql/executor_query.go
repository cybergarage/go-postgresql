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
	"github.com/cybergarage/go-postgresql/postgresql/errors"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// baseQueryExecutor represents a base query message executor.
type baseQueryExecutor struct {
	sqlExecutor SQLExecutor
}

// NewQueryExecutor returns a base query message executor.
func NewueryExecutor() QueryExecutor {
	return &baseQueryExecutor{
		sqlExecutor: nil,
	}
}

// NewSQLQueryExecutor returns a base frontend message executor.
func NewSQLQueryExecutor() SQLQueryExecutor {
	return &baseQueryExecutor{
		sqlExecutor: nil,
	}
}

// SetSQLExecutor sets a SQL executor.
func (executor *baseQueryExecutor) SetSQLExecutor(se SQLExecutor) {
	executor.sqlExecutor = se
}

// Begin handles a BEGIN query.
func (executor *baseQueryExecutor) Begin(conn Conn, stmt query.Begin) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("BEGIN")
	}

	err := executor.sqlExecutor.Commit(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// Commit handles a COMMIT query.
func (executor *baseQueryExecutor) Commit(conn Conn, stmt query.Commit) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("COMMIT")
	}

	err := executor.sqlExecutor.Commit(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// Rollback handles a ROLLBACK query.
func (executor *baseQueryExecutor) Rollback(conn Conn, stmt query.Rollback) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("ROLLBACK")
	}

	err := executor.sqlExecutor.Rollback(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// CreateDatabase handles a CREATE DATABASE query.
func (executor *baseQueryExecutor) CreateDatabase(conn Conn, stmt query.CreateDatabase) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("CREATE DATABASE")
	}

	err := executor.sqlExecutor.CreateDatabase(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// CreateTable handles a CREATE TABLE query.
func (executor *baseQueryExecutor) CreateTable(conn Conn, stmt query.CreateTable) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("CREATE TABLE")
	}

	err := executor.sqlExecutor.CreateTable(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// AlterDatabase handles a ALTER DATABASE query.
func (executor *baseQueryExecutor) AlterDatabase(conn Conn, stmt query.AlterDatabase) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("ALTER DATABASE")
	}

	err := executor.sqlExecutor.AlterDatabase(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// AlterTable handles a ALTER TABLE query.
func (executor *baseQueryExecutor) AlterTable(conn Conn, stmt query.AlterTable) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("ALTER TABLE")
	}

	err := executor.sqlExecutor.AlterTable(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// DropDatabase handles a DROP DATABASE query.
func (executor *baseQueryExecutor) DropDatabase(conn Conn, stmt query.DropDatabase) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("DROP DATABASE")
	}

	err := executor.sqlExecutor.DropDatabase(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// DropIndex handles a DROP INDEX query.
func (executor *baseQueryExecutor) DropTable(conn Conn, stmt query.DropTable) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("DROP TABLE")
	}

	err := executor.sqlExecutor.DropTable(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewCommandCompleteResponsesWith(stmt.String())
}

// Insert handles a INSERT query.
func (executor *baseQueryExecutor) Insert(conn Conn, stmt query.Insert) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("INSERT")
	}
	err := executor.sqlExecutor.Insert(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewInsertCompleteResponsesWith(1)
}

// Select handles a SELECT query.
func (executor *baseQueryExecutor) Select(conn Conn, stmt query.Select) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("SELECT")
	}

	rs, err := executor.sqlExecutor.Select(conn, stmt)
	if err != nil {
		return nil, err
	}

	return query.NewResponseFromResultSet(stmt, rs)
}

// Update handles a UPDATE query.
func (executor *baseQueryExecutor) Update(conn Conn, stmt query.Update) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("UPDATE")
	}

	rs, err := executor.sqlExecutor.Update(conn, stmt)
	if err != nil {
		return nil, err
	}

	return protocol.NewUpdateCompleteResponsesWith(int(rs.RowsAffected()))
}

// Delete handles a DELETE query.
func (executor *baseQueryExecutor) Delete(conn Conn, stmt query.Delete) (protocol.Responses, error) {
	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("DELETE")
	}

	rs, err := executor.sqlExecutor.Delete(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewDeleteCompleteResponsesWith(int(rs.RowsAffected()))
}
