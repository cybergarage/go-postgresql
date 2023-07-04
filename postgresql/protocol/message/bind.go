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

// Bind represents a bind message.
type Bind struct {
	MessageLength int32
	Portal        string
	Name          string
	NumParams     int16
	Params        []*BindParam
}

// BindParamType represents a bind parameter type.
type BindParamType int16

const (
	// BindParamTypeString represents a string type.
	BindParamTypeString BindParamType = 0
	// BindParamTypeBinary represents a binary type.
	BindParamTypeBinary BindParamType = 1
)

// BindParam represents a bind parameter.
type BindParam struct {
	Type      BindParamType
	numValues int16
	Value     any
}

// NewBind returns a new bind message.
func NewBindWith(reader *Reader) (*Bind, error) {
	msgLen, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	portal, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	num, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	params := make([]*BindParam, num)
	for n := int16(0); n < num; n++ {
		t, err := reader.ReadInt16()
		if err != nil {
			return nil, err
		}
		num, err := reader.ReadInt16()
		if err != nil {
			return nil, err
		}
		nBytes, err := reader.ReadInt32()
		if err != nil {
			return nil, err
		}
		bytes := make([]byte, nBytes)
		nRead, err := reader.Read(bytes)
		if err != nil {
			return nil, err
		}
		if nRead != int(nBytes) {
			return nil, newShortMessageErrorWith(int(nBytes), nRead)
		}
		var v any
		if BindParamType(t) == BindParamTypeString {
			v = string(bytes)
		} else {
			v = bytes
		}
		params[n] = &BindParam{
			Type:      BindParamType(t),
			numValues: num,
			Value:     v,
		}
	}

	return &Bind{
		MessageLength: msgLen,
		Portal:        portal,
		Name:          name,
		NumParams:     num,
		Params:        params,
	}, nil
}
