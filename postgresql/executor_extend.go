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
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-postgresql/postgresql/query"
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

// Prepare handles a parse message.
func (executor *BaseExtendedQueryExecutor) Parse(conn *Conn, msg *message.Parse) (message.Responses, error) {
	err := conn.SetPreparedStatement(msg)
	if err != nil {
		return nil, err
	}
	return message.NewResponsesWith(message.NewParseComplete()), nil
}

// Bind handles a bind message.
func (executor *BaseExtendedQueryExecutor) Bind(conn *Conn, msg *message.Bind) (message.Responses, error) {
	preparedQuery, err := conn.PreparedStatement(msg.StatementName)
	if err != nil {
		return nil, err
	}

	q, err := message.NewQueryWith(preparedQuery, msg)
	if err != nil {
		return nil, err
	}

	err = conn.SetPreparedPortal(msg.PortalName, q)
	if err != nil {
		return nil, err
	}

	return message.NewResponsesWith(message.NewBindComplete()), nil
}

// Describe handles a describe message.
func (executor *BaseExtendedQueryExecutor) Describe(conn *Conn, msg *message.Describe) (message.Responses, error) {
	switch msg.Type {
	case message.PreparedStatement:
		_, err := conn.PreparedStatement(msg.Name)
		if err != nil {
			return nil, err
		}
	case message.PreparedPortal:
		_, err := conn.PreparedPortal(msg.Name)
		if err != nil {
			return nil, err
		}
	}

	return message.NewResponsesWith(
		message.NewParameterDescription(),
		message.NewNoData()), nil
}

// Execute handles a execute message.
func (executor *BaseExtendedQueryExecutor) Execute(conn *Conn, msg *message.Execute) (message.Responses, error) {
	q, err := conn.PreparedPortal(msg.PortalName)
	if err != nil {
		return nil, err
	}

	return executor.Query(conn, q)
}

// Close handles a close message.
func (executor *BaseExtendedQueryExecutor) Close(conn *Conn, msg *message.Close) (message.Responses, error) {
	return nil, nil
}

// Sync handles a sync message.
func (executor *BaseExtendedQueryExecutor) Sync(conn *Conn, msg *message.Sync) (message.Responses, error) {
	return nil, nil
}

// Flush handles a flush message.
func (executor *BaseExtendedQueryExecutor) Flush(conn *Conn, msg *message.Flush) (message.Responses, error) {
	return nil, nil
}

// Query handles a query message.
func (executor *BaseExtendedQueryExecutor) Query(conn *Conn, msg *message.Query) (message.Responses, error) {
	q := msg.Query
	log.Debugf("%s %s", conn.conn.RemoteAddr(), q)

	parser := query.NewParser()
	stmts, err := parser.ParseString(q)
	if err != nil {
		res, err := executor.ErrorHandler.ParserError(conn, q, err)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	handleCopyQuery := func(conn *Conn, stmt *query.Copy) (message.Responses, error) {
		res, err := executor.BulkExecutor.Copy(conn, stmt)
		if err != nil || res.HasErrorResponse() {
			return res, err
		}
		err = conn.ResponseMessages(res)
		if err != nil {
			return nil, err
		}

		ok, err := conn.IsPeekType(message.CopyDataMessage)
		if !ok || err != nil {
			return nil, err
		}

		return executor.BulkExecutor.CopyData(conn, stmt, NewCopyStreamWithReader(conn.MessageReader))
	}

	for _, stmt := range stmts {
		var err error
		err = stmt.Bind(msg.BindParams)
		if err != nil {
			return nil, err
		}

		var res message.Responses
		switch stmt := stmt.Statement.(type) {
		case *query.Begin:
			res, err = executor.TransactionExecutor.Begin(conn, stmt)
		case *query.Commit:
			res, err = executor.TransactionExecutor.Commit(conn, stmt)
		case *query.Rollback:
			res, err = executor.TransactionExecutor.Rollback(conn, stmt)
		case *query.CreateDatabase:
			res, err = executor.QueryExecutor.CreateDatabase(conn, stmt)
		case *query.CreateTable:
			res, err = executor.QueryExecutor.CreateTable(conn, stmt)
		case *query.AlterDatabase:
			res, err = executor.QueryExecutor.AlterDatabase(conn, stmt)
		case *query.AlterTable:
			res, err = executor.QueryExecutor.AlterTable(conn, stmt)
		case *query.DropDatabase:
			res, err = executor.QueryExecutor.DropDatabase(conn, stmt)
		case *query.DropTable:
			res, err = executor.QueryExecutor.DropTable(conn, stmt)
		case *query.Insert:
			res, err = executor.QueryExecutor.Insert(conn, stmt)
		case *query.Select:
			res, err = executor.QueryExecutor.Select(conn, stmt)
		case *query.Update:
			res, err = executor.QueryExecutor.Update(conn, stmt)
		case *query.Delete:
			res, err = executor.QueryExecutor.Delete(conn, stmt)
		case *query.Truncate:
			res, err = executor.QueryExtraExecutor.Truncate(conn, stmt)
		case *query.Vacuum:
			res, err = executor.QueryExtraExecutor.Vacuum(conn, stmt)
		case *query.Copy:
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
