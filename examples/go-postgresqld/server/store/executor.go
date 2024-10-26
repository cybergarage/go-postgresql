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
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
	"github.com/cybergarage/go-postgresql/postgresql/system"
	sql "github.com/cybergarage/go-sqlparser/sql/query"
)

// Begin handles a BEGIN query.
func (store *MemStore) Begin(conn postgresql.Conn, q query.Begin) (protocol.Responses, error) {
	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// Commit handles a COMMIT query.
func (store *MemStore) Commit(con postgresql.Conn, q query.Commit) (protocol.Responses, error) {
	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// Rollback handles a ROLLBACK query.
func (store *MemStore) Rollback(postgresql.Conn, query.Rollback) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("ROLLBACK")
}

// CreateDatabase handles a CREATE DATABASE query.
func (store *MemStore) CreateDatabase(conn postgresql.Conn, q query.CreateDatabase) (protocol.Responses, error) {
	dbName := q.DatabaseName()

	_, ok := store.LookupDatabase(dbName)
	if ok {
		if q.IfNotExists() {
			return protocol.NewCommandCompleteResponsesWith(q.String())
		}
		return nil, query.NewErrDatabaseExist(dbName)
	}

	err := store.AddDatabase(NewDatabaseWithName(dbName))
	if err != nil {
		return nil, err
	}

	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// CreateTable handles a CREATE TABLE query.
func (store *MemStore) CreateTable(conn postgresql.Conn, q query.CreateTable) (protocol.Responses, error) {
	dbName := conn.Database()

	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, query.NewErrDatabaseExist(dbName)
	}

	tblName := q.TableName()
	_, ok = db.LookupTable(tblName)
	if ok {
		if q.IfNotExists() {
			return protocol.NewCommandCompleteResponsesWith(q.String())
		}
		return nil, query.NewErrTableNotExist(tblName)
	}

	tbl := NewTableWith(tblName, q.Schema())
	err := db.AddTable(tbl)
	if err != nil {
		return nil, err
	}

	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// AlterDatabase handles a ALTER DATABASE query.
func (store *MemStore) AlterDatabase(conn postgresql.Conn, q query.AlterDatabase) (protocol.Responses, error) {
	return nil, query.NewErrNotImplemented("ALTER DATABASE")
}

// AlterTable handles a ALTER TABLE query.
func (store *MemStore) AlterTable(conn postgresql.Conn, q query.AlterTable) (protocol.Responses, error) {
	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	schema := tbl.Schema

	if column, ok := q.AddColumn(); ok {
		err := schema.AddColumn(column)
		if err != nil {
			return nil, err
		}
	}

	if index, ok := q.AddIndex(); ok {
		indexColums := sql.NewColumns()
		for _, indexColumn := range index.Columns() {
			schemaColum, err := schema.ColumnByName(indexColumn.Name())
			if err != nil {
				return nil, err
			}
			indexColums = append(indexColums, schemaColum)
		}
		err := schema.AddIndex(sql.NewIndexWith(index.Name(), index.Type(), indexColums))
		if err != nil {
			return nil, err
		}
	}

	if column, ok := q.DropColumn(); ok {
		err := schema.DropColumn(column.Name())
		if err != nil {
			return nil, err
		}
	}

	if _, ok := q.RenameTo(); ok {
		return nil, query.NewErrNotImplemented(q.String())
	}

	if _, _, ok := q.RenameColumns(); ok {
		return nil, query.NewErrNotImplemented(q.String())
	}

	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// DropDatabase handles a DROP DATABASE query.
func (store *MemStore) DropDatabase(conn postgresql.Conn, q query.DropDatabase) (protocol.Responses, error) {
	dbName := q.DatabaseName()

	db, ok := store.LookupDatabase(dbName)
	if !ok {
		if q.IfExists() {
			return protocol.NewCommandCompleteResponsesWith(q.String())
		}
		return nil, query.NewErrDatabaseNotExist(dbName)
	}

	err := store.Databases.DropDatabase(db)
	if err != nil {
		return nil, err
	}

	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// DropIndex handles a DROP INDEX query.
func (store *MemStore) DropTable(conn postgresql.Conn, q query.DropTable) (protocol.Responses, error) {
	for _, dropTbl := range q.Tables() {
		db, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), dropTbl.TableName())
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

	return protocol.NewCommandCompleteResponsesWith(q.String())
}

// Insert handles a INSERT query.
func (store *MemStore) Insert(conn postgresql.Conn, q query.Insert) (protocol.Responses, error) {
	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	cols := q.Columns()
	err = tbl.Insert(cols)
	if err != nil {
		return nil, err
	}

	return protocol.NewInsertCompleteResponsesWith(1)
}

// Select handles a SELECT query.
func (store *MemStore) Select(conn postgresql.Conn, q query.Select) (protocol.Responses, error) {
	from := q.From()
	if len(from) != 1 {
		return nil, query.NewErrMultipleTableNotSupported(from.String())
	}
	tblName := from[0].TableName()

	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), tblName)
	if err != nil {
		return nil, err
	}

	rows, err := tbl.Select(q.Where())
	if err != nil {
		return nil, err
	}

	// Responses

	schema := tbl.Schema
	res := protocol.NewResponses()

	// Row description response

	selectors := q.Selectors()
	if selectors.IsSelectAll() {
		selectors = tbl.Selectors()
	}

	rowDesc := protocol.NewRowDescription()
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
		offset := q.Limit().Offset()
		limit := q.Limit().Limit()
		for rowNo, row := range rows {
			if 0 < offset && rowNo < offset {
				continue
			}
			dataRow, err := query.NewDataRowForSelectors(schema, rowDesc, selectors, row)
			if err != nil {
				return nil, err
			}
			res = res.Append(dataRow)
			nDataRow++
			if 0 < limit && limit <= nDataRow {
				break
			}
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

	cmpRes, err := protocol.NewSelectCompleteWith(nDataRow)
	if err != nil {
		return nil, err
	}
	res = res.Append(cmpRes)

	return res, nil
}

// Update handles a UPDATE query.
func (store *MemStore) Update(conn postgresql.Conn, q query.Update) (protocol.Responses, error) {
	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	n, err := tbl.Update(q.Columns(), q.Where())
	if err != nil {
		return nil, err
	}

	return protocol.NewUpdateCompleteResponsesWith(n)
}

// Delete handles a DELETE query.
func (store *MemStore) Delete(conn postgresql.Conn, q query.Delete) (protocol.Responses, error) {
	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	n, err := tbl.Delete(q.Where())
	if err != nil {
		return nil, err
	}

	return protocol.NewDeleteCompleteResponsesWith(n)
}

// SystemSelect handles a SELECT query for system tables.
func (store *MemStore) SystemSelect(conn postgresql.Conn, q query.Select) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 8.0: System Catalogs
	// https://www.postgresql.org/docs/8.0/catalogs.html
	// PostgreSQL: Documentation: 16: Part IV. Client Interfaces
	// https://www.postgresql.org/docs/current/client-interfaces.html

	selectInformationSchemaColumns := func(conn postgresql.Conn, q query.Select) (protocol.Responses, error) {
		return nil, query.NewErrNotImplemented("SELECT")
	}

	from := q.From()
	if len(from) != 1 {
		return nil, query.NewErrMultipleTableNotSupported(from.String())
	}

	tbl := from[0]
	switch tbl.Name() {
	case system.InformationSchemaColumns:
		return selectInformationSchemaColumns(conn, q)
	}

	return nil, query.NewErrNotSupported(tbl.Name())
}

// Copy handles a COPY query.
func (store *MemStore) Copy(conn postgresql.Conn, q query.Copy) (protocol.Responses, error) {
	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	return postgresql.NewCopyInResponsesFrom(q, tbl.Schema)
}

// Copy handles a COPY DATA protocol.
func (store *MemStore) CopyData(conn postgresql.Conn, q query.Copy, stream *postgresql.CopyStream) (protocol.Responses, error) {
	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return postgresql.NewCopyCompleteResponsesFrom(q, stream, conn, tbl.Schema, store)
}
