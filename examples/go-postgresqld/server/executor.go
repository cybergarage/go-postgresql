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

package server

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

import (
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// Copy handles a COPY query.
func (server *Server) Copy(conn postgresql.Conn, q query.Copy) (protocol.Responses, error) {
	_, tbl, err := server.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		return nil, err
	}

	return postgresql.NewCopyInResponsesFrom(q, tbl.Schema)
}

// Copy handles a COPY DATA protocol.
func (server *Server) CopyData(conn postgresql.Conn, q query.Copy, stream *postgresql.CopyStream) (protocol.Responses, error) {
	_, tbl, err := server.LookupDatabaseTable(conn, conn.Database(), q.TableName())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return postgresql.NewCopyCompleteResponsesFrom(q, stream, conn, tbl.Schema, server.QueryExecutor())
}
