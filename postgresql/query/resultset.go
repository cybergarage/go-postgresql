// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-sqlparser/sql/query/response/resultset"
)

// NewResponseFromResultSet creates a response from a result set.
func NewResponseFromResultSet(stmt Select, rs resultset.ResultSet) (protocol.Responses, error) {
	// Schema

	schema := rs.Schema()

	// Responses

	res := protocol.NewResponses()

	// Row description response

	selectors := stmt.Selectors()
	if selectors.IsAsterisk() {
		selectors = schema.Selectors()
	}

	rowDesc := protocol.NewRowDescription()
	for n, selector := range selectors {
		field, err := NewRowFieldFrom(schema, selector, n)
		if err != nil {
			return nil, err
		}
		rowDesc.AppendField(field)
	}
	res = res.Append(rowDesc)

	// Data row response

	nRows := 0
	if !selectors.HasAggregateFunction() {
		offset := stmt.Limit().Offset()
		limit := stmt.Limit().Limit()
		rowNo := 0
		for rs.Next() {
			rowNo++
			if 0 < offset && rowNo <= offset {
				continue
			}
			rowObj := rs.Row().Object()
			dataRow, err := NewDataRowForSelectors(schema, rowDesc, selectors, rowObj)
			if err != nil {
				return nil, err
			}
			res = res.Append(dataRow)
			nRows++
			if 0 < limit && limit <= nRows {
				break
			}
		}
	} else {
		groupBy := stmt.GroupBy().ColumnName()
		queryRows := []Row{}
		for rs.Next() {
			rowObj := rs.Row().Object()
			queryRows = append(queryRows, rowObj)
		}
		dataRows, err := NewDataRowsForAggregateFunction(schema, rowDesc, selectors, queryRows, groupBy)
		if err != nil {
			return nil, err
		}
		for _, dataRow := range dataRows {
			res = res.Append(dataRow)
			nRows++
		}
	}

	cmpRes, err := protocol.NewSelectCompleteWith(nRows)
	if err != nil {
		return nil, err
	}
	res = res.Append(cmpRes)

	return res, nil
}
