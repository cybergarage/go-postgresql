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
	"fmt"
	"reflect"
	"time"

	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-sqlparser/sql/errors"
	"github.com/cybergarage/go-sqlparser/sql/fn"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Row represents a row of a table.
type Row map[string]any

// NewRow returns a new row.
func NewRow() Row {
	return make(Row)
}

// NewRowWithResultMap returns a new row with the specified result map.
func NewRowWithResultMap(resultMap map[string]any) Row {
	row := NewRow()
	for k, v := range resultMap {
		row[k] = v
	}
	return row
}

// NewRowFromColumns returns a new row with the specified columns.
func NewRowFromColumns(table *Table, cols query.Columns) (Row, error) {
	row := NewRow()

	for _, col := range cols {
		colName := col.Name()
		schemaCol, err := table.Schema.LookupColumn(colName)
		if err != nil {
			return nil, err
		}
		var colValue any
		switch schemaCol.DataType() {
		case query.BooleanType:
			var v bool
			err = safecast.ToBool(col.Value(), &v)
			colValue = v
		case query.TextType, query.VarCharType, query.LongTextType:
			var v string
			err = safecast.ToString(col.Value(), &v)
			colValue = v
		case query.IntType, query.IntegerType, query.TinyIntType, query.SmallIntType, query.MediumIntType:
			var v int
			err = safecast.ToInt(col.Value(), &v)
			colValue = v
		case query.SerialType, query.BigSerialType, query.SmallSerialType:
			var v int
			err = safecast.ToInt(col.Value(), &v)
			colValue = v
		case query.FloatType:
			var v float32
			err = safecast.ToFloat32(col.Value(), &v)
			colValue = v
		case query.DoubleType:
			var v float64
			err = safecast.ToFloat64(col.Value(), &v)
			colValue = v
		case query.DateTimeType, query.TimeStampType:
			var v time.Time
			err = safecast.ToTime(col.Value(), &v)
			colValue = v
		}
		if err != nil {
			return nil, err
		}
		row[colName] = colValue
	}
	return row, nil
}

// IsMatched returns true if the row is matched with the specified condition.
func (row Row) IsMatched(cond query.Condition) bool {
	if !cond.HasConditions() {
		return true
	}

	eq := func(name string, v any) bool {
		rv, ok := row[name]
		if !ok {
			return false
		}
		return safecast.Equal(rv, v)
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
func (row Row) Update(colums []query.Column) error {
	for _, col := range colums {
		colName := col.Name()
		if col.HasValue() {
			row[colName] = col.Value()
		}
	}
	for _, col := range colums {
		colName := col.Name()
		if fx, ok := col.Function(); ok {
			if exe, err := fx.Executor(); err == nil {
				v, err := exe.Execute(fn.NewMapWithMap(row))
				if err != nil {
					return err
				}
				row[colName] = v
			}
		}
	}
	return nil
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
		return nil, fmt.Errorf("row (%s) %w", name, errors.ErrNotExist)
	}
	return v, nil
}
