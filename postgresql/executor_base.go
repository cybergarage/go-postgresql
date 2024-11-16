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

// BaseExecutor represents a base frontend message executor.
type BaseExecutor struct {
	ProtocolStartupHander
	ProtocolQueryHandler
	QueryExecutor
	ExQueryExecutor
	BulkQueryExecutor
	ErrorHandler
	SystemQueryExecutor
	sqlExecutor SQLExecutor
}

// NewBaseExecutor returns a base frontend message executor.
func NewBaseExecutor() *BaseExecutor {
	executor := &BaseExecutor{
		ProtocolStartupHander: newProtocolStartupHandler(),
		ProtocolQueryHandler:  nil,
		QueryExecutor:         NewBaseQueryExecutor(),
		ExQueryExecutor:       nil,
		BulkQueryExecutor:     NewBaseBulkExecutor(),
		ErrorHandler:          NewBaseErrorHandler(),
		SystemQueryExecutor:   NewBaseSystemQueryExecutor(),
		sqlExecutor:           nil,
	}
	executor.ExQueryExecutor = newDMOExtraExecutorWith(executor)
	executor.ProtocolQueryHandler = newProtocolQueryHandlerWith(executor)
	return executor
}

// SetSQLExecutor sets a SQL executor.
func (executor *BaseExecutor) SetSQLExecutor(se SQLExecutor) {
	executor.sqlExecutor = se
}

// SetQueryExecutor sets a user query executor.
func (executor *BaseExecutor) SetQueryExecutor(qe QueryExecutor) {
	executor.QueryExecutor = qe
}

// SetExQueryExecutor sets a user query extra executor.
func (executor *BaseExecutor) SetExQueryExecutor(qe ExQueryExecutor) {
	executor.ExQueryExecutor = qe
}

// SetBulkQueryExecutor sets a user bulk executor.
func (executor *BaseExecutor) SetBulkQueryExecutor(be BulkQueryExecutor) {
	executor.BulkQueryExecutor = be
}

// SetErrorHandler sets a user error handler.
func (executor *BaseExecutor) SetErrorHandler(eh ErrorHandler) {
	executor.ErrorHandler = eh
}

// SetSystemQueryExecutor sets a system query executor.
func (executor *BaseExecutor) SetSystemQueryExecutor(sq SystemQueryExecutor) {
	executor.SystemQueryExecutor = sq
}
