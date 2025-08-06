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
	"github.com/cybergarage/go-sqlparser/sql/fn"
	"github.com/cybergarage/go-sqlparser/sql/query"
	"github.com/cybergarage/go-sqlparser/sql/query/response/resultset"
)

// Row represents a row of a result set.
type Row = map[string]any

// NewDataRowForSelectors returns a new DataRow from the specified row.
func NewDataRowForSelectors(schema resultset.Schema, rowDesc *protocol.RowDescription, selectors query.Selectors, row Row) (*protocol.DataRow, error) {
	dataRow := protocol.NewDataRow()

	for n, selector := range selectors {
		field := rowDesc.Field(n)
		name := selector.String()

		v, ok := row[name]
		if !ok {
			dataRow.AppendData(field, nil)
			continue
		}

		dataRow.AppendData(field, v)
	}

	return dataRow, nil
}

// NewDataRowsForAggregator returns a new DataRow list from the specified rows.
func NewDataRowsForAggregator(schema resultset.Schema, rowDesc *protocol.RowDescription, selectors query.Selectors, rows []Row, groupBy string) ([]*protocol.DataRow, error) {
	// Sets aggregate functions
	aggrSet, err := selectors.Aggregators()
	if err != nil {
		return nil, err
	}

	err = aggrSet.Reset(fn.GroupBy(groupBy))
	if err != nil {
		return nil, err
	}

	// Executes aggregate functions

	for _, row := range rows {
		err := aggrSet.Aggregate(row)
		if err != nil {
			return nil, err
		}
	}

	// Finalizes aggregate functions

	resultSet, err := aggrSet.Finalize()
	if err != nil {
		return nil, err
	}

	// Creates DataRow list from the result set

	dataRows := []*protocol.DataRow{}

	for resultSet.Next() {
		m, err := resultSet.Map()
		if err != nil {
			return nil, err
		}

		dataRow, err := NewDataRowForSelectors(schema, rowDesc, selectors, m)
		if err != nil {
			return nil, err
		}

		dataRows = append(dataRows, dataRow)
	}

	return dataRows, nil
}
