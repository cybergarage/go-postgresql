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
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-postgresql/postgresql/query"
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
	err := server.open()
	if err != nil {
		return err
	}

	go server.serve()

	addr := net.JoinHostPort(server.addr, strconv.Itoa(server.port))
	log.Infof("%s/%s (%s) started", PackageName, Version, addr)

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	err := server.close()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(server.addr, strconv.Itoa(server.port))
	log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)

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
	addr := net.JoinHostPort(server.addr, strconv.Itoa(server.port))
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
func (server *Server) receive(netConn net.Conn) error { //nolint:gocyclo,maintidx
	defer netConn.Close()

	log.Debugf("%s/%s (%s) accepted", PackageName, Version, netConn.RemoteAddr().String())

	responseMessage := func(resMsg message.Response) error {
		if resMsg == nil {
			return nil
		}
		resBytes, err := resMsg.Bytes()
		if err != nil {
			return err
		}
		if _, err := netConn.Write(resBytes); err != nil {
			return err
		}
		return nil
	}

	responseMessages := func(resMsgs message.Responses) error {
		if resMsgs == nil {
			return nil
		}
		for _, resMsg := range resMsgs {
			err := responseMessage(resMsg)
			if err != nil {
				return err
			}
		}
		return nil
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
		err = exConn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ParameterStatus (S) message.
		reses, err := server.Executor.ParameterStatuses(exConn)
		if err != nil {
			return err
		}
		err = exConn.ResponseMessages(reses)
		if err != nil {
			return err
		}
		// Return BackendKeyData (K) message.
		res, err = server.Executor.BackendKeyData(exConn)
		if err != nil {
			return err
		}
		err = exConn.ResponseMessage(res)
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

	handleParseMessage := func(conn *Conn) error {
		parseMsg, err := message.NewParseWithReader(conn.MessageReader)
		if err != nil {
			return err
		}

		res, err := server.Executor.Parse(conn, parseMsg)
		if err != nil {
			return err
		}

		err = conn.ResponseMessages(res)
		if err != nil {
			return err
		}

		return nil
	}

	handleBindMessage := func(conn *Conn) error {
		bindMsg, err := message.NewBindWithReader(conn.MessageReader)
		if err != nil {
			return err
		}

		res, err := server.Executor.Bind(conn, bindMsg)
		if err != nil {
			return err
		}

		err = conn.ResponseMessages(res)
		if err != nil {
			return err
		}

		return nil
	}

	handleDescMessage := func(conn *Conn) error {
		descMsg, err := message.NewDescribeWithReader(conn.MessageReader)
		if err != nil {
			return err
		}

		res, err := server.Executor.Describe(conn, descMsg)
		if err != nil {
			return err
		}

		err = responseMessages(res)
		if err != nil {
			return err
		}

		return nil
	}

	handleExecuteMessage := func(conn *Conn) error {
		execMsg, err := message.NewExecuteWithReader(conn.MessageReader)
		if err != nil {
			return err
		}

		res, err := server.Executor.Execute(conn, execMsg)
		if err != nil {
			return err
		}

		err = responseMessages(res)
		if err != nil {
			return err
		}

		return nil
	}

	handleCopyQuery := func(conn *Conn, stmt *query.Copy) (message.Responses, error) {
		res, err := server.Executor.Copy(conn, stmt)
		if err != nil || res.HasErrorResponse() {
			return res, err
		}
		err = responseMessages(res)
		if err != nil {
			return nil, err
		}

		ok, err := conn.IsPeekType(message.CopyDataMessage)
		if !ok || err != nil {
			return nil, err
		}

		return server.Executor.CopyData(conn, stmt, NewCopyStreamWithReader(conn.MessageReader))
	}

	// executeQuery executes a query and returns the result.
	executeQuery := func(conn *Conn, queryMsg *message.Query) error {
		q := queryMsg.Query
		log.Debugf("%s %s", conn.conn.RemoteAddr(), q)

		parser := query.NewParser()
		stmts, err := parser.ParseString(q)
		if err != nil {
			res, err := server.Executor.ParserError(conn, q, err)
			if err != nil {
				return err
			}
			return responseMessages(res)
		}

		for _, stmt := range stmts {
			var err error
			err = stmt.Bind(queryMsg.BindParams)
			if err != nil {
				return err
			}

			var res message.Responses
			switch stmt := stmt.Statement.(type) {
			case *query.Begin:
				res, err = server.Executor.Begin(conn, stmt)
			case *query.Commit:
				res, err = server.Executor.Commit(conn, stmt)
			case *query.Rollback:
				res, err = server.Executor.Rollback(conn, stmt)
			case *query.CreateDatabase:
				res, err = server.Executor.CreateDatabase(conn, stmt)
			case *query.CreateTable:
				res, err = server.Executor.CreateTable(conn, stmt)
			case *query.AlterDatabase:
				res, err = server.Executor.AlterDatabase(conn, stmt)
			case *query.AlterTable:
				res, err = server.Executor.AlterTable(conn, stmt)
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
			case *query.Truncate:
				res, err = server.Executor.Truncate(conn, stmt)
			case *query.Vacuum:
				res, err = server.Executor.Vacuum(conn, stmt)
			case *query.Copy:
				res, err = handleCopyQuery(conn, stmt)
			}

			if err != nil {
				return err
			}

			err = responseMessages(res)
			if err != nil {
				return err
			}
		}
		return nil
	}

	skipMessage := func(reader *message.MessageReader) error {
		msg, err := message.NewMessageWithReader(reader)
		if err != nil {
			return err
		}
		_, err = msg.ReadMessageData()
		if err != nil {
			return err
		}
		return nil
	}

	conn := NewConnWith(netConn)

	// Checks the SSLRequest message.

	startupMsgLength, err := conn.PeekInt32()
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	if startupMsgLength == 8 {
		_, err := message.NewSSLRequestWithReader(conn.MessageReader)
		if err != nil {
			conn.ResponseError(err)
			return err
		}
		err = responseMessage(message.NewSSLResponseWith(message.SSLDisabled))
		if err != nil {
			return err
		}
	}

	// Handle a Start-up message.

	startupMsg, err := message.NewStartupWithReader(conn.MessageReader)
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	err = handleStartupMessage(netConn, startupMsg)
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	// Handle the request messages after the Start-up message.

	dbname := ""
	if db, ok := startupMsg.Database(); ok {
		dbname = db
	}

	conn.SetDatabase(dbname)

	for {
		loopSpan := server.Tracer.StartSpan(PackageName)
		loopSpan.StartSpan("parse")
		conn.SetSpanContext(loopSpan)

		var reqErr error
		var reqType message.Type
		reqType, reqErr = conn.PeekType()
		if reqErr != nil {
			conn.ResponseError(reqErr)
			loopSpan.FinishSpan()
			break
		}

		switch reqType { // nolint:exhaustive
		case message.ParseMessage:
			reqErr = handleParseMessage(conn)
		case message.BindMessage:
			var queryMsg *message.Query
			reqErr = handleBindMessage(conn)
			if reqErr == nil {
				reqErr = executeQuery(conn, queryMsg)
			}
		case message.DescribeMessage:
			reqErr = handleDescMessage(conn)
		case message.QueryMessage:
			var queryMsg *message.Query
			queryMsg, reqErr = message.NewQueryWithReader(conn.MessageReader)
			if reqErr == nil {
				reqErr = executeQuery(conn, queryMsg)
			}
		case message.ExecuteMessage:
			reqErr = handleExecuteMessage(conn)
		case message.SyncMessage:
			// Ignore the Sync message.
			_, reqErr = message.NewTerminateWithReader(conn.MessageReader)
		case message.TerminateMessage:
			_, reqErr := message.NewSyncWithReader(conn.MessageReader)
			if reqErr == nil {
				return nil
			}
		default:
			reqErr = skipMessage(conn.MessageReader)
			if reqErr == nil {
				reqErr = message.NewErrMessageNotSuppoted(reqType)
				log.Warnf(reqErr.Error())
			}
		}

		if reqErr != nil {
			conn.ResponseError(reqErr)
			log.Error(reqErr)
		}

		// Return ReadyForQuery (B) message.
		err := readyForMessage(conn, message.TransactionIdle)

		loopSpan.FinishSpan()

		if err != nil {
			return err
		}
	}

	return nil
}
