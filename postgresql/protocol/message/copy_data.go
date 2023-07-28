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

// CopyData represents a copy data message.
type CopyData struct {
	MessageLength int32
	Data          []byte
}

// NewCopyDataWithReader returns a new copy data message with the specified reader.
func NewCopyDataWithReader(reader *Reader) (*CopyData, error) {
	msgLen, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	dataLen := msgLen - 4
	if dataLen < 0 {
		return nil, newInvalidLengthError(int(dataLen))
	}

	data := make([]byte, dataLen)
	_, err = reader.ReadBytes(data)
	if err != nil {
		return nil, err
	}

	return &CopyData{
		MessageLength: msgLen,
		Data:          data,
	}, nil
}
