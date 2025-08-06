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

// Database represents a destination or source database of query.
type Database struct {
	name   string
	tables map[string]*Table
}

// NewDatabaseWithName returns a new database with the specified string.
func NewDatabaseWithName(name string) *Database {
	db := &Database{
		name:   name,
		tables: map[string]*Table{},
	}

	return db
}

// NewDatabase returns a new database.
func NewDatabase() *Database {
	return NewDatabaseWithName("")
}

// Name returns the database name.
func (db *Database) Name() string {
	return db.name
}

// AddTable adds a specified table into the database.
func (db *Database) AddTable(table *Table) {
	tableName := table.Name
	db.tables[tableName] = table
}

// AddTables adds a specified tables into the database.
func (db *Database) AddTables(tables []*Table) {
	for _, table := range tables {
		db.AddTable(table)
	}
}

// DropTable remove the specified table.
func (db *Database) DropTable(table *Table) bool {
	name := table.TableName()
	delete(db.tables, name)
	_, ok := db.tables[name]

	return !ok
}

// LookupTable returns a table with the specified name.
func (db *Database) LookupTable(name string) (*Table, bool) {
	table, ok := db.tables[name]
	return table, ok
}

// Tables returns all tables in the database.
func (db *Database) Tables() []*Table {
	tables := make([]*Table, 0)
	for _, table := range db.tables {
		tables = append(tables, table)
	}

	return tables
}
