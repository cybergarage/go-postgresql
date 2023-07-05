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
	"bufio"
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

// SetExecutor sets a executor.
func (server *Server) SetExecutor(e Executor) {
	server.Executor = e
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
func (server *Server) receive(conn net.Conn) error {
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

	handleStartupMessage := func(exConn *Conn, startupMsg *message.Startup) error {
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
		readyMsg, err := message.NewReadyForQueryWith(message.TransactionIdle)
		if err != nil {
			return err
		}
		err = responseMessage(readyMsg)
		if err != nil {
			return err
		}
		return nil
	}

	handleParseBindMessage := func(exConn *Conn, reqMsg *message.RequestMessage) (*message.Query, error) {
		parseMsg, err := reqMsg.ParseParseMessage()
		if err != nil {
			return nil, err
		}

		resMsg, err := server.Executor.Parse(exConn, parseMsg)
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
			return nil, message.NewMessageNotSuppoted(reqType)
		}

		bindMsg, err := reqMsg.ParseBindMessageWith(parseMsg)
		if err != nil {
			return nil, err
		}

		resMsg, err = server.Executor.Bind(exConn, parseMsg, bindMsg)
		if err != nil {
			return nil, err
		}

		err = responseMessage(resMsg)
		if err != nil {
			return nil, err
		}

		q := &message.Query{}

		return q, nil
	}

	isStartupMessage := true
	var lastErr error
	for lastErr == nil {
		loopSpan := server.Tracer.StartSpan(PackageName)
		exConn := NewConnWith(conn, loopSpan)
		loopSpan.StartSpan("parse")

		reqMsg := message.NewRequestMessageWith(bufio.NewReader(conn))

		// Handle a Start-up message.

		if isStartupMessage {
			isStartupMessage = false
			msg, err := reqMsg.ParseStartupMessage()
			if err != nil {
				// Return ErrorResponse (B) and close the error connection.
				responseError(err)
				return err
			}

			err = handleStartupMessage(exConn, msg)
			if err != nil {
				// Return ErrorResponse (B) and close the error connection.
				responseError(err)
				return err
			}
			continue
		}

		// Handle the request messages after the Start-up message.

		var reqType message.Type
		reqType, lastErr = reqMsg.ReadType()
		if lastErr != nil {
			responseError(lastErr)
			break
		}

		var resMsg message.Response
		switch reqType { // nolint:exhaustive
		case message.ParseMessage:
			var queryMsg *message.Query
			queryMsg, err := handleParseBindMessage(exConn, reqMsg)
			if err == nil {
				log.Infof("Query: %s", queryMsg.Query)
				lastErr = err
			}
		case message.QueryMessage:
			var queryMsg *message.Query
			queryMsg, lastErr = reqMsg.ParseQueryMessage()
			if lastErr == nil {
				log.Infof("Query: %s", queryMsg.Query)
				lastErr = message.NewMessageNotSuppoted(reqType)
			}
		case message.TerminateMessage:
			return nil
		default:
			lastErr = message.NewMessageNotSuppoted(reqType)
			log.Warnf(lastErr.Error())
		}

		if lastErr == nil {
			if resMsg != nil {
				lastErr = responseMessage(resMsg)
			}
		} else {
			// Return ErrorResponse (B)
			responseError(lastErr)
		}

		loopSpan.FinishSpan()
	}

	return lastErr
}
