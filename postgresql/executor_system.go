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
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// BaseSystemQueryExecutor represents a base query message executor.
type BaseSystemQueryExecutor struct {
}

// NewBaseSystemQueryExecutor returns a base frontend message executor.
func NewBaseSystemQueryExecutor() *BaseSystemQueryExecutor {
	return &BaseSystemQueryExecutor{}
}

// SystemSelect handles a SELECT query for system tables.
func (executor *BaseSystemQueryExecutor) SystemSelect(*Conn, *query.Select) (message.Responses, error) {
	// PostgreSQL: Documentation: 8.0: System Catalogs
	// https://www.postgresql.org/docs/8.0/catalogs.html
	// PostgreSQL: Documentation: 16: Part IV. Client Interfaces
	// https://www.postgresql.org/docs/current/client-interfaces.html
	return nil, query.NewErrNotImplemented("SELECT")
}
