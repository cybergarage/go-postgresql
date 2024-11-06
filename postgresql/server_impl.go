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
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-tracing/tracer"
)

// server represents a PostgreSQL protocol server.
type server struct {
	Config
	*ConnManager
	tracer.Tracer
	tcpListener net.Listener
	*BaseExecutor
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		Config:       NewDefaultConfig(),
		ConnManager:  NewConnManager(),
		Tracer:       tracer.NullTracer,
		tcpListener:  nil,
		BaseExecutor: NewBaseExecutor(),
	}
	return server
}

// SetTracer sets a tracing tracer.
func (server *server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// Start starts the server.
func (server *server) Start() error {
	err := server.ConnManager.Start()
	if err != nil {
		return err
	}

	err = server.open()
	if err != nil {
		return err
	}

	go server.serve()

	addr := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	log.Infof("%s/%s (%s) started", PackageName, Version, addr)

	return nil
}

// Stop stops the server.
func (server *server) Stop() error {
	if err := server.ConnManager.Stop(); err != nil {
		return err
	}

	err := server.close()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)

	return nil
}

// Restart restarts the server.
func (server *server) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	return server.Start()
}

// open opens a listen socket.
func (server *server) open() error {
	var err error
	addr := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	server.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

// close closes a listening socket.
func (server *server) close() error {
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
func (server *server) serve() error {
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
func (server *server) receive(netConn net.Conn) error { //nolint:gocyclo,maintidx
	defer netConn.Close()

	log.Debugf("%s/%s (%s) accepted", PackageName, Version, netConn.RemoteAddr().String())

	handleStartupMessage := func(conn Conn, startupMsg *protocol.Startup) error {
		// PostgreSQL: Documentation: 16: 55.2. Message Flow
		// https://www.postgresql.org/docs/16/protocol-flow.html
		// Handle the Start-up message and return an Authentication message or error protocol.
		res, err := server.BaseExecutor.Authenticate(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ParameterStatus (S) protocol.
		reses, err := server.BaseExecutor.ParameterStatuses(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessages(reses)
		if err != nil {
			return err
		}
		// Return BackendKeyData (K) protocol.
		res, err = server.BaseExecutor.BackendKeyData(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ReadyForQuery (B) protocol.
		err = conn.ReadyForMessage(protocol.TransactionIdle)
		if err != nil {
			return err
		}
		return nil
	}

	conn := NewConnWith(netConn)
	defer func() {
		conn.Close()
	}()

	// Checks the SSLRequest protocol.

	startupMsgLength, err := conn.MessageReader().PeekInt32()
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	// PostgreSQL: Documentation: 16: 55.2.10. SSL Session Encryption
	// https://www.postgresql.org/docs/16/protocol-flow.html
	if startupMsgLength == 8 {
		_, err := protocol.NewSSLRequestWithReader(conn.MessageReader())
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
			err = conn.ResponseMessage(protocol.NewSSLResponseWith(protocol.SSLDisabled))
			if err != nil {
				return err
			}
		}
		err = conn.ResponseMessage(protocol.NewSSLResponseWith(protocol.SSLEnabled))
		if err != nil {
			return err
		}
		tlsConn := tls.Server(conn, tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			conn.ResponseError(err)
			return err
		}
		tlsConnState := tlsConn.ConnectionState()
		conn = NewConnWith(tlsConn, WithTLSConnectionState(&tlsConnState))
	}

	// Handle a Start-up protocol.

	startupMsg, err := protocol.NewStartupWithReader(conn.MessageReader())
	if err != nil {
		conn.ResponseError(err)
		return err
	}
	conn.SetStartupMessage(startupMsg)

	err = handleStartupMessage(conn, startupMsg)
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	// Handle the request messages after the Start-up protocol.

	dbname := ""
	if db, ok := startupMsg.Database(); ok {
		dbname = db
	}
	conn.SetDatabase(dbname)

	// Add the connection to the connection manager.

	server.AddConn(conn)
	defer func() {
		server.RemoveConn(conn)
	}()

	for {
		var reqErr error
		var reqType protocol.Type
		reqType, reqErr = conn.MessageReader().PeekType()
		if reqErr != nil {
			conn.ResponseError(reqErr)
			break
		}

		loopSpan := server.Tracer.StartSpan(PackageName)
		conn.SetSpanContext(loopSpan)
		conn.StartSpan(reqType.String())

		var resMsgs protocol.Responses

		switch reqType { // nolint:exhaustive
		case protocol.ParseMessage:
			var reqMsg *protocol.Parse
			reqMsg, reqErr = protocol.NewParseWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Parse(conn, reqMsg)
			}
		case protocol.BindMessage:
			var reqMsg *protocol.Bind
			reqMsg, reqErr = protocol.NewBindWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Bind(conn, reqMsg)
			}
		case protocol.DescribeMessage:
			var reqMsg *protocol.Describe
			reqMsg, reqErr = protocol.NewDescribeWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Describe(conn, reqMsg)
			}
		case protocol.QueryMessage:
			var reqMsg *protocol.Query
			reqMsg, reqErr = protocol.NewQueryWithReader(conn.MessageReader())
			if reqErr == nil {
				resMsgs, reqErr = server.BaseExecutor.Query(conn, reqMsg)
			}
		case protocol.ExecuteMessage:
			var reqMsg *protocol.Execute
			reqMsg, err := protocol.NewExecuteWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Execute(conn, reqMsg)
			}
		case protocol.CloseMessage:
			var reqMsg *protocol.Close
			reqMsg, err := protocol.NewCloseWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Close(conn, reqMsg)
			}
		case protocol.SyncMessage:
			var reqMsg *protocol.Sync
			reqMsg, err := protocol.NewSyncWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Sync(conn, reqMsg)
			}
		case protocol.FlushMessage:
			var reqMsg *protocol.Flush
			reqMsg, err := protocol.NewFlushWithReader(conn.MessageReader())
			if err == nil {
				resMsgs, reqErr = server.BaseExecutor.Flush(conn, reqMsg)
			}
		case protocol.TerminateMessage:
			_, reqErr = protocol.NewTerminateWithReader(conn.MessageReader())
			if reqErr == nil {
				conn.FinishSpan()
				loopSpan.Span().Finish()
				return nil
			}
		default:
			reqErr = conn.SkipMessage()
			if reqErr == nil {
				reqErr = protocol.NewErrMessageNotSuppoted(reqType)
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

		// Return ReadyForQuery (B) protocol.
		conn.StartSpan("ready")
		err := conn.ReadyForMessage(protocol.TransactionIdle)
		conn.FinishSpan()
		if err != nil {
			loopSpan.Span().Finish()
			return err
		}

		loopSpan.Span().Finish()
	}

	return nil
}
