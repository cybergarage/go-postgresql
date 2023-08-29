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
	"bufio"

	util "github.com/cybergarage/go-postgresql/postgresql/encoding/bytes"
)

// Reader represents a message reader.
type Reader struct {
	*bufio.Reader
}

// NewReader returns a new message reader.
func NewReaderWith(reader *bufio.Reader) *Reader {
	return &Reader{
		Reader: reader,
	}
}

// PeekInt32 reads a 32-bit integer.
func (reader *Reader) PeekInt32() (int32, error) {
	int32Bytes, err := reader.Peek(4)
	if err != nil {
		return 0, err
	}
	return util.BytesToInt32(int32Bytes), nil
}

// ReadInt32 reads a 32-bit integer.
func (reader *Reader) ReadInt32() (int32, error) {
	int32Bytes := make([]byte, 4)
	nRead, err := reader.Read(int32Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 4 {
		return 0, newShortMessageError(4, nRead)
	}
	return util.BytesToInt32(int32Bytes), nil
}

// ReadInt16 reads a 16-bit integer.
func (reader *Reader) ReadInt16() (int16, error) {
	int16Bytes := make([]byte, 2)
	nRead, err := reader.Read(int16Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 2 {
		return 0, newShortMessageError(2, nRead)
	}
	return util.BytesToInt16(int16Bytes), nil
}

// ReadByte reads a byte array data into the specified buffer.
func (reader *Reader) ReadBytes(buf []byte) (int, error) {
	return reader.Read(buf)
}

// ReadString reads a string.
func (reader *Reader) ReadString() (string, error) {
	b, err := reader.Reader.ReadBytes(0x00)
	if err != nil {
		return "", err
	}
	return string(b[:len(b)-1]), nil
}
