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

package message

import "fmt"

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// CloseType represents a close type.
type CloseType int

const (
	// CloseTypeStatementByte represents a close statement.
	CloseTypeStatementByte = 'S'
	// CloseTypePortalByte represents a close portal.
	CloseTypePortalByte = 'P'
)

const (
	// CloseTypeStatement represents a close statement.
	CloseTypeStatement CloseType = iota
	// CloseTypePortal represents a close portal.
	CloseTypePortal
)

// Close represents a close message.
type Close struct {
	*RequestMessage
	Type CloseType
	Name string
}

// NewCloseWithReader returns a new close message with the specified reader.
func NewCloseWithReader(reader *MessageReader) (*Close, error) {
	msg, err := NewRequestMessageWithReader(reader)
	if err != nil {
		return nil, err
	}

	bt, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	var dt CloseType
	switch bt {
	case CloseTypeStatementByte:
		dt = CloseTypeStatement
	case CloseTypePortalByte:
		dt = CloseTypePortal
	default:
		return nil, fmt.Errorf("%w close type (%d)", ErrInvalid, bt)
	}

	return &Close{
		RequestMessage: msg,
		Type:           dt,
		Name:           name,
	}, nil
}

// IsStatement returns true whether the close type is statement.
func (msg *Close) IsStatement() bool {
	return msg.Type == CloseTypeStatement
}

// IsPortal returns true whether the close type is portal.
func (msg *Close) IsPortal() bool {
	return msg.Type == CloseTypePortal
}
