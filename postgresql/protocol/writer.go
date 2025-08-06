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
	"bufio"
	"bytes"

	encording "github.com/cybergarage/go-postgresql/postgresql/encoding/bytes"
)

// Writer represents a message writer.
type Writer struct {
	*bytes.Buffer
	*bufio.Writer
}

// NewWriter returns a new message writer.
func NewWriter() *Writer {
	buffer := &bytes.Buffer{}

	return &Writer{
		Buffer: buffer,
		Writer: bufio.NewWriter(buffer),
	}
}

// AppendByte appends the specified byte.
func (writer *Writer) AppendByte(c byte) error {
	return writer.Writer.WriteByte(c)
}

// AppendBytes appends the specified bytes.
func (writer *Writer) AppendBytes(p []byte) error {
	_, err := writer.Writer.Write(p)
	return err
}

// AppendInt8 appends the specified int8 value.
func (writer *Writer) AppendInt8(v int8) error {
	return writer.AppendBytes(encording.Int8ToBytes(v))
}

// AppendInt16 appends the specified int16 value.
func (writer *Writer) AppendInt16(v int16) error {
	return writer.AppendBytes(encording.Int16ToBytes(v))
}

// AppendInt32 appends the specified int32 value.
func (writer *Writer) AppendInt32(v int32) error {
	return writer.AppendBytes(encording.Int32ToBytes(v))
}

// AppendInt64 appends the specified int64 value.
func (writer *Writer) AppendInt64(v int64) error {
	return writer.AppendBytes(encording.Int64ToBytes(v))
}

// AppendString appends the specified string.
func (writer *Writer) AppendString(s string) error {
	if 0 < len(s) {
		_, err := writer.Writer.WriteString(s)
		if err != nil {
			return err
		}
	}

	return writer.AppendTerminator()
}

// AppendFloat32 appends the specified float32 value.
func (writer *Writer) AppendFloat32(v float32) error {
	return writer.AppendBytes(encording.Float32ToBytes(v))
}

// AppendFloat64 appends the specified float64 value.
func (writer *Writer) AppendFloat64(v float64) error {
	return writer.AppendBytes(encording.Float64ToBytes(v))
}

// AppendTerminator appends a null terminator.
func (writer *Writer) AppendTerminator() error {
	return writer.Writer.WriteByte(0x00)
}

// Bytes returns the message bytes.
func (writer *Writer) Bytes() ([]byte, error) {
	err := writer.Flush()
	if err != nil {
		return nil, err
	}

	return writer.Buffer.Bytes(), nil
}
