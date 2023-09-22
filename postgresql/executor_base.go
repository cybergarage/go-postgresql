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
	Authenticator
	StartupHandler
	QueryExecutor
	QueryExtraExecutor
	TransactionExecutor
	ExtendedQueryExecutor
	BulkExecutor
	ErrorHandler
}

// NewBaseExecutor returns a base frontend message executor.
func NewBaseExecutor() *BaseExecutor {
	executor := &BaseExecutor{
		Authenticator:         NewBaseAuthenticator(),
		StartupHandler:        NewBaseProtocolExecutor(),
		QueryExecutor:         NewBaseQueryExecutor(),
		QueryExtraExecutor:    nil,
		ExtendedQueryExecutor: nil,
		TransactionExecutor:   NewBaseTransactionExecutor(),
		BulkExecutor:          NewBaseBulkExecutor(),
		ErrorHandler:          NewBaseErrorHandler(),
	}
	executor.QueryExtraExecutor = NewBaseSugarExecutorWith(executor)
	executor.ExtendedQueryExecutor = NewBaseExtendedQueryExecutorWith(executor)
	return executor
}

// SetAuthenticator sets a user authenticator.
func (executor *BaseExecutor) SetAuthenticator(at Authenticator) {
	executor.Authenticator = at
}

// SetStartupHandler sets a user startup handler.
func (executor *BaseExecutor) SetStartupHandler(sh StartupHandler) {
	executor.StartupHandler = sh
}

// SetQueryExecutor sets a user query executor.
func (executor *BaseExecutor) SetQueryExecutor(qe QueryExecutor) {
	executor.QueryExecutor = qe
}

// SetQueryExtraExecutor sets a user query extra executor.
func (executor *BaseExecutor) SetQueryExtraExecutor(qe QueryExtraExecutor) {
	executor.QueryExtraExecutor = qe
}

// SetTransactionExecutor sets a user transaction executor.
func (executor *BaseExecutor) SetTransactionExecutor(te TransactionExecutor) {
	executor.TransactionExecutor = te
}

// SetExtendedQueryExecutor sets a user extended query executor.
func (executor *BaseExecutor) SetExtendedQueryExecutor(eqe ExtendedQueryExecutor) {
	executor.ExtendedQueryExecutor = eqe
}

// SetBulkExecutor sets a user bulk executor.
func (executor *BaseExecutor) SetBulkExecutor(be BulkExecutor) {
	executor.BulkExecutor = be
}

// SetErrorHandler sets a user error handler.
func (executor *BaseExecutor) SetErrorHandler(eh ErrorHandler) {
	executor.ErrorHandler = eh
}
