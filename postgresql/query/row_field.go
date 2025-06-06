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
	"github.com/cybergarage/go-sqlparser/sql"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// NewRowFieldFrom returns a new RowField from the specified selector.
func NewRowFieldFrom(schema sql.ResultSetSchema, selector query.Selector, idx int) (*protocol.RowField, error) {
	columnName := selector.String()
	schemaColumn, err := schema.LookupColumn(columnName)
	if err != nil {
		return nil, err
	}
	dt, err := NewDataTypeFrom(schemaColumn.DataType())
	if err != nil {
		return nil, err
	}
	return protocol.NewRowFieldWith(columnName,
		protocol.WithRowFieldNumber(int16(idx+1)),
		protocol.WithRowFieldDataType(dt),
	), nil
}
