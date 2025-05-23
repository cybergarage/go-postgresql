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

package query

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Statement represents a statement instance.
type Statement struct {
	obj query.Statement
}

// NewStatementWith returns a new statement.
func NewStatementWith(stmt query.Statement) *Statement {
	return &Statement{
		obj: stmt,
	}
}

// Statement returns a statement object.
func (stmt *Statement) Object() query.Statement {
	return stmt.obj
}

// Bind binds the statement with the specified parameters.
func (stmt *Statement) Bind(bindParams protocol.BindParams) error {
	updateBindColumns := func(columns []query.Column, params protocol.BindParams) error {
		for _, column := range columns {
			if !column.HasValue() {
				continue
			}
			v, ok := column.Value().(query.BindParam)
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

	switch stmt := stmt.obj.(type) {
	case query.Update:
		err := updateBindColumns(stmt.Columns(), bindParams)
		if err != nil {
			return err
		}
	case query.Insert:
		for _, colums := range stmt.Values() {
			err := updateBindColumns(colums, bindParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
