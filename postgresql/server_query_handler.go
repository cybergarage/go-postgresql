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
	stderrors "errors"
	"strings"

	"github.com/cybergarage/go-postgresql/postgresql/errors"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
	"github.com/cybergarage/go-postgresql/postgresql/stmt"
	"github.com/cybergarage/go-postgresql/postgresql/system"
	"github.com/cybergarage/go-safecast/safecast"
	sqlparser "github.com/cybergarage/go-sqlparser/sql/parser"
	sql "github.com/cybergarage/go-sqlparser/sql/query"
	sql_system "github.com/cybergarage/go-sqlparser/sql/system"
)

// protocolQueryHandler represents a protocol query server.
type protocolQueryHandler struct {
	*stmt.PreparedManager
}

// newProtocolQueryHandlerWith returns a new protocol query server.
func newProtocolQueryHandler() *protocolQueryHandler {
	return &protocolQueryHandler{
		PreparedManager: stmt.NewPreparedManager(),
	}
}

// Prepare handles a parse protocol.
func (server *server) Parse(conn Conn, msg *protocol.Parse) (protocol.Responses, error) {
	err := server.SetPreparedStatement(conn, msg)
	if err != nil {
		return nil, err
	}
	return protocol.NewResponsesWith(protocol.NewParseComplete()), nil
}

// Bind handles a bind protocol.
func (server *server) Bind(conn Conn, msg *protocol.Bind) (protocol.Responses, error) {
	prepStmt, err := server.PreparedStatement(conn, msg.StatementName)
	if err != nil {
		return nil, err
	}

	q, err := protocol.NewQueryWith(prepStmt.Parse, msg)
	if err != nil {
		return nil, err
	}

	err = server.SetPreparedPortal(conn, msg.PortalName, q)
	if err != nil {
		return nil, err
	}

	return protocol.NewResponsesWith(protocol.NewBindComplete()), nil
}

// Describe handles a describe protocol.
func (server *server) Describe(conn Conn, msg *protocol.Describe) (protocol.Responses, error) {
	newSystemSelectQuery := func(stmt query.Select) (sql.Select, error) {
		tables := stmt.From().Tables()
		if len(tables) != 1 {
			return nil, errors.NewErrMultipleTableNotSupported(stmt.From().String())
		}
		table := tables[0]
		sysStmt, err := sql_system.NewSchemaColumnsStatement(
			sql_system.WithSchemaColumnsStatementDatabaseName(conn.Database()),
			sql_system.WithSchemaColumnsStatementTableNames([]string{table.TableName()}),
		)
		return sysStmt.Statement(), err
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
		res, err := server.systemQueryExecutor.SystemSelect(conn, query)
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

	switch msg.PreparedType() {
	case protocol.PreparedStatement:
		prepStmt, err := server.PreparedStatement(conn, msg.Name())
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
		_, err := server.PreparedPortal(conn, msg.Name())
		if err != nil {
			return nil, err
		}
		return protocol.NewResponsesWith(
			protocol.NewNoData()), nil
	}
	return nil, nil
}

// Execute handles a execute protocol.
func (server *server) Execute(conn Conn, msg *protocol.Execute) (protocol.Responses, error) {
	q, err := server.PreparedPortal(conn, msg.PortalName)
	if err != nil {
		return nil, err
	}

	return server.Query(conn, q)
}

// Close handles a close protocol.
func (server *server) Close(conn Conn, msg *protocol.Close) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// https://www.postgresql.org/docs/16/protocol-flow.html
	// The Close message closes an existing prepared statement or portal and releases resources.
	// It is not an error to issue Close against a nonexistent statement or portal name.

	switch msg.Type {
	case protocol.PreparedStatement:
		_ = server.RemovePreparedStatement(conn, msg.Name)
	case protocol.PreparedPortal:
		_ = server.RemovePreparedPortal(conn, msg.Name)
	}

	return protocol.NewResponsesWith(protocol.NewCloseComplete()), nil
}

// Sync handles a sync protocol.
func (server *server) Sync(conn Conn, msg *protocol.Sync) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// https://www.postgresql.org/docs/16/protocol-flow.html
	// At completion of each series of extended-query messages, the frontend should issue a Sync protocol.
	return nil, nil
}

// Flush handles a flush protocol.
func (server *server) Flush(conn Conn, msg *protocol.Flush) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.2. Message Flow
	// https://www.postgresql.org/docs/16/protocol-flow.html
	// The Flush message does not cause any specific output to be generated,
	// but forces the backend to deliver any data pending in its output buffers.
	return nil, nil
}

// Query handles a query protocol.
func (server *server) Query(conn Conn, msg *protocol.Query) (protocol.Responses, error) {
	conn.StartSpan("parse")
	stmts, err := msg.Statements()
	conn.FinishSpan()
	if err != nil {
		// Is it a empty query for ping?
		if stderrors.Is(err, sqlparser.ErrEmptyQuery) {
			return protocol.NewEmptyCompleteResponses()
		}
		res, err := server.errorHandler.ParserError(conn, msg.String(), err)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	handleCopyQuery := func(conn Conn, stmt query.Copy) (protocol.Responses, error) {
		res, err := server.bulkQueryExecutor.Copy(conn, stmt)
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

		return server.bulkQueryExecutor.CopyData(conn, stmt, NewCopyStreamWithReader(conn.MessageReader()))
	}

	for _, stmt := range stmts {
		var res protocol.Responses

		// nolint: forcetypeassert
		switch stmt.StatementType() {
		case sql.BeginStatement:
			err = conn.LockTransaction()
			if err == nil {
				stmt := stmt.(query.Begin)
				res, err = server.queryExecutor.Begin(conn, stmt)
			}
		case sql.CommitStatement:
			stmt := stmt.(query.Commit)
			res, err = server.queryExecutor.Commit(conn, stmt)
			conn.UnlockTransaction()
		case sql.RollbackStatement:
			stmt := stmt.(query.Rollback)
			res, err = server.queryExecutor.Rollback(conn, stmt)
			conn.UnlockTransaction()
		case sql.CreateDatabaseStatement:
			stmt := stmt.(query.CreateDatabase)
			res, err = server.queryExecutor.CreateDatabase(conn, stmt)
		case sql.CreateTableStatement:
			stmt := stmt.(query.CreateTable)
			res, err = server.queryExecutor.CreateTable(conn, stmt)
		case sql.CreateIndexStatement:
			stmt := stmt.(query.CreateIndex)
			res, err = server.exQueryExecutor.CreateIndex(conn, stmt)
		case sql.AlterDatabaseStatement:
			stmt := stmt.(query.AlterDatabase)
			res, err = server.queryExecutor.AlterDatabase(conn, stmt)
		case sql.AlterTableStatement:
			stmt := stmt.(query.AlterTable)
			res, err = server.queryExecutor.AlterTable(conn, stmt)
		case sql.DropDatabaseStatement:
			stmt := stmt.(query.DropDatabase)
			res, err = server.queryExecutor.DropDatabase(conn, stmt)
		case sql.DropTableStatement:
			stmt := stmt.(query.DropTable)
			res, err = server.queryExecutor.DropTable(conn, stmt)
		case sql.DropIndexStatement:
			stmt := stmt.(query.DropIndex)
			res, err = server.exQueryExecutor.DropIndex(conn, stmt)
		case sql.InsertStatement:
			stmt := stmt.(query.Insert)
			res, err = server.queryExecutor.Insert(conn, stmt)
		case sql.SelectStatement:
			stmt := stmt.(query.Select)
			if stmt.From().HasSchemaTable(system.SystemSchemaNames...) {
				res, err = server.systemQueryExecutor.SystemSelect(conn, stmt)
			} else {
				res, err = server.queryExecutor.Select(conn, stmt)
			}
		case sql.UpdateStatement:
			stmt := stmt.(query.Update)
			res, err = server.queryExecutor.Update(conn, stmt)
		case sql.DeleteStatement:
			stmt := stmt.(query.Delete)
			res, err = server.queryExecutor.Delete(conn, stmt)
		case sql.TruncateStatement:
			stmt := stmt.(query.Truncate)
			res, err = server.exQueryExecutor.Truncate(conn, stmt)
		case sql.VacuumStatement:
			stmt := stmt.(query.Vacuum)
			res, err = server.exQueryExecutor.Vacuum(conn, stmt)
		case sql.CopyStatement:
			stmt := stmt.(query.Copy)
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
