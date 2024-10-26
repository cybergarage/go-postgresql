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

package store

// Databases represents a collection of databases.
type Databases map[string]*Database

// NewDatabases returns a databases instance.
func NewDatabases() Databases {
	return Databases{}
}

// AddDatabase adds a specified database.
func (dbs Databases) AddDatabase(db *Database) error {
	dbName := db.Name()
	dbs[dbName] = db
	return nil
}

// DropDatabase remove the specified database.
func (dbs Databases) DropDatabase(db *Database) error {
	name := db.Name()
	delete(dbs, name)
	return nil
}

// LookupDatabase returns a database with the specified name.
func (dbs Databases) LookupDatabase(name string) (*Database, bool) {
	ks, ok := dbs[name]
	return ks, ok
}

// GetTableWithDatabase returns a specified table in a specified database.
func (dbs *Databases) GetTableWithDatabase(dbName string, tableName string) (*Table, bool) {
	db, ok := dbs.LookupDatabase(dbName)
	if !ok {
		return nil, false
	}
	return db.LookupTable(tableName)
}
