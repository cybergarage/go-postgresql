// Copyright (C) 2020 The go-postgresql Authors. All rights reserved.
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

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-sqlparser/sql"
	"github.com/cybergarage/go-sqlparser/sql/errors"
	"github.com/cybergarage/go-sqlparser/sql/fn"
	"github.com/cybergarage/go-sqlparser/sql/net"
	"github.com/cybergarage/go-sqlparser/sql/query"
	"github.com/cybergarage/go-sqlparser/sql/query/response/resultset"
	"github.com/cybergarage/go-sqlparser/sql/system"
)

// Store represents a data store.
type Store struct {
	Databases
}

// NewStore returns a new store instance.
func NewStore() *Store {
	store := &Store{
		Databases: NewDatabases(),
	}
	return store
}

func (store *Store) LookupDatabaseTable(conn net.Conn, dbName string, tblName string) (*Database, *Table, error) {
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, nil, errors.NewErrDatabaseNotExist(dbName)
	}

	tbl, ok := db.LookupTable(tblName)
	if !ok {
		return nil, nil, errors.NewErrTableNotExist(tblName)
	}

	return db, tbl, nil
}

// Begin should handle a BEGIN statement.
func (store *Store) Begin(conn net.Conn, stmt query.Begin) error {
	log.Debugf("%v", stmt)
	return nil
}

// Commit should handle a COMMIT statement.
func (store *Store) Commit(conn net.Conn, stmt query.Commit) error {
	log.Debugf("%v", stmt)
	return nil
}

// Rollback should handle a ROLLBACK statement.
func (store *Store) Rollback(conn net.Conn, stmt query.Rollback) error {
	log.Debugf("%v", stmt)
	return nil
}

// Use should handle a USE statement.
func (store *Store) Use(conn net.Conn, stmt query.Use) error {
	log.Debugf("%v", stmt)
	conn.SetDatabase(stmt.DatabaseName())
	return nil
}

// CreateDatabase should handle a CREATE database statement.
func (store *Store) CreateDatabase(conn net.Conn, stmt query.CreateDatabase) error {
	log.Debugf("%v", stmt)

	dbName := stmt.DatabaseName()
	_, ok := store.LookupDatabase(dbName)
	if ok {
		if stmt.IfNotExists() {
			return nil
		}
		return errors.NewErrDatabaseExist(dbName)
	}

	return store.AddDatabase(NewDatabaseWithName(dbName))
}

// AlterDatabase should handle a ALTER database statement.
func (store *Store) AlterDatabase(conn net.Conn, stmt query.AlterDatabase) error {
	log.Debugf("%v", stmt)
	return errors.ErrNotImplemented
}

// DropDatabase should handle a DROP database statement.
func (store *Store) DropDatabase(conn net.Conn, stmt query.DropDatabase) error {
	log.Debugf("%v", stmt)

	dbName := stmt.DatabaseName()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		if stmt.IfExists() {
			return nil
		}
		return errors.NewErrDatabaseNotExist(dbName)
	}
	return store.Databases.DropDatabase(db)
}

// CreateTable should handle a CREATE table statement.
func (store *Store) CreateTable(conn net.Conn, stmt query.CreateTable) error {
	log.Debugf("%v", stmt)

	dbName := conn.Database()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return errors.NewErrDatabaseNotExist(dbName)
	}
	tableName := stmt.TableName()
	_, ok = db.LookupTable(tableName)
	if !ok {
		table := NewTableWith(tableName, stmt.Schema())
		db.AddTable(table)
	} else {
		if !stmt.IfNotExists() {
			return errors.NewErrTableExist(tableName)
		}
	}
	return nil
}

// AlterTable should handle a ALTER table statement.
func (store *Store) AlterTable(conn net.Conn, stmt query.AlterTable) error {
	log.Debugf("%v", stmt)

	dbName := conn.Database()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return errors.NewErrDatabaseNotExist(dbName)
	}
	tableName := stmt.TableName()
	tbl, ok := db.LookupTable(tableName)
	if !ok {
		return errors.NewErrTableExist(tableName)
	}
	return tbl.Schema.Alter(stmt)
}

// DropTable should handle a DROP table statement.
func (store *Store) DropTable(conn net.Conn, stmt query.DropTable) error {
	log.Debugf("%v", stmt)

	dbName := conn.Database()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return errors.NewErrDatabaseNotExist(dbName)
	}
	for _, table := range stmt.Tables() {
		tableName := table.TableName()
		table, ok := db.LookupTable(tableName)
		if !ok {
			if stmt.IfExists() {
				continue
			}
			return errors.NewErrTableNotExist(tableName)
		}

		if !db.DropTable(table) {
			return fmt.Errorf("%s could not deleted", table.TableName())
		}
	}
	return nil
}

// Insert should handle a INSERT statement.
func (store *Store) Insert(conn net.Conn, stmt query.Insert) error {
	log.Debugf("%v", stmt)

	dbName := conn.Database()
	tableName := stmt.TableName()
	table, ok := store.LookupTableWithDatabase(dbName, tableName)
	if !ok {
		return errors.NewErrTableNotExist(tableName)
	}

	table.Lock()
	defer table.Unlock()

	for _, value := range stmt.Values() {
		row, err := NewRowFromColumns(table, value.Columns())
		if err != nil {
			return err
		}
		table.Rows = append(table.Rows, row)
	}

	return nil
}

// Update should handle a UPDATE statement.
func (store *Store) Update(conn net.Conn, stmt query.Update) (sql.ResultSet, error) {
	log.Debugf("%v", stmt)

	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), stmt.TableName())
	if err != nil {
		return nil, err
	}

	n, err := tbl.Update(stmt.Columns(), stmt.Where())
	if err != nil {
		return nil, err
	}

	return resultset.NewResultSet(
		resultset.WithResultSetRowsAffected(uint64(n)),
	), nil
}

// Delete should handle a DELETE statement.
func (store *Store) Delete(conn net.Conn, stmt query.Delete) (sql.ResultSet, error) {
	log.Debugf("%v", stmt)

	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), stmt.TableName())
	if err != nil {
		return nil, err
	}

	n, err := tbl.Delete(stmt.Where())
	if err != nil {
		return nil, err
	}

	return resultset.NewResultSet(
		resultset.WithResultSetRowsAffected(uint64(n)),
	), nil
}

// Select should handle a SELECT statement.
func (store *Store) Select(conn net.Conn, stmt query.Select) (sql.ResultSet, error) {
	log.Debugf("%v", stmt)

	// Select the target table

	from := stmt.From()
	if len(from) != 1 {
		return nil, errors.NewErrMultipleTableNotSupported(from.String())
	}

	tblName := from[0].TableName()

	_, tbl, err := store.LookupDatabaseTable(conn, conn.Database(), tblName)
	if err != nil {
		return nil, err
	}

	// Selectors

	selectors := stmt.Selectors()
	if selectors.IsAsterisk() {
		selectors = tbl.Selectors()
	}

	// Select rows from a target table

	rows, err := tbl.Select(stmt.Where())
	if err != nil {
		return nil, err
	}

	// Aggregate

	if stmt.HasAggregator() {
		aggrSet, err := selectors.Aggregators()
		if err != nil {
			return nil, err
		}

		resetOpts := []any{}
		if stmt.GroupBy() != nil {
			resetOpts = append(resetOpts, fn.GroupBy(stmt.GroupBy().ColumnName()))
		}

		err = aggrSet.Reset(resetOpts...)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			err := aggrSet.Aggregate(row)
			if err != nil {
				return nil, err
			}
		}

		resultSet, err := aggrSet.Finalize()
		if err != nil {
			return nil, err
		}

		rows = []Row{}
		for resultSet.Next() {
			rowMap, err := resultSet.Map()
			if err != nil {
				return nil, err
			}
			rows = append(rows, NewRowWithResultMap(rowMap))
		}
	}

	// Row description response

	schema := tbl.Schema
	rsSchemaColums := []sql.ResultSetColumn{}
	for _, selector := range selectors {
		var rsCchemaColumn resultset.Column
		fx, ok := selector.Function()
		if !ok {
			selectorName := selector.Name()
			shemaColumn, err := schema.LookupColumn(selectorName)
			if err != nil {
				return nil, err
			}
			rsCchemaColumn, err = resultset.NewColumnFrom(shemaColumn)
			if err != nil {
				return nil, err
			}
		} else {
			dataType, err := query.NewDataTypeForFunction(fx)
			if err != nil {
				return nil, err
			}
			rsCchemaColumn = resultset.NewColumn(
				resultset.WithColumnName(selector.String()),
				resultset.WithColumnType(dataType),
				resultset.WithColumnFunction(fx),
			)

		}
		rsSchemaColums = append(rsSchemaColums, rsCchemaColumn)
	}

	rsSchema := resultset.NewSchema(
		resultset.WithSchemaDatabaseName(conn.Database()),
		resultset.WithSchemaTableName(tblName),
		resultset.WithSchemaColumns(rsSchemaColums),
	)

	// offset and limit

	offset := stmt.Limit().Offset()
	if 0 < offset && len(rows) <= offset {
		rows = rows[offset:]
	}

	limit := stmt.Limit().Limit()
	if 0 < limit && limit < len(rows) {
		rows = rows[:limit]
	}

	// Data row response

	rsRows := []sql.ResultSetRow{}
	for _, row := range rows {
		rowValues := []any{}
		for _, selector := range selectors {
			var rowValue any
			rowValue = nil
			if fx, ok := selector.Function(); ok {
				if executor, err := fx.Executor(); err == nil {
					rowValue, err = executor.Execute(fn.NewMapWithMap(row))
					if err != nil {
						return nil, err
					}
				}
			}
			if rowValue == nil {
				selectorName := selector.Name()
				rowValue, err = row.ValueByName(selectorName)
				if err != nil {
					return nil, err
				}
			}
			rowValues = append(rowValues, rowValue)
		}
		rsRow := resultset.NewRow(
			resultset.WithRowSchema(rsSchema),
			resultset.WithRowValues(rowValues),
		)
		rsRows = append(rsRows, rsRow)
	}

	// Return a result set

	rs := resultset.NewResultSet(
		resultset.WithResultSetSchema(rsSchema),
		resultset.WithResultSetRowsAffected(uint64(len(rsRows))),
		resultset.WithResultSetRows(rsRows),
	)

	return rs, nil
}

// SystemSelect should handle a system SELECT statement.
func (store *Store) SystemSelect(conn net.Conn, stmt query.Select) (sql.ResultSet, error) {
	q := stmt.String()
	log.Debugf("%v", q)

	switch {
	case system.IsSchemaColumsQuery(stmt):
		sysStmt, err := system.NewSchemaColumnsStatement(
			system.WithSchemaColumnsStatementSelect(stmt),
		)
		if err != nil {
			return nil, err
		}
		dbName := sysStmt.DatabaseName()
		tblNames := sysStmt.TableNames()
		schemas := []query.Schema{}
		for _, tblName := range tblNames {
			_, tbl, err := store.LookupDatabaseTable(conn, dbName, tblName)
			if err != nil {
				return nil, err
			}
			schemas = append(schemas, tbl.Schema)
		}
		return system.NewSchemaColumnsResultSetFromSchemas(schemas)
	}

	return nil, errors.NewErrNotImplemented(fmt.Sprintf("SystemSelect: %s", stmt.String()))
}
