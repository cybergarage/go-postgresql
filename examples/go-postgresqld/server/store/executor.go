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

// CreateDatabase handles a CREATE DATABASE query.
func (store *MemStore) CreateDatabase(conn *postgresql.Conn, q *query.CreateDatabase) (message.Responses, error) {
	dbName := q.DatabaseName()

	_, ok := store.GetDatabase(dbName)
	if ok {
		if q.IfNotExists() {
			return message.NewCommandCompleteResponsesWith(q.String())
		}
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
	if ok {
		if q.IfNotExists() {
			return message.NewCommandCompleteResponsesWith(q.String())
		}
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
	if !ok {
		if q.IfExists() {
			return message.NewCommandCompleteResponsesWith(q.String())
		}
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
	for _, dropTbl := range q.Tables() {
		db, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), dropTbl.TableName())
		if err != nil {
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
	_, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), q.TableName())
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

	// Row description response

	colums := q.Columns()
	if colums.IsSelectAll() {
		colums = tbl.Columns()
	}

	schema := tbl.Schema
	names := colums.Names()

	res := message.NewResponses()

	rowDesc := message.NewRowDescription()
	for n, name := range names {
		schemaColumn, err := schema.ColumnByName(name)
		if err != nil {
			return nil, err
		}
		dt := int32(query.DataTypeFrom(schemaColumn.DataType()))
		fc := int32(query.FormatCodeFrom(schemaColumn.DataType()))
		field := message.NewRowFieldWith(name,
			message.WithNumber(int16(n+1)),
			message.WithDataTypeID(dt),
			message.WithFormatCode(int16(fc)),
		)
		rowDesc.AppendField(field)
	}
	res = res.Append(rowDesc)

	// Data row response

	for _, row := range rows {
		dataRow := message.NewDataRow()
		for _, name := range names {
			v, err := row.ValueByName(name)
			if err != nil {
				dataRow.AppendData(nil)
				continue
			}
			dataRow.AppendData(v)
		}
		res = res.Append(dataRow)
	}

	cmpRes, err := message.NewSelectCompleteWith(len(rows))
	if err != nil {
		return nil, err
	}
	res = res.Append(cmpRes)

	return res, nil
}

// Update handles a UPDATE query.
func (store *MemStore) Update(conn *postgresql.Conn, q *query.Update) (message.Responses, error) {
	_, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), q.TableName())
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
	_, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), q.TableName())
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
	_, tbl, err := store.GetDatabaseTable(conn, conn.DatabaseName(), q.TableName())
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
			return nil, postgresql.NewErrColumnsNotEqual(len(copyColums), len(schemaColumns))
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
