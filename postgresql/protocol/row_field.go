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

// RowField represents a row description field.
type RowField struct {
	Name         string
	TableID      int32
	Number       int16
	ObjectID     int32
	DataTypeSize int16
	TypeModifier int32
	FormatCode   int16
}

// RowFieldOption represents a row description field option.
type RowFieldOption = func(*RowField)

// NewRowField returns a new row description field.
func NewRowFieldWith(name string, opts ...RowFieldOption) *RowField {
	field := &RowField{
		Name:         name,
		TableID:      0,
		Number:       0,
		ObjectID:     0,
		DataTypeSize: 0,
		TypeModifier: 0,
		FormatCode:   0,
	}
	for _, opt := range opts {
		opt(field)
	}

	return field
}

// WithRowFieldNumber sets a number.
func WithRowFieldNumber(number int16) func(*RowField) {
	return func(fileld *RowField) {
		fileld.Number = number
	}
}

// WithRowFieldTableID sets a table ID.
func WithRowFieldTableID(tableID int32) func(*RowField) {
	return func(fileld *RowField) {
		fileld.TableID = tableID
	}
}

// WithRowFieldObjectID sets a data type ID.
func WithRowFieldObjectID(dataTypeID int32) func(*RowField) {
	return func(fileld *RowField) {
		fileld.ObjectID = dataTypeID
	}
}

// WithRowFieldSize sets a data type size.
func WithRowFieldSize(dataTypeSize int16) func(*RowField) {
	return func(fileld *RowField) {
		fileld.DataTypeSize = dataTypeSize
	}
}

// WithRowFieldModifier sets a type modifier.
func WithRowFieldModifier(typeModifier int32) func(*RowField) {
	return func(fileld *RowField) {
		fileld.TypeModifier = typeModifier
	}
}

// WithRowFieldFormatCode sets a format code.
func WithRowFieldFormatCode(formatCode int16) func(*RowField) {
	return func(fileld *RowField) {
		fileld.FormatCode = formatCode
	}
}

// WithRowFieldDataType sets a data type.
func WithRowFieldDataType(dt *DataType) func(*RowField) {
	return func(fileld *RowField) {
		fileld.ObjectID = dt.OID()
		fileld.DataTypeSize = int16(dt.Size())
		fileld.FormatCode = dt.FormatCode()
	}
}

// WirteBytes appends a row field elements.
func (field *RowField) WirteBytes(w *Writer) error {
	if err := w.AppendString(field.Name); err != nil {
		return err
	}

	if err := w.AppendInt32(field.TableID); err != nil {
		return err
	}

	if err := w.AppendInt16(field.Number); err != nil {
		return err
	}

	if err := w.AppendInt32(field.ObjectID); err != nil {
		return err
	}

	if err := w.AppendInt16(field.DataTypeSize); err != nil {
		return err
	}

	if err := w.AppendInt32(field.TypeModifier); err != nil {
		return err
	}

	if err := w.AppendInt16(field.FormatCode); err != nil {
		return err
	}

	return nil
}
