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
)

// Table represents a destination or source database of query.
type Table struct {
	sync.Mutex
	value string
}

// NewTableWithNameAndSchema returns a new database with the specified string.
func NewTableWithNameAndSchema(name string) *Table {
	tbl := &Table{
		value: name,
	}
	return tbl
}

// NewTable returns a new database.
func NewTable() *Table {
	return NewTableWithNameAndSchema("")
}

// Name returns the table name.
func (tbl *Table) Name() string {
	return tbl.value
}

// String returns the string representation.
func (tbl *Table) String() string {
	return tbl.value
}