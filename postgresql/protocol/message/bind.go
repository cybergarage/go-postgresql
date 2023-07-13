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
func NewBindWith(reader *Reader, parse *Parse) (*Bind, error) {
	msgLen, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	// The name of the destination portal (an empty string selects the unnamed portal).
	portal, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// The name of the source prepared statement (an empty string selects the unnamed prepared statement).
	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// The number of parameter format codes that follow (denoted C below).
	// This can be zero to indicate that there are no parameters or that the parameters all use the default format (text); or one, in which case the specified format code is applied to all parameters; or it can equal the actual number of parameters.
	_, err = reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	// The parameter format codes. Each must presently be zero (text) or one (binary).
	t, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	// The number of parameter values that follow (possibly zero). This must match the number of parameters needed by the query.
	num, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	params := make([]*BindParam, num)
	for n := int16(0); n < num; n++ {
		// The length of the parameter value, in bytes (this count does not include itself).
		// Can be zero. As a special case, -1 indicates a NULL parameter value. No value bytes follow in the NULL case.
		nBytes, err := reader.ReadInt32()
		if err != nil {
			return nil, err
		}
		// The value of the parameter, in the format indicated by the associated format code. n is the above length.
		var bytes []byte
		switch nBytes {
		case 0:
			bytes = []byte{}
		case -1:
			bytes = nil
		default:
			if nBytes <= 0 {
				return nil, newInvalidLengthError(int(nBytes))
			}
			bytes = make([]byte, nBytes)
			nRead, err := reader.Read(bytes)
			if err != nil {
				return nil, err
			}
			if nRead != int(nBytes) {
				return nil, newShortMessageError(int(nBytes), nRead)
			}
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

	// The number of result-column format codes that follow (denoted R below).
	_, err = reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	// The result-column format codes. Each must presently be zero (text) or one (binary).
	_, err = reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	return &Bind{
		MessageLength: msgLen,
		Portal:        portal,
		Name:          name,
		NumParams:     num,
		Params:        params,
	}, nil
}
