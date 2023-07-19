// Copyright (C) 2019 Satoshi Konno. All rights reserved.
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
	"sync"

	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// Table represents a destination or source database of query.
type Table struct {
	sync.Mutex
	Name string
	*query.Schema
	Rows []Row
}

// NewTable returns a new table.
func NewTableWith(name string, schema *query.Schema) *Table {
	tbl := &Table{
		Name:   name,
		Schema: schema,
		Rows:   []Row{},
	}
	return tbl
}

// Select returns rows matched to the specified condition.
func (tbl *Table) Select(cond *query.Condition) ([]Row, error) {
	rows := []Row{}
	return rows, nil
}

// Insert inserts a new row.
func (tbl *Table) Insert(cols []*query.Column) error {
	row := NewRowWith(cols)
	tbl.Lock()
	tbl.Rows = append(tbl.Rows, row)
	tbl.Unlock()
	return nil
}

// String returns the string representation.
func (tbl *Table) String() string {
	return tbl.Schema.String()
}
