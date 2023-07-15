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

package postgresql

import (
	"bufio"
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-sqlparser/sql"
	"github.com/cybergarage/go-sqlparser/sql/query"
	"github.com/cybergarage/go-tracing/tracer"
)

// Server represents a PostgreSQL protocol server.
type Server struct {
	*Config
	tracer.Tracer
	tcpListener net.Listener
	Executor
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Config:      NewDefaultConfig(),
		Tracer:      tracer.NullTracer,
		tcpListener: nil,
		Executor:    NewBaseExecutor(),
	}
	return server
}

// SetTracer sets a tracing tracer.
func (server *Server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// Start starts the server.
func (server *Server) Start() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	err = server.open()
	if err != nil {
		return err
	}

	go server.serve()

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	err := server.close()
	if err != nil {
		return err
	}

	return nil
}

// Restart restarts the server.
func (server *Server) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	return server.Start()
}

// open opens a listen socket.
func (server *Server) open() error {
	var err error
	addr := net.JoinHostPort(server.host, strconv.Itoa(server.port))
	server.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

// close closes a listening socket.
func (server *Server) close() error {
	if server.tcpListener != nil {
		err := server.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	server.tcpListener = nil

	return nil
}

// serve handles client requests.
func (server *Server) serve() error {
	defer server.close()

	l := server.tcpListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go server.receive(conn)
	}

	return nil
}

// receive handles client messages.
func (server *Server) receive(conn net.Conn) error { //nolint:gocyclo,maintidx
	defer conn.Close()

	responseMessage := func(resMsg message.Response) error {
		resBytes, err := resMsg.Bytes()
		if err != nil {
			return err
		}
		if _, err := conn.Write(resBytes); err != nil {
			return err
		}
		return nil
	}

	responseError := func(err error) {
		errMsg, err := message.NewErrorResponseWith(err)
		if err != nil {
			log.Error(err)
			return
		}
		errBytes, err := errMsg.Bytes()
		if err != nil {
			log.Error(err)
			return
		}
		if _, err := conn.Write(errBytes); err != nil {
			log.Error(err)
		}
	}

	readyForMessage := func(exConn *Conn, status message.TransactionStatus) error {
		readyMsg, err := message.NewReadyForQueryWith(status)
		if err != nil {
			return err
		}
		err = responseMessage(readyMsg)
		if err != nil {
			return err
		}
		return nil
	}

	handleStartupMessage := func(conn net.Conn, startupMsg *message.Startup) error {
		exConn := NewConnWith(conn)

		// PostgreSQL: Documentation: 16: 55.2. Message Flow
		// https://www.postgresql.org/docs/16/protocol-flow.html
		// Handle the Start-up message and return an Authentication message or error message.
		res, err := server.Executor.Authenticate(exConn, startupMsg)
		if err != nil {
			return err
		}
		err = responseMessage(res)
		if err != nil {
			return err
		}
		// Return ParameterStatus (S) message.
		res, err = server.Executor.ParameterStatus(exConn)
		if err != nil {
			return err
		}
		err = responseMessage(res)
		if err != nil {
			return err
		}
		// Return BackendKeyData (K) message.
		res, err = server.Executor.BackendKeyData(exConn)
		if err != nil {
			return err
		}
		err = responseMessage(res)
		if err != nil {
			return err
		}
		// Return ReadyForQuery (B) message.
		err = readyForMessage(exConn, message.TransactionIdle)
		if err != nil {
			return err
		}
		return nil
	}

	handleParseBindMessage := func(conn *Conn, reqMsg *message.RequestMessage) (*message.Query, error) {
		parseMsg, err := reqMsg.ParseParseMessage()
		if err != nil {
			return nil, err
		}

		resMsg, err := server.Executor.Parse(conn, parseMsg)
		if err != nil {
			return nil, err
		}

		err = responseMessage(resMsg)
		if err != nil {
			return nil, err
		}

		reqType, err := reqMsg.ReadType()
		if err != nil {
			return nil, err
		}

		if reqType != message.BindMessage {
			return nil, message.NewErrMessageNotSuppoted(reqType)
		}

		bindMsg, err := reqMsg.ParseBindMessageWith(parseMsg)
		if err != nil {
			return nil, err
		}

		resMsg, err = server.Executor.Bind(conn, parseMsg, bindMsg)
		if err != nil {
			return nil, err
		}

		q, err := message.NewQueryWith(parseMsg, bindMsg)
		if err != nil {
			return nil, err
		}

		err = responseMessage(resMsg)
		if err != nil {
			return nil, err
		}

		return q, nil
	}

	// executeQuery executes a query and returns the result.
	executeQuery := func(conn *Conn, queryMsg *message.Query) error {
		parser := sql.NewParser()
		stmts, err := parser.ParseString(queryMsg.Query)
		if err != nil {
			return err
		}

		bindStmt := func(stmt query.Statement, params message.BindParams) error {
			updateBindColumns := func(columns []*query.Column, params message.BindParams) error {
				for _, column := range columns {
					v, ok := column.Value().(*query.BindParam)
					if !ok {
						continue
					}
					param, err := params.FindBindParam(v.Name())
					if err != nil {
						return err
					}
					column.SetValue(param.Value)
				}
				return nil
			}

			switch stmt := stmt.(type) {
			case *query.Insert:
				err := updateBindColumns(stmt.Columns(), params)
				if err != nil {
					return err
				}
			case *query.Update:
				err := updateBindColumns(stmt.Columns(), params)
				if err != nil {
					return err
				}
			}
			return nil
		}

		for _, stmt := range stmts {
			var err error
			err = bindStmt(stmt, queryMsg.BindParams)
			if err != nil {
				return err
			}

			var res message.Responses
			switch stmt := stmt.(type) {
			case *query.CreateDatabase:
				res, err = server.Executor.CreateDatabase(conn, stmt)
			case *query.CreateTable:
				res, err = server.Executor.CreateTable(conn, stmt)
			case *query.CreateIndex:
				res, err = server.Executor.CreateIndex(conn, stmt)
			case *query.DropDatabase:
				res, err = server.Executor.DropDatabase(conn, stmt)
			case *query.DropTable:
				res, err = server.Executor.DropTable(conn, stmt)
			case *query.Insert:
				res, err = server.Executor.Insert(conn, stmt)
			case *query.Select:
				res, err = server.Executor.Select(conn, stmt)
			case *query.Update:
				res, err = server.Executor.Update(conn, stmt)
			case *query.Delete:
				res, err = server.Executor.Delete(conn, stmt)
			}

			for _, r := range res {
				err := responseMessage(r)
				if err != nil {
					return err
				}
			}

			if err != nil {
				return err
			}
		}
		return nil
	}

	// Handle a Start-up message.

	reqMsg := message.NewRequestMessageWith(bufio.NewReader(conn))

	startupMsg, err := reqMsg.ParseStartupMessage()
	if err != nil {
		responseError(err)
		return err
	}

	err = handleStartupMessage(conn, startupMsg)
	if err != nil {
		responseError(err)
		return err
	}

	// Handle the request messages after the Start-up message.

	dbname := ""
	if db, ok := startupMsg.Database(); ok {
		dbname = db
	}

	for {
		loopSpan := server.Tracer.StartSpan(PackageName)
		exConn := NewConnWith(conn, WithConnTracer(loopSpan), WithConnDatabase(dbname))
		loopSpan.StartSpan("parse")

		reqMsg := message.NewRequestMessageWith(bufio.NewReader(conn))

		var reqErr error
		var reqType message.Type
		reqType, reqErr = reqMsg.ReadType()
		if reqErr != nil {
			responseError(reqErr)
			loopSpan.FinishSpan()
			break
		}

		switch reqType { // nolint:exhaustive
		case message.ParseMessage:
			var queryMsg *message.Query
			queryMsg, reqErr = handleParseBindMessage(exConn, reqMsg)
			if reqErr == nil {
				reqErr = executeQuery(exConn, queryMsg)
			}
		case message.QueryMessage:
			var queryMsg *message.Query
			queryMsg, reqErr = reqMsg.ParseQueryMessage()
			if reqErr == nil {
				reqErr = executeQuery(exConn, queryMsg)
			}
		case message.TerminateMessage:
			return nil
		default:
			reqErr = message.NewErrMessageNotSuppoted(reqType)
			log.Warnf(reqErr.Error())
		}

		if reqErr != nil {
			// Return ErrorResponse (B)
			responseError(reqErr)
		}

		// Return ReadyForQuery (B) message.
		err := readyForMessage(exConn, message.TransactionIdle)

		loopSpan.FinishSpan()

		if err != nil {
			return err
		}
	}

	return nil
}
