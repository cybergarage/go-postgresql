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

package store

import (
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// Row represents a row of a table.
type Row map[string]any

// NewRow returns a new row.
func NewRow() Row {
	return make(Row)
}

// NewRowWith returns a new row with the specified columns.
func NewRowWith(cols []*query.Column) Row {
	row := NewRow()
	for _, col := range cols {
		row[col.Name()] = col.Value()
	}
	return row
}

func (row Row) IsMatched(cond *query.Condition) bool {
	if cond.IsEmpty() {
		return true
	}

	eq := func(name string, v any) bool {
		return false
	}

	expr := cond.Expr()
	switch expr := expr.(type) {
	case *query.CmpExpr:
		name := expr.Left().Name()
		value := expr.Right().Value()
		switch expr.Operator() {
		case query.EQ:
			return eq(name, value)
		default:
			return false
		}
	}

	return true
}

// ValueByColumnName returns a value of the specified column name.
func (row Row) ValueByColumnName(name string) (any, error) {
	v, ok := row[name]
	if !ok {
		return nil, postgresql.NewErrNotExist(name)
	}
	return v, nil
}
