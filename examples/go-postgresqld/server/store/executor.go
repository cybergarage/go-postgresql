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
	"fmt"

	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// CreateDatabase handles a CREATE DATABASE query.
func (store *MemStore) CreateDatabase(conn *postgresql.Conn, q *query.CreateDatabase) (message.Responses, error) {
	dbName := q.DatabaseName()

	_, ok := store.GetDatabase(dbName)
	if ok && !q.IfNotExists() {
		return nil, postgresql.NewErrDatabaseExist(dbName)
	}

	err := store.AddDatabase(NewDatabaseWithName(dbName))
	if err != nil {
		return nil, err
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// CreateTable handles a CREATE TABLE query.
func (store *MemStore) CreateTable(conn *postgresql.Conn, q *query.CreateTable) (message.Responses, error) {
	dbName := conn.DatabaseName()

	db, ok := store.GetDatabase(dbName)
	if !ok {
		return nil, postgresql.NewErrDatabaseExist(dbName)
	}

	tblName := q.TableName()
	_, ok = db.GetTable(tblName)
	if ok && !q.IfNotExists() {
		return nil, postgresql.NewErrTableNotExist(tblName)
	}

	tbl := NewTableWith(tblName, q.Schema())
	err := db.AddTable(tbl)
	if err != nil {
		return nil, err
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// CreateIndex handles a CREATE INDEX query.
func (store *MemStore) CreateIndex(conn *postgresql.Conn, q *query.CreateIndex) (message.Responses, error) {
	return nil, postgresql.NewErrNotImplemented("CREATE INDEX")
}

// DropDatabase handles a DROP DATABASE query.
func (store *MemStore) DropDatabase(conn *postgresql.Conn, q *query.DropDatabase) (message.Responses, error) {
	dbName := q.DatabaseName()

	db, ok := store.GetDatabase(dbName)
	if !ok && !q.IfExists() {
		return nil, postgresql.NewErrDatabaseNotExist(dbName)
	}

	err := store.Databases.DropDatabase(db)
	if err != nil {
		return nil, err
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// DropIndex handles a DROP INDEX query.
func (store *MemStore) DropTable(conn *postgresql.Conn, q *query.DropTable) (message.Responses, error) {
	dbName := conn.DatabaseName()

	db, ok := store.GetDatabase(dbName)
	if !ok {
		return nil, postgresql.NewErrDatabaseNotExist(dbName)
	}

	tblName := q.TableName()
	tbl, ok := db.GetTable(tblName)
	if !ok && !q.IfExists() {
		return nil, postgresql.NewErrTableNotExist(tblName)
	}

	err := db.DropTable(tbl)
	if err != nil {
		return nil, err
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// Insert handles a INSERT query.
func (store *MemStore) Insert(conn *postgresql.Conn, q *query.Insert) (message.Responses, error) {
	_, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), q.TableName())
	if err != nil {
		return nil, err
	}

	cols := q.Columns()
	err = tbl.Insert(cols)
	if err != nil {
		return nil, err
	}

	return message.NewInsertCompleteResponsesWith(1)
}

// Select handles a SELECT query.
func (store *MemStore) Select(conn *postgresql.Conn, q *query.Select) (message.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// 55.2.2. Simple Query
	// https://www.postgresql.org/docs/16/protocol-flow.html#PROTOCOL-FLOW-SIMPLE-QUERY
	tbls := q.Tables()
	if len(tbls) != 1 {
		return nil, postgresql.NewErrNotImplemented(fmt.Sprintf("Multiple tables (%v)", tbls.String()))
	}
	tblName := tbls[0].TableName()

	_, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), tblName)
	if err != nil {
		return nil, err
	}

	rows, err := tbl.Select(q.Where())
	if err != nil {
		return nil, err
	}

	colums := q.Columns()
	if colums.IsSelectAll() {
		colums = tbl.Columns()
	}

	// Row description response

	schema := tbl.Schema
	names := colums.Names()

	rowDesc := message.NewRowDescription()
	for n, name := range names {
		_, err := schema.ColumnByName(name)
		if err != nil {
			return nil, err
		}
		field := message.NewRowFieldWith(name,
			message.WithNumber(int16(n+1)),
		)
		rowDesc.AppendField(field)
	}

	// Data row response

	dataRow := message.NewDataRow()
	for _, row := range rows {
		for _, name := range names {
			v, err := row.ValueByColumnName(name)
			if err != nil {
				return nil, err
			}
			dataRow.AppendData(v)
		}
	}

	return message.NewResponsesWith(rowDesc, dataRow), postgresql.NewErrNotImplemented("SELECT")
}

// Update handles a UPDATE query.
func (store *MemStore) Update(conn *postgresql.Conn, q *query.Update) (message.Responses, error) {
	return nil, postgresql.NewErrNotImplemented("UPDATE")
}

// Delete handles a DELETE query.
func (store *MemStore) Delete(conn *postgresql.Conn, q *query.Delete) (message.Responses, error) {
	return nil, postgresql.NewErrNotImplemented("DELETE")
}
