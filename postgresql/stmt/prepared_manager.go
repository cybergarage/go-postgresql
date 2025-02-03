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
	"github.com/google/uuid"
)

// PreparedManager represents a prepared manager.
type PreparedManager struct {
	stateMap map[uuid.UUID]PreparedStatementMap
	potalMap map[uuid.UUID]PreparedPortalMap
}

// NewPreparedManager returns a new prepared manager.
func NewPreparedManager() *PreparedManager {
	return &PreparedManager{
		stateMap: make(map[uuid.UUID]PreparedStatementMap),
		potalMap: make(map[uuid.UUID]PreparedPortalMap),
	}
}

// PreparedStatement returns a prepared statement.
func (mgr *PreparedManager) PreparedStatement(conn Conn, name string) (*PreparedStatement, error) {
	preState, ok := mgr.stateMap[conn.UUID()]
	if !ok {
		return nil, errors.NewErrPreparedStatementNotExist(name)
	}
	return preState.PreparedStatement(name)
}

// SetPreparedStatement sets a prepared statement.
func (mgr *PreparedManager) SetPreparedStatement(conn Conn, msg *protocol.Parse) error {
	preState, ok := mgr.stateMap[conn.UUID()]
	if !ok {
		preState = NewPreparedStatementMap()
		mgr.stateMap[conn.UUID()] = preState
	}
	return preState.SetPreparedStatement(msg)
}

// RemovePreparedStatement removes a prepared statement.
func (mgr *PreparedManager) RemovePreparedStatement(conn Conn, name string) error {
	preState, ok := mgr.stateMap[conn.UUID()]
	if !ok {
		return errors.NewErrPreparedStatementNotExist(name)
	}
	return preState.RemovePreparedStatement(name)
}

// PreparedPortal returns a prepared query statement.
func (mgr *PreparedManager) PreparedPortal(conn Conn, name string) (*PreparedPortal, error) {
	prePortal, ok := mgr.potalMap[conn.UUID()]
	if !ok {
		return nil, errors.NewErrPreparedPortalNotExist(name)
	}
	return prePortal.PreparedPortal(name)
}

// SetPreparedPortal sets a prepared query statement.
func (mgr *PreparedManager) SetPreparedPortal(conn Conn, name string, query *PreparedPortal) error {
	prePortal, ok := mgr.potalMap[conn.UUID()]
	if !ok {
		prePortal = NewPreparedPortalMap()
		mgr.potalMap[conn.UUID()] = prePortal
	}
	return prePortal.SetPreparedPortal(name, query)
}

// RemovePreparedPortal removes a prepared query statement.
func (mgr *PreparedManager) RemovePreparedPortal(conn Conn, name string) error {
	prePortal, ok := mgr.potalMap[conn.UUID()]
	if !ok {
		return errors.NewErrPreparedPortalNotExist(name)
	}
	return prePortal.RemovePreparedPortal(name)
}
