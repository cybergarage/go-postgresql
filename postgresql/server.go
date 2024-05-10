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
	"crypto/tls"
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
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

	handleStartupMessage := func(conn *Conn, startupMsg *message.Startup) error {
		// PostgreSQL: Documentation: 16: 55.2. Message Flow
		// https://www.postgresql.org/docs/16/protocol-flow.html
		// Handle the Start-up message and return an Authentication message or error message.
		res, err := server.Executor.Authenticate(conn, startupMsg)
		if err != nil {
			return err
		}
		err = conn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ParameterStatus (S) message.
		reses, err := server.Executor.ParameterStatuses(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessages(reses)
		if err != nil {
			return err
		}
		// Return BackendKeyData (K) message.
		res, err = server.Executor.BackendKeyData(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ReadyForQuery (B) message.
		err = conn.ReadyForMessage(message.TransactionIdle)
		if err != nil {
			return err
		}
		return nil
	}

	conn := NewConnWith(netConn)

	// Checks the SSLRequest message.

	startupMsgLength, err := conn.MessageReader().PeekInt32()
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	// PostgreSQL: Documentation: 16: 55.2.10. SSL Session Encryption
	// https://www.postgresql.org/docs/16/protocol-flow.html
	if startupMsgLength == 8 {
		_, err := message.NewSSLRequestWithReader(conn.MessageReader())
		if err != nil {
			conn.ResponseError(err)
			return err
		}
		tlsConfig, err := server.Config.TLSConfig()
		if err != nil {
			conn.ResponseError(err)
			return err
		}
		if tlsConfig == nil {
			err = conn.ResponseMessage(message.NewSSLResponseWith(message.SSLDisabled))
			if err != nil {
				return err
			}
		}
		err = conn.ResponseMessage(message.NewSSLResponseWith(message.SSLEnabled))
		if err != nil {
			return err
		}
		conn = NewConnWith(tls.Client(conn, tlsConfig))
	}

	// Handle a Start-up message.

	startupMsg, err := message.NewStartupWithReader(conn.MessageReader())
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	err = handleStartupMessage(conn, startupMsg)
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
		var reqErr error
		var reqType message.Type
		reqType, reqErr = conn.MessageReader().PeekType()
		if reqErr != nil {
			conn.ResponseError(reqErr)
			break
		}

		loopSpan := server.Tracer.StartSpan(PackageName)
		conn.SetSpanContext(loopSpan)
		conn.StartSpan(reqType.String())

		var resMsgs message.Responses

		switch reqType { // nolint:exhaustive
		case message.ParseMessage:
			var reqMsg *message.Parse
			reqMsg, reqErr = message.NewParseWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Parse(conn, reqMsg)
			}
		case message.BindMessage:
			var reqMsg *message.Bind
			reqMsg, reqErr = message.NewBindWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Bind(conn, reqMsg)
			}
		case message.DescribeMessage:
			var reqMsg *message.Describe
			reqMsg, reqErr = message.NewDescribeWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Describe(conn, reqMsg)
			}
		case message.QueryMessage:
			var reqMsg *message.Query
			reqMsg, reqErr = message.NewQueryWithReader(conn.MessageReader())
			if reqErr == nil {
				resMsgs, reqErr = server.Executor.Query(conn, reqMsg)
			}
		case message.ExecuteMessage:
			var reqMsg *message.Execute
			reqMsg, err := message.NewExecuteWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Execute(conn, reqMsg)
			}
		case message.CloseMessage:
			var reqMsg *message.Close
			reqMsg, err := message.NewCloseWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Close(conn, reqMsg)
			}
		case message.SyncMessage:
			var reqMsg *message.Sync
			reqMsg, err := message.NewSyncWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Sync(conn, reqMsg)
			}
		case message.FlushMessage:
			var reqMsg *message.Flush
			reqMsg, err := message.NewFlushWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.Executor.Flush(conn, reqMsg)
			}
		case message.TerminateMessage:
			_, reqErr = message.NewTerminateWithReader(conn.MessageReader())
			if reqErr == nil {
				conn.FinishSpan()
				loopSpan.Span().Finish()
				return nil
			}
		default:
			reqErr = conn.SkipMessage()
			if reqErr == nil {
				reqErr = message.NewErrMessageNotSuppoted(reqType)
				log.Warnf(reqErr.Error())
			}
		}

		conn.FinishSpan()

		conn.StartSpan("response")
		err = conn.ResponseMessages(resMsgs)
		conn.FinishSpan()
		if err != nil {
			loopSpan.Span().Finish()
			return err
		}

		if reqErr != nil {
			err := conn.ResponseError(reqErr)
			if err != nil {
				loopSpan.Span().Finish()
				return err
			}
		}

		// Return ReadyForQuery (B) message.
		conn.StartSpan("ready")
		err := conn.ReadyForMessage(message.TransactionIdle)
		conn.FinishSpan()
		if err != nil {
			loopSpan.Span().Finish()
			return err
		}

		loopSpan.Span().Finish()
	}

	return nil
}
