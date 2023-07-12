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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// Parse represents a parse message.
type Parse struct {
	Name          string
	Query         string
	MessageLength int32
	NumDataTypes  int16
	DataTypes     []int32
}

// NewParse returns a new parse message.
func NewParseWithReader(reader *Reader) (*Parse, error) {
	msgLen, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	query, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	num, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	types := make([]int32, num)
	for n := int16(0); n < num; n++ {
		typ, err := reader.ReadInt32()
		if err != nil {
			return nil, err
		}
		types[n] = typ
	}

	return &Parse{
		MessageLength: msgLen,
		Name:          name,
		Query:         query,
		NumDataTypes:  num,
		DataTypes:     types,
	}, nil
}
