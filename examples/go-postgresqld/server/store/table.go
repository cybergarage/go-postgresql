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
	"sync"

	"github.com/cybergarage/go-sqlparser/sql/query"
)

// Table represents a destination or source database of query.
type Table struct {
	sync.Mutex
	Name string
	query.Schema
	Rows []Row
}

// NewTable returns a new table.
func NewTableWith(name string, schema query.Schema) *Table {
	tbl := &Table{
		Name:   name,
		Schema: schema,
		Rows:   []Row{},
	}
	return tbl
}

// Select returns rows matched to the specified condition.
func (tbl *Table) Select(cond query.Condition) ([]Row, error) {
	tbl.Lock()
	defer tbl.Unlock()

	if !cond.HasConditions() {
		return tbl.Rows, nil
	}

	rows := []Row{}
	for _, row := range tbl.Rows {
		if !row.IsMatched(cond) {
			continue
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// Insert inserts a new row.
func (tbl *Table) Insert(cols []query.Column) error {
	row := NewRowWith(tbl, cols)
	tbl.Lock()
	defer tbl.Unlock()
	tbl.Rows = append(tbl.Rows, row)
	return nil
}

// Update updates rows matched to the specified condition.
func (tbl *Table) Update(cols []query.Column, cond query.Condition) (int, error) {
	rows, err := tbl.Select(cond)
	if err != nil {
		return 0, err
	}

	tbl.Lock()
	defer tbl.Unlock()

	for _, row := range rows {
		row.Update(cols)
	}

	return len(rows), nil
}

// Delete deletes rows matched to the specified condition.
func (tbl *Table) Delete(cond query.Condition) (int, error) {
	rows, err := tbl.Select(cond)
	if err != nil {
		return 0, err
	}

	tbl.Lock()
	defer tbl.Unlock()

	for _, row := range rows {
		for n, r := range tbl.Rows {
			if !row.IsEqual(r) {
				continue
			}
			tbl.Rows = append(tbl.Rows[:n], tbl.Rows[n+1:]...)
		}
	}

	return len(rows), nil
}
