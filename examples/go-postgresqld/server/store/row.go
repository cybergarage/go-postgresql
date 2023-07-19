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
	"reflect"

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

// IsMatched returns true if the row is matched with the specified condition.
func (row Row) IsMatched(cond *query.Condition) bool {
	if cond.IsEmpty() {
		return true
	}

	eq := func(name string, v any) bool {
		rv, ok := row[name]
		if !ok {
			return false
		}
		return reflect.DeepEqual(rv, v)
	}

	expr := cond.Expr()
	switch expr := expr.(type) {
	case *query.CmpExpr:
		name := expr.Left().Name()
		value := expr.Right().Value()
		switch expr.Operator() {
		case query.EQ:
			return eq(name, value)
		case query.NEQ:
			return !eq(name, value)
		default:
			return false
		}
	}

	return true
}

// Update updates the row with the specified columns.
func (row Row) Update(colums []*query.Column) {
	for _, col := range colums {
		row[col.Name()] = col.Value()
	}
}

// IsEqual returns true if the row is equal to the specified row.
func (row Row) IsEqual(other Row) bool {
	if len(row) != len(other) {
		return false
	}

	for k, v := range row {
		if ov, ok := other[k]; !ok || !reflect.DeepEqual(v, ov) {
			return false
		}
	}

	return true
}

// ValueByName returns a value of the specified column name.
func (row Row) ValueByName(name string) (any, error) {
	v, ok := row[name]
	if !ok {
		return nil, postgresql.NewErrNotExist(name)
	}
	return v, nil
}
