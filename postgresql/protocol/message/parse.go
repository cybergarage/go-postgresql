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

// NewParseWithReader returns a new parse message with the specified reader.
func NewParseWithReader(reader *Reader) (*Parse, error) {
	msgLen, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	// The name of the destination prepared statement (an empty string selects the unnamed prepared statement).
	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// The query string to be parsed.
	query, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// The number of parameter data types specified (can be zero).
	// Note that this is not an indication of the number of parameters that might appear in the query string, only the number that the frontend wants to prespecify types for.
	num, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	types := make([]int32, num)
	for n := int16(0); n < num; n++ {
		// Specifies the object ID of the parameter data type. Placing a zero here is equivalent to leaving the type unspecified.
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
