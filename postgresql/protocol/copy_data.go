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

import (
	"bytes"
	"io"
	"strings"
)

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: COPY
// https://www.postgresql.org/docs/16/sql-copy.html

const (
	tabSep     = '\t'
	newLineSep = "\r\n"
)

// CopyData represents a copy data protocol.
type CopyData struct {
	*RequestMessage

	Data []string
}

// NewCopyDataWithReader returns a new copy data message with the specified reader.
func NewCopyDataWithReader(reader *MessageReader) (*CopyData, error) {
	msg, err := NewRequestMessageWithReader(reader)
	if err != nil {
		return nil, err
	}

	dataLen := msg.MessageDataLength()
	if dataLen < 0 {
		return nil, newInvalidLengthError(int(dataLen))
	}

	dataBytes := make([]byte, dataLen)
	_, err = reader.ReadBytes(dataBytes)
	if err != nil {
		return nil, err
	}

	dataBytes = bytes.TrimRight(dataBytes, newLineSep)

	isEOFData := func(data []byte) bool {
		for _, b := range data {
			if b != 0x5C && b != 0x2E {
				return false
			}
		}
		return true
	}

	if isEOFData(dataBytes) {
		return nil, io.EOF
	}

	data := strings.Split(string(dataBytes), string(tabSep))

	return &CopyData{
		RequestMessage: msg,
		Data:           data,
	}, nil
}
