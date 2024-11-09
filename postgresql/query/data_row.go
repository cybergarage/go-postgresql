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
	"fmt"

	"github.com/cybergarage/go-postgresql/postgresql/errors"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Row represents a row of a result set.
type Row = map[string]any

// NewDataRowForSelectors returns a new DataRow from the specified row.
func NewDataRowForSelectors(schema query.Schema, rowDesc *protocol.RowDescription, selectors query.SelectorList, row Row) (*protocol.DataRow, error) {
	dataRow := protocol.NewDataRow()
	for n, selector := range selectors {
		field := rowDesc.Field(n)
		switch selector := selector.(type) {
		case query.Column:
			name := selector.Name()
			v, ok := row[name]
			if !ok {
				dataRow.AppendData(field, nil)
				continue
			}
			dataRow.AppendData(field, v)
		case query.Function:
			executor, err := selector.Executor()
			if err != nil {
				return nil, err
			}
			args := []any{}
			for _, arg := range selector.Arguments() {
				v, ok := row[arg.Name()]
				if !ok {
					v = nil
				}
				args = append(args, v)
			}
			v, err := executor.Execute(args...)
			if err != nil {
				return nil, err
			}
			dataRow.AppendData(field, v)
		}
	}
	return dataRow, nil
}

// NewDataRowsForAggregateFunction returns a new DataRow list from the specified rows.
func NewDataRowsForAggregateFunction(schema query.Schema, rowDesc *protocol.RowDescription, selectors query.SelectorList, rows []Row, groupBy string) ([]*protocol.DataRow, error) {
	// Setups aggregate functions
	aggrFns := []query.Function{}
	aggrExecutors := []*query.AggregateFunction{}
	for _, selector := range selectors {
		fn, ok := selector.(query.Function)
		if !ok {
			continue
		}
		executor, err := fn.Executor()
		if err != nil {
			return nil, err
		}
		aggrExecutor, ok := executor.(*query.AggregateFunction)
		if !ok {
			return nil, fmt.Errorf("invalid aggregate function (%s)", fn.Name())
		}
		aggrFns = append(aggrFns, fn)
		aggrExecutors = append(aggrExecutors, aggrExecutor)
	}
	// Executes aggregate functions
	for _, row := range rows {
		for n, aggrFn := range aggrFns {
			var groupKey any
			groupKey = query.GroupByNone
			if 0 < len(groupBy) {
				v, ok := row[groupBy]
				if !ok {
					return nil, errors.NewErrGroupByColumnValueNotFound(groupBy)
				}
				groupKey = v
			}
			args := []any{
				groupKey,
			}
			for _, arg := range aggrFn.Arguments() {
				if arg.IsAsterisk() {
					args = append(args, arg.Name())
					continue
				}
				v, ok := row[arg.Name()]
				if !ok {
					return nil, errors.NewErrColumnValueNotExist(arg.Name())
				}
				args = append(args, v)
			}
			_, err := aggrExecutors[n].Execute(args...)
			if err != nil {
				return nil, err
			}
		}
	}
	// Gets aggregate group keys
	aggrResultSets := map[string]query.AggregateResultSet{}
	groupKeys := []any{}
	for _, aggaggrExecutor := range aggrExecutors {
		aggResultSet := aggaggrExecutor.ResultSet()
		aggrResultSets[aggaggrExecutor.Name()] = aggResultSet
		for aggrResultKey := range aggResultSet {
			hasGroupKey := false
			for _, groupKey := range groupKeys {
				if groupKey == aggrResultKey {
					hasGroupKey = true
				}
			}
			if hasGroupKey {
				continue
			}
			groupKeys = append(groupKeys, aggrResultKey)
		}
	}
	// Add aggregate results
	dataRows := []*protocol.DataRow{}
	if 0 < len(groupKeys) { // ResultSet is not empty
		for _, groupKey := range groupKeys {
			dataRow := protocol.NewDataRow()
			for n, selector := range selectors {
				field := rowDesc.Field(n)
				name := selector.Name()
				switch selector.(type) {
				case query.Column:
					if name != groupBy {
						return nil, fmt.Errorf("invalid column (%s)", name)
					}
					dataRow.AppendData(field, groupKey)
				case query.Function:
					aggResultSet, ok := aggrResultSets[name]
					if !ok {
						return nil, fmt.Errorf("invalid aggregate function (%s)", name)
					}
					aggResult, ok := aggResultSet[groupKey]
					if ok {
						dataRow.AppendData(field, aggResult)
					} else {
						dataRow.AppendData(field, nil)
					}
				}
			}
			dataRows = append(dataRows, dataRow)
		}
	} else { // ResultSet is empty
		dataRow := protocol.NewDataRow()
		for n, selector := range selectors {
			field := rowDesc.Field(n)
			name := selector.Name()
			switch selector.(type) {
			case query.Function:
				switch name {
				case query.CountFunctionName:
					dataRow.AppendData(field, 0)
				default:
					dataRow.AppendData(field, nil)
				}
			}
		}
		dataRows = append(dataRows, dataRow)
	}
	return dataRows, nil
}
