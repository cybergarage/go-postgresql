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

package stmt

import (
	"github.com/cybergarage/go-postgresql/postgresql/errors"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// PreparedStatement represents a prepared statement.
type PreparedStatement struct {
	*protocol.Parse

	ParsedStatement *query.Statement
}

// Name returns a prepared statement name.
func (stmt *PreparedStatement) Name() string {
	return stmt.Parse.Name
}

// PreparedStatementMap represents a prepared statement map.
type PreparedStatementMap map[string]*PreparedStatement

// NewPreparedStatementMap returns a new prepared statement map.
func NewPreparedStatementMap() PreparedStatementMap {
	return make(PreparedStatementMap)
}

// PreparedStatement returns a prepared statement.
func (stmtMap PreparedStatementMap) PreparedStatement(name string) (*PreparedStatement, error) {
	q, oK := stmtMap[name]
	if !oK {
		return nil, errors.NewErrPreparedStatementNotExist(name)
	}

	return q, nil
}

// SetPreparedStatement sets a prepared statement.
func (stmtMap PreparedStatementMap) SetPreparedStatement(msg *protocol.Parse) error {
	parser := query.NewParser()

	stmts, err := parser.ParseString(msg.Query)
	if err != nil {
		return err
	}

	if 1 < len(stmts) {
		return errors.NewErrMultiplePreparedStatementNotSupported(msg.Query)
	}

	stmt := &PreparedStatement{
		Parse:           msg,
		ParsedStatement: stmts[0],
	}
	stmtMap[msg.Name] = stmt

	return nil
}

// RemovePreparedStatement removes a prepared statement.
func (stmtMap PreparedStatementMap) RemovePreparedStatement(name string) error {
	_, oK := stmtMap[name]
	if !oK {
		return errors.NewErrPreparedStatementNotExist(name)
	}

	delete(stmtMap, name)

	return nil
}

// PreparedPortal represents a prepared query statement.
type PreparedPortal = protocol.Query

// PreparedPortalMap represents a prepared query statement map.
type PreparedPortalMap map[string]PreparedPortal

// NewPreparedPortalMap returns a new prepared query statement map.
func NewPreparedPortalMap() PreparedPortalMap {
	return make(PreparedPortalMap)
}

// PreparedPortal returns a prepared query statement.
func (portalMap PreparedPortalMap) PreparedPortal(name string) (*PreparedPortal, error) {
	q, oK := portalMap[name]
	if !oK {
		return nil, errors.NewErrPreparedPortalNotExist(name)
	}

	return &q, nil
}

// SetPreparedPortal sets a prepared query statement.
func (portalMap PreparedPortalMap) SetPreparedPortal(name string, query *PreparedPortal) error {
	portalMap[name] = *query
	return nil
}

// RemovePreparedPortal removes a prepared query statement.
func (portalMap PreparedPortalMap) RemovePreparedPortal(name string) error {
	_, oK := portalMap[name]
	if !oK {
		return errors.NewErrPreparedPortalNotExist(name)
	}

	delete(portalMap, name)

	return nil
}
