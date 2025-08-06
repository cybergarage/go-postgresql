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
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

// RowDescription represents a row description protocol.
type RowDescription struct {
	*ResponseMessage
	fileds []*RowField
}

// NewRowDescription returns a new row description message instance.
func NewRowDescription() *RowDescription {
	return &RowDescription{
		ResponseMessage: NewResponseMessageWith(RowDescriptionMessage),
		fileds:          []*RowField{},
	}
}

// AppendField appends a field to the protocol.
func (msg *RowDescription) AppendField(field *RowField) {
	msg.fileds = append(msg.fileds, field)
}

// Field returns a field at the specified index.
func (msg *RowDescription) Field(n int) *RowField {
	return msg.fileds[n]
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *RowDescription) Bytes() ([]byte, error) {
	msg.AppendInt16(int16(len(msg.fileds)))
	for _, field := range msg.fileds {
		err := field.WirteBytes(msg.Writer)
		if err != nil {
			return nil, err
		}
	}
	return msg.ResponseMessage.Bytes()
}
