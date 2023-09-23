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

package query

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Statement represents a statement instance.
type Statement struct {
	query.Statement
}

// NewStatement returns a new statement.
func NewStatement(stmt query.Statement) *Statement {
	return &Statement{
		Statement: stmt,
	}
}

// Bind binds the statement with the specified parameters.
func (stmt *Statement) Bind(bindParams message.BindParams) error {
	updateBindColumns := func(columns []*query.Column, params message.BindParams) error {
		for _, column := range columns {
			if !column.HasLiteral() {
				continue
			}
			v, ok := column.Value().(*query.BindParam)
			if !ok {
				continue
			}
			bindParam, err := params.FindBindParam(v.Name())
			if err != nil {
				return err
			}
			column.SetValue(bindParam.Value)
		}
		return nil
	}

	switch stmt := stmt.Statement.(type) {
	case *query.Insert:
		err := updateBindColumns(stmt.Columns(), bindParams)
		if err != nil {
			return err
		}
	case *query.Update:
		err := updateBindColumns(stmt.Columns(), bindParams)
		if err != nil {
			return err
		}
	}
	return nil
}
