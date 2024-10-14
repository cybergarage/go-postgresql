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

import (
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

type MemStore struct {
	Databases
	*postgresql.BaseExecutor
}

// NewMemStore returns an in-memory storeinstance.
func NewMemStore() *MemStore {
	store := &MemStore{
		Databases:    NewDatabases(),
		BaseExecutor: postgresql.NewBaseExecutor(),
	}
	return store
}

func (store *MemStore) LookupDatabaseTable(conn postgresql.Conn, dbName string, tblName string) (*Database, *Table, error) {
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, nil, query.NewErrDatabaseNotExist(dbName)
	}

	tbl, ok := db.GetTable(tblName)
	if !ok {
		return nil, nil, query.NewErrTableNotExist(tblName)
	}

	return db, tbl, nil
}
