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

package protocol

import (
	"crypto/tls"
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/auth"
	pgnet "github.com/cybergarage/go-postgresql/postgresql/net"
	"github.com/cybergarage/go-tracing/tracer"
)

// server represents a PostgreSQL protocol server.
type server struct {
	Config
	*pgnet.ConnManager
	tracer.Tracer
	tcpListener net.Listener
	MessageHandler
	auth.Manager
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		Config:         NewDefaultConfig(),
		ConnManager:    pgnet.NewConnManager(),
		Tracer:         tracer.NullTracer,
		tcpListener:    nil,
		MessageHandler: nil,
		Manager:        auth.NewManager(),
	}
	return server
}

// SetTracer sets a tracing tracer.
func (server *server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

func (server *server) SetMessageHandler(h MessageHandler) {
	server.MessageHandler = h
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
	log.Infof("%s/%s (%s) started", server.ProductName(), server.ProductVersion(), addr)

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
	log.Infof("%s/%s (%s) terminated", server.ProductName(), server.ProductVersion(), addr)

	return nil
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

	log.Debugf("%s/%s (%s) accepted", server.ProductName(), server.ProductVersion(), netConn.RemoteAddr().String())

	authenticateStartup := func(conn Conn, startupMsg *Startup) (Response, error) {
		authenticate := func(conn Conn, startupMsg *Startup) (bool, error) {
			clientUsername, ok := startupMsg.User()
			if !ok {
				return false, nil
			}
			authMsg, err := NewAuthenticationCleartextPassword()
			if err != nil {
				return false, err
			}
			err = conn.ResponseMessage(authMsg)
			if err != nil {
				return false, err
			}
			msg, err := NewPasswordWithReader(conn.MessageReader())
			if err != nil {
				return false, err
			}
			q, err := auth.NewQuery(
				auth.WithQueryUsername(clientUsername),
				auth.WithQueryPassword(msg.Password),
			)
			if err != nil {
				return false, err
			}
			return server.VerifyCredential(conn, q)
		}

		ok, err := authenticate(conn, startupMsg)
		if err != nil {
			return nil, err
		}
		if ok {
			return NewAuthenticationOk()
		}
		return NewErrorResponse(), nil
	}

	handleStartupMessage := func(conn Conn, startupMsg *Startup) error {
		// PostgreSQL: Documentation: 16: 55.2. Message Flow
		// https://www.postgresql.org/docs/16/protocol-flow.html
		// Handle the Start-up message and return an Authentication message or error
		res, err := authenticateStartup(conn, startupMsg)
		if err != nil {
			return err
		}
		err = conn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ParameterStatus (S)
		reses, err := server.MessageHandler.ParameterStatuses(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessages(reses)
		if err != nil {
			return err
		}
		// Return BackendKeyData (K)
		res, err = server.MessageHandler.BackendKeyData(conn)
		if err != nil {
			return err
		}
		err = conn.ResponseMessage(res)
		if err != nil {
			return err
		}
		// Return ReadyForQuery (B)
		err = conn.ReadyForMessage(TransactionIdle)
		if err != nil {
			return err
		}
		return nil
	}

	conn := NewConnWith(netConn)
	defer func() {
		conn.Close()
	}()

	// Checks the SSLRequest

	startupMsgLength, err := conn.MessageReader().PeekInt32()
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	// PostgreSQL: Documentation: 16: 55.2.10. SSL Session Encryption
	// https://www.postgresql.org/docs/16/protocol-flow.html
	if startupMsgLength == 8 {
		_, err := NewSSLRequestWithReader(conn.MessageReader())
		if err != nil {
			conn.ResponseError(err)
			return err
		}
		tlsConfig, err := server.Config.TLSConfig()
		if err != nil {
			conn.ResponseError(err)
			return err
		}
		if tlsConfig != nil {
			err = conn.ResponseMessage(NewSSLResponseWith(SSLEnabled))
			if err != nil {
				return err
			}
			tlsConn := tls.Server(conn, tlsConfig)
			if err := tlsConn.Handshake(); err != nil {
				conn.ResponseError(err)
				return err
			}
			ok, err := server.Manager.VerifyCertificate(tlsConn)
			if !ok {
				conn.ResponseError(err)
				return err
			}
			conn = NewConnWith(tlsConn, WithConnTLSConn(tlsConn))
		} else {
			err = conn.ResponseMessage(NewSSLResponseWith(SSLDisabled))
			if err != nil {
				return err
			}
		}
	}

	// Get reader from the existing or new TLS connection

	reader := conn.MessageReader()

	// Handle a Start-up

	startupMsg, err := NewStartupWithReader(reader)
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	err = handleStartupMessage(conn, startupMsg)
	if err != nil {
		conn.ResponseError(err)
		return err
	}

	// Handle the request messages after the Start-up

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
		var reqType Type
		reqType, reqErr = reader.PeekType()
		if reqErr != nil {
			conn.ResponseError(reqErr)
			break
		}

		loopSpan := server.Tracer.StartSpan(server.ProductName())
		conn.SetSpanContext(loopSpan)
		conn.StartSpan(reqType.String())

		var resMsgs Responses

		switch reqType { // nolint:exhaustive
		case ParseMessage:
			var reqMsg *Parse
			reqMsg, reqErr = NewParseWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Parse(conn, reqMsg)
			}
		case BindMessage:
			var reqMsg *Bind
			reqMsg, reqErr = NewBindWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Bind(conn, reqMsg)
			}
		case DescribeMessage:
			var reqMsg *Describe
			reqMsg, reqErr = NewDescribeWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Describe(conn, reqMsg)
			}
		case QueryMessage:
			var reqMsg *Query
			reqMsg, reqErr = NewQueryWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Query(conn, reqMsg)
			}
		case ExecuteMessage:
			var reqMsg *Execute
			reqMsg, reqErr = NewExecuteWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Execute(conn, reqMsg)
			}
		case CloseMessage:
			var reqMsg *Close
			reqMsg, reqErr = NewCloseWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Close(conn, reqMsg)
			}
		case SyncMessage:
			var reqMsg *Sync
			reqMsg, reqErr = NewSyncWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Sync(conn, reqMsg)
			}
		case FlushMessage:
			var reqMsg *Flush
			reqMsg, reqErr = NewFlushWithReader(reader)
			if reqErr == nil {
				resMsgs, reqErr = server.MessageHandler.Flush(conn, reqMsg)
			}
		case TerminateMessage:
			_, reqErr = NewTerminateWithReader(reader)
			if reqErr == nil {
				conn.FinishSpan()
				loopSpan.Span().Finish()
				return nil
			}
		default:
			reqErr = conn.SkipMessage()
			if reqErr == nil {
				reqErr = NewErrMessageNotSuppoted(reqType)
				log.Warn(reqErr.Error())
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

		// Return ReadyForQuery (B)

		conn.StartSpan("ready")
		err := conn.ReadyForMessage(TransactionIdle)
		conn.FinishSpan()
		if err != nil {
			loopSpan.Span().Finish()
			return err
		}

		loopSpan.Span().Finish()
	}

	return nil
}
