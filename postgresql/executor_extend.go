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

// BaseExtendedQueryExecutor represents a base extended query message executor.
type BaseExtendedQueryExecutor struct {
}

// NewBaseExtendedQueryExecutor returns a base extended query message executor.
func NewBaseExtendedQueryExecutor() *BaseExtendedQueryExecutor {
	return &BaseExtendedQueryExecutor{}
}

// Prepare handles a parse message.
func (executor *BaseExtendedQueryExecutor) Parse(*Conn, *message.Parse) (message.Responses, error) {
	return nil, query.NewErrNotImplemented("Parse")
}

// Bind handles a bind message.
func (executor *BaseExtendedQueryExecutor) Bind(*Conn, *message.Bind) (message.Responses, error) {
	return nil, query.NewErrNotImplemented("Bind")
}

// Describe handles a describe message.
func (executor *BaseExtendedQueryExecutor) Describe(*Conn, *message.Describe) (message.Responses, error) {
	return nil, query.NewErrNotImplemented("Describe")
}
