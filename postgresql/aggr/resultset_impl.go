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

package aggr

type resultSet struct {
	columns []string
	rows    []Row
	index   int
}

type resultSetOption func(*resultSet)

// WithColumns sets the columns of the result set.
func WithColumns(columns []string) resultSetOption {
	return func(rs *resultSet) {
		rs.columns = columns
	}
}

// WithRows sets the rows of the result set.
func WithRows(rows []Row) resultSetOption {
	return func(rs *resultSet) {
		rs.rows = rows
	}
}

// NewResultSet creates a new result set with the given options.
func NewResultSet(opts ...resultSetOption) ResultSet {
	rs := &resultSet{
		columns: []string{},
		rows:    []Row{},
		index:   -1,
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

// Columns returns the column names of the result set.
func (rs *resultSet) Columns() []string {
	return rs.columns
}

// Next moves to the next row in the result set.
func (rs *resultSet) Next() bool {
	if len(rs.rows) <= (rs.index + 1) {
		return false
	}
	rs.index++
	return true
}

// Row returns the current row of the result set.
func (rs *resultSet) Row() (Row, error) {
	if rs.index < 0 || rs.index >= len(rs.rows) {
		return nil, ErrNoData
	}
	return rs.rows[rs.index], nil
}
