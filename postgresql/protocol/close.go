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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// Close represents a close protocol.
type Close struct {
	*RequestMessage
	Type PreparedType
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

	dt, err := NewPreparedTypeWithByte(bt)
	if err != nil {
		return nil, err
	}

	return &Close{
		RequestMessage: msg,
		Type:           dt,
		Name:           name,
	}, nil
}
