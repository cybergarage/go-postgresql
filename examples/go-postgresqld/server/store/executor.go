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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

import (
	"fmt"

	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// Begin handles a BEGIN query.
func (store *MemStore) Begin(*postgresql.Conn, *query.Begin) (message.Responses, error) {
	return nil, query.NewErrNotImplemented("BEGIN")
}

// Commit handles a COMMIT query.
func (store *MemStore) Commit(*postgresql.Conn, *query.Commit) (message.Responses, error) {
	return nil, query.NewErrNotImplemented("COMMIT")
}

// Rollback handles a ROLLBACK query.
func (store *MemStore) Rollback(*postgresql.Conn, *query.Rollback) (message.Responses, error) {
	return nil, query.NewErrNotImplemented("ROLLBACK")
}

// CreateDatabase handles a CREATE DATABASE query.
func (store *MemStore) CreateDatabase(conn *postgresql.Conn, q *query.CreateDatabase) (message.Responses, error) {
	dbName := q.DatabaseName()

	_, ok := store.GetDatabase(dbName)
	if ok {
		if q.IfNotExists() {
			return message.NewCommandCompleteResponsesWith(q.String())
		}
		return nil, query.NewErrDatabaseExist(dbName)
	}

	err := store.AddDatabase(NewDatabaseWithName(dbName))
	if err != nil {
		return nil, err
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// CreateTable handles a CREATE TABLE query.
func (store *MemStore) CreateTable(conn *postgresql.Conn, q *query.CreateTable) (message.Responses, error) {
	dbName := conn.Database()

	db, ok := store.GetDatabase(dbName)
	if !ok {
		return nil, query.NewErrDatabaseExist(dbName)
	}

	tblName := q.TableName()
	_, ok = db.GetTable(tblName)
	if ok {
		if q.IfNotExists() {
			return message.NewCommandCompleteResponsesWith(q.String())
		}
		return nil, query.NewErrTableNotExist(tblName)
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
	return nil, query.NewErrNotImplemented("CREATE INDEX")
}

// DropDatabase handles a DROP DATABASE query.
func (store *MemStore) DropDatabase(conn *postgresql.Conn, q *query.DropDatabase) (message.Responses, error) {
	dbName := q.DatabaseName()

	db, ok := store.GetDatabase(dbName)
	if !ok {
		if q.IfExists() {
			return message.NewCommandCompleteResponsesWith(q.String())
		}
		return nil, query.NewErrDatabaseNotExist(dbName)
	}

	err := store.Databases.DropDatabase(db)
	if err != nil {
		return nil, err
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// DropIndex handles a DROP INDEX query.
func (store *MemStore) DropTable(conn *postgresql.Conn, q *query.DropTable) (message.Responses, error) {
	for _, dropTbl := range q.Tables() {
		db, tbl, err := store.GetDatabaseTable(conn, conn.Database(), dropTbl.TableName())
		if err != nil {
			if q.IfExists() {
				continue
			}
			return nil, err
		}
		err = db.DropTable(tbl)
		if err != nil {
			return nil, err
		}
	}

	return message.NewCommandCompleteResponsesWith(q.String())
}

// Insert handles a INSERT query.
func (store *MemStore) Insert(conn *postgresql.Conn, q *query.Insert) (message.Responses, error) {
	_, tbl, err := store.GetDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	err = q.SetSchema(tbl.Schema)
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
	tbls := q.Tables()
	if len(tbls) != 1 {
		return nil, query.NewErrNotImplemented(fmt.Sprintf("Multiple tables (%v)", tbls.String()))
	}
	tblName := tbls[0].TableName()

	_, tbl, err := store.GetDatabaseTable(conn, conn.Database(), tblName)
	if err != nil {
		return nil, err
	}

	rows, err := tbl.Select(q.Where())
	if err != nil {
		return nil, err
	}

	// Responses

	schema := tbl.Schema
	res := message.NewResponses()

	// Row description response

	selectors := q.Selectors()
	if selectors.IsSelectAll() {
		selectors = tbl.Selectors()
	}

	rowDesc := message.NewRowDescription()
	for n, selector := range selectors {
		field, err := query.NewRowFieldFrom(schema, selector, n)
		if err != nil {
			return nil, err
		}
		rowDesc.AppendField(field)
	}
	res = res.Append(rowDesc)

	// Data row response

	nDataRow := 0
	if !selectors.HasAggregateFunction() {
		for _, row := range rows {
			dataRow, err := query.NewDataRowForSelectors(schema, rowDesc, selectors, row)
			if err != nil {
				return nil, err
			}
			res = res.Append(dataRow)
			nDataRow++
		}
	} else {
		groupBy := q.GroupBy().Column()
		queryRows := []query.Row{}
		for _, row := range rows {
			queryRows = append(queryRows, row)
		}
		dataRows, err := query.NewDataRowsForAggregateFunction(schema, rowDesc, selectors, queryRows, groupBy)
		if err != nil {
			return nil, err
		}
		for _, dataRow := range dataRows {
			res = res.Append(dataRow)
			nDataRow++
		}
	}

	cmpRes, err := message.NewSelectCompleteWith(nDataRow)
	if err != nil {
		return nil, err
	}
	res = res.Append(cmpRes)

	return res, nil
}

// Update handles a UPDATE query.
func (store *MemStore) Update(conn *postgresql.Conn, q *query.Update) (message.Responses, error) {
	_, tbl, err := store.GetDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	err = q.SetSchema(tbl.Schema)
	if err != nil {
		return nil, err
	}

	n, err := tbl.Update(q.Columns(), q.Where())
	if err != nil {
		return nil, err
	}

	return message.NewUpdateCompleteResponsesWith(n)
}

// Delete handles a DELETE query.
func (store *MemStore) Delete(conn *postgresql.Conn, q *query.Delete) (message.Responses, error) {
	_, tbl, err := store.GetDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	n, err := tbl.Delete(q.Where())
	if err != nil {
		return nil, err
	}

	return message.NewDeleteCompleteResponsesWith(n)
}

// Copy handles a COPY query.
func (store *MemStore) Copy(conn *postgresql.Conn, q *query.Copy, stream *postgresql.CopyStream) (message.Responses, error) {
	_, tbl, err := store.GetDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	// PostgreSQL: Documentation: 16: COPY
	// https://www.postgresql.org/docs/16/sql-copy.html

	newQueryWith := func(schema *query.Schema, stream *postgresql.CopyStream) (*query.Insert, error) {
		copyData, err := stream.CopyData()
		if err != nil {
			return nil, err
		}
		copyColums := copyData.Data
		schemaColumns := schema.Columns()
		// COPY FROM will raise an error if any line of the input file contains
		//  more or fewer columns than are expected.
		if len(copyColums) != len(schemaColumns) {
			return nil, query.NewErrColumnsNotEqual(len(copyColums), len(schemaColumns))
		}
		columns := schemaColumns.Copy()
		for idx, column := range columns {
			v := copyColums[idx]
			if err := column.SetValue(v); err != nil {
				return nil, err
			}
		}

		return query.NewInsertWith(schema.SchemaTable(), columns), nil
	}

	copyData := func(schema *query.Schema, stream *postgresql.CopyStream) error {
		q, err := newQueryWith(schema, stream)
		if err != nil {
			return err
		}
		_, err = store.Insert(conn, q)
		return err
	}

	schema := tbl.Schema
	nCopy := 0
	nFail := 0
	ok, err := stream.Next()
	for {
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if err := copyData(schema, stream); err == nil {
			nCopy++
		} else {
			nFail++
		}
		ok, err = stream.Next()
	}

	return message.NewCopyCompleteResponsesWith(nCopy)
}
