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

// DataRow represents a data row message.
type DataRow struct {
	*ResponseMessage
	Data []any
}

// NewDataRow returns a new data row message instance.
func NewDataRow() *DataRow {
	return &DataRow{
		ResponseMessage: NewResponseMessageWith(DataRowMessage),
	}
}

// AppendData appends a column value to the data row message.
func (msg *DataRow) AppendData(v any) error {
	msg.Data = append(msg.Data, v)
	return nil
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *DataRow) Bytes() ([]byte, error) {
	err := msg.AppendInt16(int16(len(msg.Data)))
	if err != nil {
		return nil, err
	}
	for _, v := range msg.Data {
		switch v := v.(type) {
		case []byte:
			if err := msg.AppendInt32(int32(len(v))); err != nil {
				return nil, err
			}
			if err := msg.AppendBytes(v); err != nil {
				return nil, err
			}
		case string:
			if err := msg.AppendInt32(int32(len(v))); err != nil {
				return nil, err
			}
			if err := msg.AppendBytes([]byte(v)); err != nil {
				return nil, err
			}
		case int8:
			if err := msg.AppendInt32(1); err != nil {
				return nil, err
			}
			if err := msg.AppendInt8(v); err != nil {
				return nil, err
			}
		case int16:
			if err := msg.AppendInt32(2); err != nil {
				return nil, err
			}
			if err := msg.AppendInt16(v); err != nil {
				return nil, err
			}
		case int32:
			if err := msg.AppendInt32(4); err != nil {
				return nil, err
			}
			if err := msg.AppendInt32(v); err != nil {
				return nil, err
			}
		case int64:
			if err := msg.AppendInt32(8); err != nil {
				return nil, err
			}
			if err := msg.AppendInt64(v); err != nil {
				return nil, err
			}
		case int:
			if err := msg.AppendInt32(8); err != nil {
				return nil, err
			}
			if err := msg.AppendInt64(int64(v)); err != nil {
				return nil, err
			}
		case nil:
			if err := msg.AppendInt32(-1); err != nil {
				return nil, err
			}
		default:
			return nil, newColumnTypeNotSuppotedError(v)
		}
	}
	return msg.ResponseMessage.Bytes()
}
