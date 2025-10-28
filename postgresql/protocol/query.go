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

package protocol

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

import (
	"fmt"

	"github.com/cybergarage/go-sqlparser/sql/stmt"
)

// Query represents a parse protocol.
type Query struct {
	*RequestMessage
	Query string
	BindParams
	stmt.BindStatement
}

// NewQueryWithReader returns a new query message with specified reader.
func NewQueryWithReader(reader *MessageReader) (*Query, error) {
	msg, err := NewRequestMessageWithReader(reader)
	if err != nil {
		return nil, err
	}
	query, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	q := &Query{
		RequestMessage: msg,
		Query:          query,
		BindParams:     BindParams{},
		BindStatement:  nil,
	}
	q.BindStatement = stmt.NewBindStatement(
		stmt.WithBindStatementQuery(q.Query),
	)
	return q, nil
}

// NewQueryWith returns a new query message with specified parameters.
func NewQueryWith(parseMsg *Parse, bindMsg *Bind) (*Query, error) {
	q := &Query{
		RequestMessage: nil,
		Query:          parseMsg.Query,
		BindParams:     bindMsg.Params,
		BindStatement:  nil,
	}
	bindParams := stmt.BindParams{}
	for _, param := range bindMsg.Params {
		bindParams = append(bindParams, stmt.NewBindParam(param.Value))
	}
	q.BindStatement = stmt.NewBindStatement(
		stmt.WithBindStatementQuery(q.Query),
		stmt.WithBindStatementParams(bindParams),
	)
	return q, nil
}

// Statement returns the bind statement of the query.
func (q *Query) Statements() ([]stmt.Statement, error) {
	return q.BindStatement.Statements()
}

// String returns the string representation of the query.
func (q *Query) String() string {
	s := q.Query
	if 0 < len(q.BindParams) {
		s += " WITH BIND PARAMS: "
		for i, param := range q.BindParams {
			if 0 < i {
				s += ", "
			}
			s += fmt.Sprintf("%s", param.Value)
		}
	}
	return s
}
