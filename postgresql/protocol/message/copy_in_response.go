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
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

// CopyFormat represents a copy format.
type CopyFormat = int8

const (
	// TextCopy represents a textual copy format.
	TextCopy = CopyFormat(0)
	// BinaryCopy represents a binary copy format.
	BinaryCopy = CopyFormat(1)
)

// CopyInResponse represents a command complete message.
type CopyInResponse struct {
	*ResponseMessage
	formatCodes []int16
}

// NewCopyInResponse returns a new command complete message instance.
func NewCopyInResponseWith(fmt CopyFormat) *CopyInResponse {
	msg := &CopyInResponse{
		ResponseMessage: NewResponseMessageWith(CopyInResponseMessage),
		formatCodes:     []int16{},
	}
	msg.AppendInt8(fmt)
	return msg
}

// AppendFormatCode appends a format code.
func (msg *CopyInResponse) AppendFormatCode(formatCode int16) {
	msg.formatCodes = append(msg.formatCodes, formatCode)
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *CopyInResponse) Bytes() ([]byte, error) {
	msg.AppendInt16(int16(len(msg.formatCodes)))
	for _, field := range msg.formatCodes {
		err := msg.AppendInt16(field)
		if err != nil {
			return nil, err
		}
	}
	return msg.ResponseMessage.Bytes()
}
