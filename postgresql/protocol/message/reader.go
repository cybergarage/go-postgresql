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

import (
	"bufio"
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

// ReadInt32 reads a 32-bit integer.
func (reader *Reader) ReadInt32() (int, error) {
	int32Bytes := make([]byte, 4)
	nRead, err := reader.Read(int32Bytes)
	if err != nil {
		return 0, err
	}
	if nRead != 4 {
		return 0, newShortMessageErrorWith(4, nRead)
	}
	return Int32BytesToInt(int32Bytes), nil
}

// ReadString reads a string.
func (reader *Reader) ReadString() (string, error) {
	return reader.Reader.ReadString(0x00)
}
