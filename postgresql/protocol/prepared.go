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

import "fmt"

// PostgreSQL: Documentation: 16: 55.1. Overview
// https://www.postgresql.org/docs/current/protocol-overview.html
// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// PreparedType represents a prepared type.
type PreparedType int

const (
	// PreparedStatementByte represents a prepared statement.
	PreparedStatementByte = 'S'
	// PreparedPortalByte represents a prepared portal.
	PreparedPortalByte = 'P'
)

const (
	// PreparedStatement represents a prepared statement.
	PreparedStatement PreparedType = iota
	// PreparedPortal represents a prepared portal.
	PreparedPortal
)

// NewPreparedTypeWithByte returns a new prepared type with the specified byte.
func NewPreparedTypeWithByte(bt byte) (PreparedType, error) {
	switch bt {
	case PreparedStatementByte:
		return PreparedStatement, nil
	case PreparedPortalByte:
		return PreparedPortal, nil
	}
	return 0, fmt.Errorf("%w prepared type (%d)", ErrInvalid, bt)
}
