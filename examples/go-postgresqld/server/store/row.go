// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-sqlparser/sql/errors"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Row represents a row of a table.
type Row map[string]any

// NewRow returns a new row.
func NewRow() Row {
	return make(Row)
}

// NewRowWith returns a new row with the specified columns.
func NewRowWith(table *Table, cols query.Columns) Row {
	row := NewRow()
	for _, schemaCols := range table.Schema.Columns() {
		var colValue any
		colName := schemaCols.Name()
		col, err := cols.LookupColumn(colName)
		if err == nil {
			colValue = col.Value()
		} else {
			colValue = nil
		}
		row[colName] = colValue
	}
	return row
}

// IsMatched returns true if the row is matched with the specified condition.
func (row Row) IsMatched(cond query.Condition) bool {
	if !cond.HasConditions() {
		return true
	}

	deepEqual := func(r1 any, r2 any) bool {
		switch v1 := r1.(type) {
		case string:
			var v2 string
			err := safecast.ToString(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case int:
			var v2 int
			err := safecast.ToInt(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case int8:
			var v2 int8
			err := safecast.ToInt8(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case int16:
			var v2 int16
			err := safecast.ToInt16(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case int32:
			var v2 int32
			err := safecast.ToInt32(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case float32:
			var v2 float32
			err := safecast.ToFloat32(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case float64:
			var v2 float64
			err := safecast.ToFloat64(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		case bool:
			var v2 bool
			err := safecast.ToBool(r2, &v2)
			if err == nil {
				if v1 == v2 {
					return true
				}
			}
		}
		return reflect.DeepEqual(r1, r2)
	}

	eq := func(name string, v any) bool {
		rv, ok := row[name]
		if !ok {
			return false
		}
		return deepEqual(rv, v)
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
func (row Row) Update(colums []query.Column) {
	for _, col := range colums {
		colName := col.Name()
		if fn, ok := col.IsFunction(); ok {
			v, err := fn.Execute(col, row)
			if err != nil {
				continue
			}
			row[colName] = v
		} else {
			if col.HasValue() {
				row[colName] = col.Value()
			}
		}
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
		return nil, errors.NewErrNotExis(name)
	}
	return v, nil
}
