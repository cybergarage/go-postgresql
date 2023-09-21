// Copyright (C) 2019 Satoshi Konno. All rights reserved.
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

// DescribeType represents a describe type.
type DescribeType int

const (
	// DescribeTypeStatementByte represents a describe statement.
	DescribeTypeStatementByte = 'S'
	// DescribeTypePortalByte represents a describe portal.
	DescribeTypePortalByte = 'P'
)

const (
	// DescribeTypeStatement represents a describe statement.
	DescribeTypeStatement DescribeType = iota
	// DescribeTypePortal represents a describe portal.
	DescribeTypePortal
)

// Describe represents a describe message.
type Describe struct {
	*RequestMessage
	Type DescribeType
	Name string
}

// NewDescribeWithReader returns a new describe message with the specified reader.
func NewDescribeWithReader(reader *MessageReader) (*Describe, error) {
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

	var dt DescribeType
	switch bt {
	case DescribeTypeStatementByte:
		dt = DescribeTypeStatement
	case DescribeTypePortalByte:
		dt = DescribeTypePortal
	default:
		return nil, fmt.Errorf("%w describe type (%d)", ErrInvalid, bt)
	}

	return &Describe{
		RequestMessage: msg,
		Type:           dt,
		Name:           name,
	}, nil
}

// IsStatement returns true whether the describe type is statement.
func (msg *Describe) IsStatement() bool {
	return msg.Type == DescribeTypeStatement
}

// IsPortal returns true whether the describe type is portal.
func (msg *Describe) IsPortal() bool {
	return msg.Type == DescribeTypePortal
}
