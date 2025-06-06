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

// defaultSystemQueryExecutor represents a base system query executor.
type defaultSystemQueryExecutor struct {
	sqlExecutor SQLExecutor
}

// NewNullSystemQueryExecutor returns a new null SystemQueryExecutor.
func NewNullSystemQueryExecutor() SystemQueryExecutor {
	return &defaultSystemQueryExecutor{
		sqlExecutor: nil,
	}
}

// NewDefaultSystemQueryExecutor returns a default SystemQueryExecutor based on SQLExecutor.
func NewDefaultSystemQueryExecutor() SystemQueryExecutor {
	return &defaultSystemQueryExecutor{
		sqlExecutor: nil,
	}
}

// SetSQLExecutor sets a SQL executor.
func (executor *defaultSystemQueryExecutor) SetSQLExecutor(se SQLExecutor) {
	executor.sqlExecutor = se
}

// SystemSelect handles a SELECT query for system tables.
func (executor *defaultSystemQueryExecutor) SystemSelect(conn Conn, stmt query.Select) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 8.0: System Catalogs
	// https://www.postgresql.org/docs/8.0/catalogs.html
	// PostgreSQL: Documentation: 16: Part IV. Client Interfaces
	// https://www.postgresql.org/docs/current/client-interfaces.html

	if executor.sqlExecutor == nil {
		return nil, errors.NewErrNotImplemented("SELECT")
	}

	rs, err := executor.sqlExecutor.SystemSelect(conn, stmt)
	if err != nil {
		return nil, err
	}

	return query.NewResponseFromResultSet(rs)
}
