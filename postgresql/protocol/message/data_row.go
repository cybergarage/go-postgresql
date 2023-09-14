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
	"time"

	"github.com/cybergarage/go-postgresql/postgresql/system"
	"github.com/cybergarage/go-safecast/safecast"
)

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
func (msg *DataRow) AppendData(rowField *RowField, v any) error {
	switch rowField.FormatCode { //nolint:exhaustive
	case system.TextFormat:
		switch sv := v.(type) {
		case string:
			// v = sv
		case time.Time:
			v = sv.Format(system.TimestampFormat)
		case nil:
			v = nil
		default:
			var to string
			if err := safecast.ToString(v, &to); err != nil {
				return err
			}
			v = to
		}
	case system.BinaryFormat:
		switch rowField.DataTypeID { //nolint:exhaustive
		case system.Bool:
			if _, ok := v.(bool); !ok {
				var to bool
				if err := safecast.ToBool(v, &to); err != nil {
					return err
				}
				v = to
			}
		case system.Int2:
			if _, ok := v.(int16); !ok {
				var to int16
				if err := safecast.ToInt16(v, &to); err != nil {
					return err
				}
				v = to
			}
		case system.Int4:
			if _, ok := v.(int32); !ok {
				var to int32
				if err := safecast.ToInt32(v, &to); err != nil {
					return err
				}
				v = to
			}
		case system.Int8:
			if _, ok := v.(int64); !ok {
				var to int64
				if err := safecast.ToInt64(v, &to); err != nil {
					return err
				}
				v = to
			}
		case system.Float4:
			if _, ok := v.(float32); !ok {
				var to float32
				if err := safecast.ToFloat32(v, &to); err != nil {
					return err
				}
				v = to
			}
		case system.Float8:
			if _, ok := v.(float64); !ok {
				var to float64
				if err := safecast.ToFloat64(v, &to); err != nil {
					return err
				}
				v = to
			}
		case system.Text, system.Varchar:
			if _, ok := v.(string); !ok {
				var to string
				if err := safecast.ToString(v, &to); err != nil {
					return err
				}
				v = to
			}
		}
	}
	msg.Data = append(msg.Data, v)
	return nil
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *DataRow) Bytes() ([]byte, error) { // nolint:gocyclo
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
		case float32:
			if err := msg.AppendInt32(4); err != nil {
				return nil, err
			}
			if err := msg.AppendFloat32(v); err != nil {
				return nil, err
			}
		case float64:
			if err := msg.AppendInt32(8); err != nil {
				return nil, err
			}
			if err := msg.AppendFloat64(v); err != nil {
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
