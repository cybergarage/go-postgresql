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

package postgresql

import (
	"errors"
	"strings"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
	"github.com/cybergarage/go-postgresql/postgresql/system"
	"github.com/cybergarage/go-safecast/safecast"
	sqlparser "github.com/cybergarage/go-sqlparser/sql/parser"
	sql "github.com/cybergarage/go-sqlparser/sql/query"
)

// BaseExtendedQueryExecutor represents a base extended query message executor.
type BaseExtendedQueryExecutor struct {
	*BaseExecutor
}

// NewBaseExtendedQueryExecutorWith returns a base extended query message executor.
func NewBaseExtendedQueryExecutorWith(executor *BaseExecutor) *BaseExtendedQueryExecutor {
	return &BaseExtendedQueryExecutor{
		BaseExecutor: executor,
	}
}

// Prepare handles a parse protocol.
func (executor *BaseExtendedQueryExecutor) Parse(conn Conn, msg *protocol.Parse) (protocol.Responses, error) {
	err := conn.SetPreparedStatement(msg)
	if err != nil {
		return nil, err
	}
	return protocol.NewResponsesWith(protocol.NewParseComplete()), nil
}

// Bind handles a bind protocol.
func (executor *BaseExtendedQueryExecutor) Bind(conn Conn, msg *protocol.Bind) (protocol.Responses, error) {
	prepStmt, err := conn.PreparedStatement(msg.StatementName)
	if err != nil {
		return nil, err
	}

	q, err := protocol.NewQueryWith(prepStmt.Parse, msg)
	if err != nil {
		return nil, err
	}

	err = conn.SetPreparedPortal(msg.PortalName, q)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponsesWith(protocol.NewBindComplete()), nil
}

// Describe handles a describe protocol.
func (executor *BaseExtendedQueryExecutor) Describe(conn Conn, msg *protocol.Describe) (protocol.Responses, error) {
	newSystemSelectQuery := func(stmt query.Select) (*sql.Select, error) {
		tables := stmt.From().Tables()
		if len(tables) != 1 {
			return nil, query.NewErrMultipleTableNotSupported(stmt.From().String())
		}
		table := tables[0]
		return sql.NewSelectWith(
			sql.NewSelectorsWith(
				sql.NewColumnWithName(system.InformationSchemaColumnsColumnName),
				sql.NewColumnWithName(system.InformationSchemaColumnsDataType),
			),
			sql.NewTablesWith(sql.NewTableWith(system.InformationSchemaColumns)),
			sql.NewConditionWith(
				sql.NewEQWith(
					sql.NewColumnWithOptions(sql.WithColumnName(system.InformationSchemaColumnsTableName)),
					sql.NewLiteralWith(table.TableName()),
				),
			),
		), nil
	}

	selectObjectIds := func(stmt query.Select) ([]int32, error) {
		objIDFromResponses := func(responses protocol.Responses, colName string) (int32, bool) {
			for r, res := range responses {
				if r == 0 {
					continue
				}
				dataRow, ok := res.(*protocol.DataRow)
				if !ok {
					continue
				}
				if len(dataRow.Data) < 2 {
					continue
				}
				columnName, ok := dataRow.Data[0].(string)
				if !ok {
					continue
				}
				if !strings.EqualFold(columnName, colName) {
					continue
				}
				var objID int32
				err := safecast.ToInt32(dataRow.Data[1], &objID)
				if err != nil {
					continue
				}
				return objID, true
			}
			return 0, false
		}

		query, err := newSystemSelectQuery(stmt)
		if err != nil {
			return nil, err
		}
		res, err := executor.SystemQueryExecutor.SystemSelect(conn, query)
		if err != nil {
			return nil, err
		}
		sels := stmt.Selectors()
		objIDs := make([]int32, len(sels))
		for n, sel := range sels {
			selName := sel.Name()
			objID, ok := objIDFromResponses(res, selName)
			if !ok {
				objIDs[n] = 0
				continue
			}
			objIDs[n] = objID
		}
		return objIDs, nil
	}

	switch msg.Type {
	case protocol.PreparedStatement:
		prepStmt, err := conn.PreparedStatement(msg.Name)
		if err != nil {
			return nil, err
		}
		objIDs := []int32{}
		switch stmt := prepStmt.ParsedStatement.Object().(type) {
		case query.Select:
			objIDs, err = selectObjectIds(stmt)
			if err != nil {
				return nil, err
			}
		}
		paramDesc, err := protocol.NewParameterDescriptionWith(objIDs...)
		if err != nil {
			return nil, err
		}
		return protocol.NewResponsesWith(
			paramDesc,
			protocol.NewNoData()), nil
	case protocol.PreparedPortal:
		_, err := conn.PreparedPortal(msg.Name)
		if err != nil {
			return nil, err
		}
		return protocol.NewResponsesWith(
			protocol.NewNoData()), nil
	}
	return nil, nil
}

// Execute handles a execute protocol.
func (executor *BaseExtendedQueryExecutor) Execute(conn Conn, msg *protocol.Execute) (protocol.Responses, error) {
	q, err := conn.PreparedPortal(msg.PortalName)
	if err != nil {
		return nil, err
	}

	return executor.Query(conn, q)
}

// Close handles a close protocol.
func (executor *BaseExtendedQueryExecutor) Close(conn Conn, msg *protocol.Close) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// https://www.postgresql.org/docs/16/protocol-flow.html
	// The Close message closes an existing prepared statement or portal and releases resources.
	// It is not an error to issue Close against a nonexistent statement or portal name.

	switch msg.Type {
	case protocol.PreparedStatement:
		_ = conn.RemovePreparedStatement(msg.Name)
	case protocol.PreparedPortal:
		_ = conn.RemovePreparedPortal(msg.Name)
	}

	return protocol.NewResponsesWith(protocol.NewCloseComplete()), nil
}

// Sync handles a sync protocol.
func (executor *BaseExtendedQueryExecutor) Sync(conn Conn, msg *protocol.Sync) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// https://www.postgresql.org/docs/16/protocol-flow.html
	// At completion of each series of extended-query messages, the frontend should issue a Sync protocol.
	return nil, nil
}

// Flush handles a flush protocol.
func (executor *BaseExtendedQueryExecutor) Flush(conn Conn, msg *protocol.Flush) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// https://www.postgresql.org/docs/16/protocol-flow.html
	// The Flush message does not cause any specific output to be generated,
	// but forces the backend to deliver any data pending in its output buffers.
	return nil, nil
}

// Query handles a query protocol.
func (executor *BaseExtendedQueryExecutor) Query(conn Conn, msg *protocol.Query) (protocol.Responses, error) {
	q := msg.Query
	log.Debugf("%s %s", conn.RemoteAddr(), q)

	conn.StartSpan("parse")
	parser := query.NewParser()
	stmts, err := parser.ParseString(q)
	conn.FinishSpan()
	if err != nil {
		// Is it a empty query for ping?
		if errors.Is(err, sqlparser.ErrEmptyQuery) {
			return protocol.NewEmptyCompleteResponses()
		}
		res, err := executor.ErrorHandler.ParserError(conn, q, err)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	handleCopyQuery := func(conn Conn, stmt query.Copy) (protocol.Responses, error) {
		res, err := executor.BulkExecutor.Copy(conn, stmt)
		if err != nil || res.HasErrorResponse() {
			return res, err
		}
		err = conn.ResponseMessages(res)
		if err != nil {
			return nil, err
		}

		ok, err := conn.MessageReader().IsPeekType(protocol.CopyDataMessage)
		if !ok || err != nil {
			return nil, err
		}

		return executor.BulkExecutor.CopyData(conn, stmt, NewCopyStreamWithReader(conn.MessageReader()))
	}

	for _, stmt := range stmts {
		var err error
		err = stmt.Bind(msg.BindParams)
		if err != nil {
			return nil, err
		}

		var res protocol.Responses

		// nolint: forcetypeassert
		switch stmt.Object().StatementType() {
		case sql.BeginStatement:
			stmt := stmt.Object().(query.Begin)
			res, err = executor.TCLExecutor.Begin(conn, stmt)
		case sql.CommitStatement:
			stmt := stmt.Object().(query.Commit)
			res, err = executor.TCLExecutor.Commit(conn, stmt)
		case sql.RollbackStatement:
			stmt := stmt.Object().(query.Rollback)
			res, err = executor.TCLExecutor.Rollback(conn, stmt)
		case sql.CreateDatabaseStatement:
			stmt := stmt.Object().(query.CreateDatabase)
			res, err = executor.QueryExecutor.CreateDatabase(conn, stmt)
		case sql.CreateTableStatement:
			stmt := stmt.Object().(query.CreateTable)
			res, err = executor.QueryExecutor.CreateTable(conn, stmt)
		case sql.AlterDatabaseStatement:
			stmt := stmt.Object().(query.AlterDatabase)
			res, err = executor.QueryExecutor.AlterDatabase(conn, stmt)
		case sql.AlterTableStatement:
			stmt := stmt.Object().(query.AlterTable)
			res, err = executor.QueryExecutor.AlterTable(conn, stmt)
		case sql.DropDatabaseStatement:
			stmt := stmt.Object().(query.DropDatabase)
			res, err = executor.QueryExecutor.DropDatabase(conn, stmt)
		case sql.DropTableStatement:
			stmt := stmt.Object().(query.DropTable)
			res, err = executor.QueryExecutor.DropTable(conn, stmt)
		case sql.InsertStatement:
			stmt := stmt.Object().(query.Insert)
			res, err = executor.QueryExecutor.Insert(conn, stmt)
		case sql.SelectStatement:
			stmt := stmt.Object().(query.Select)
			if stmt.From().HasSchemaTable(system.SystemSchemaNames...) {
				res, err = executor.SystemQueryExecutor.SystemSelect(conn, stmt)
			} else {
				res, err = executor.QueryExecutor.Select(conn, stmt)
			}
		case sql.UpdateStatement:
			stmt := stmt.Object().(query.Update)
			res, err = executor.QueryExecutor.Update(conn, stmt)
		case sql.DeleteStatement:
			stmt := stmt.Object().(query.Delete)
			res, err = executor.QueryExecutor.Delete(conn, stmt)
		case sql.TruncateStatement:
			stmt := stmt.Object().(query.Truncate)
			res, err = executor.QueryExtraExecutor.Truncate(conn, stmt)
		case sql.VacuumStatement:
			stmt := stmt.Object().(query.Vacuum)
			res, err = executor.QueryExtraExecutor.Vacuum(conn, stmt)
		case sql.CopyStatement:
			stmt := stmt.Object().(query.Copy)
			res, err = handleCopyQuery(conn, stmt)
		}

		if 0 < len(res) {
			err = conn.ResponseMessages(res)
			if err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
