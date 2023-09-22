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
)

// BaseExtendedQueryExecutor represents a base extended query message executor.
type BaseExtendedQueryExecutor struct {
}

// NewBaseExtendedQueryExecutor returns a base extended query message executor.
func NewBaseExtendedQueryExecutor() *BaseExtendedQueryExecutor {
	return &BaseExtendedQueryExecutor{}
}

// Prepare handles a parse message.
func (executor *BaseExtendedQueryExecutor) Parse(conn *Conn, msg *message.Parse) (message.Responses, error) {
	err := conn.SetPreparedStatement(msg)
	if err != nil {
		return nil, err
	}
	return message.NewResponsesWith(message.NewParseComplete()), nil
}

// Bind handles a bind message.
func (executor *BaseExtendedQueryExecutor) Bind(conn *Conn, msg *message.Bind) (message.Responses, *message.Query, error) {
	preparedQuery, err := conn.PreparedStatement(msg.Name)
	if err != nil {
		return nil, nil, err
	}

	q, err := message.NewQueryWith(preparedQuery, msg)
	if err != nil {
		return nil, nil, err
	}

	return message.NewResponsesWith(message.NewBindComplete()), q, nil
}

// Describe handles a describe message.
func (executor *BaseExtendedQueryExecutor) Describe(conn *Conn, msg *message.Describe) (message.Responses, error) {
	if msg.IsPortal() {
		return message.NewResponsesWith(
			message.NewNoData()), nil
	}

	_, err := conn.PreparedStatement(msg.Name)
	if err != nil {
		return nil, err
	}

	return message.NewResponsesWith(
		message.NewParameterDescription(),
		message.NewNoData()), nil
}
