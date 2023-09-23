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

import (
	"strconv"
	"strings"
)

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

const (
	bindParamPrefix = "$"
)

// Bind represents a bind message.
type Bind struct {
	*RequestMessage
	PortalName    string
	StatementName string
	NumParams     int16
	Params        BindParams
}

// BindParam represents a bind parameter.
type BindParam struct {
	FormatCode int16
	Value      any
}

// BindParams represents bind parameters.
type BindParams []*BindParam

// NewBind returns a new bind message.
func NewBindWithReader(reader *MessageReader) (*Bind, error) {
	msg, err := NewRequestMessageWithReader(reader)
	if err != nil {
		return nil, err
	}

	// The name of the destination portal (an empty string selects the unnamed portal).
	portal, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// The stmt of the source prepared statement (an empty string selects the unnamed prepared statement).
	stmt, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// The number of parameter format codes that follow (denoted C below).
	// This can be zero to indicate that there are no parameters or that the parameters all use the default format (text); or one, in which case the specified format code is applied to all parameters; or it can equal the actual number of parameters.
	paramFmtNum, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	// The parameter format codes. Each must presently be zero (text) or one (binary).
	paramFmts := make([]int16, paramFmtNum)
	for n := 0; n < int(paramFmtNum); n++ {
		fmt, err := reader.ReadInt16()
		if err != nil {
			return nil, err
		}
		paramFmts[n] = fmt
	}

	// The number of parameter values that follow (possibly zero). This must match the number of parameters needed by the query.
	paramValNum, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	params := make([]*BindParam, paramValNum)
	for n := 0; n < int(paramValNum); n++ {
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

		paramFmt := TextFormat
		if n < len(paramFmts) {
			paramFmt = paramFmts[n]
		}
		var paramVal any
		switch paramFmt {
		case TextFormat:
			paramVal = string(bytes)
		case BinaryFormat:
			paramVal = bytes
		}
		params[n] = &BindParam{
			FormatCode: paramFmt,
			Value:      paramVal,
		}
	}

	// The number of result-column format codes that follow (denoted R below).
	resFmtNum, err := reader.ReadInt16()
	if err != nil {
		return nil, err
	}

	// The result-column format codes. Each must presently be zero (text) or one (binary).
	resFmts := make([]int16, resFmtNum)
	for n := 0; n < int(resFmtNum); n++ {
		fmt, err := reader.ReadInt16()
		if err != nil {
			return nil, err
		}
		resFmts[n] = fmt
	}

	return &Bind{
		RequestMessage: msg,
		PortalName:     portal,
		StatementName:  stmt,
		NumParams:      paramFmtNum,
		Params:         params,
	}, nil
}

// FindBindParam returns a bind parameter with specified id.
func (params BindParams) FindBindParam(id string) (*BindParam, error) {
	if !strings.HasPrefix(id, bindParamPrefix) {
		return nil, NewErrNotExist(id)
	}
	idx, err := strconv.Atoi(id[len(bindParamPrefix):])
	if err != nil {
		return nil, NewErrNotExist(id)
	}
	if len(params) < idx {
		return nil, NewErrNotExist(id)
	}
	return params[idx-1], nil
}
