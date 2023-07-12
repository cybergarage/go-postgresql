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

// BaseExecutor represents a base frontend message executor.
type BaseExecutor struct {
	Authenticator
	ProtocolExecutor
	QueryExecutor
}

// NewBaseExecutor returns a base frontend message executor.
func NewBaseExecutor() *BaseExecutor {
	return &BaseExecutor{
		Authenticator:    NewBaseAuthenticator(),
		ProtocolExecutor: NewBaseProtocolExecutor(),
		QueryExecutor:    NewBaseQueryExecutor(),
	}
}

// SetAuthenticator sets a user authenticator.
func (executor *BaseExecutor) SetAuthenticator(at Authenticator) {
	executor.Authenticator = at
}

// SetProtocolExecutor sets a user protocol executor.
func (executor *BaseExecutor) SetProtocolExecutor(pe ProtocolExecutor) {
	executor.ProtocolExecutor = pe
}

// SetQueryExecutor sets a user query executor.
func (executor *BaseExecutor) SetQueryExecutor(qe QueryExecutor) {
	executor.QueryExecutor = qe
}
