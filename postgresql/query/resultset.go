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
func NewResponseFromResultSet(rs resultset.ResultSet) (protocol.Responses, error) {
	// Schema

	schema := rs.Schema()
	selectors := schema.Selectors()

	// Responses

	res := protocol.NewResponses()

	// Row description response

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
	for rs.Next() {
		rsRow, err := rs.Row()
		if err != nil {
			return nil, err
		}
		rowObj := rsRow.Object()
		dataRow, err := NewDataRowForSelectors(schema, rowDesc, selectors, rowObj)
		if err != nil {
			return nil, err
		}
		res = res.Append(dataRow)
		nRows++
	}

	cmpRes, err := protocol.NewSelectCompleteWith(nRows)
	if err != nil {
		return nil, err
	}
	res = res.Append(cmpRes)

	return res, nil
}
