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
	"net"
	"time"

	util "github.com/cybergarage/go-postgresql/postgresql/encoding/bytes"
)

// Reader represents a message reader.
type Reader struct {
	conn net.Conn
	io.Reader
	peekBuf []byte
}

// ReaderOptionFunc is a function that modifies the Reader.
type ReaderOption func(*Reader)

// WithReaderFile sets the file for the Reader.
func WithReaderConn(conn net.Conn) ReaderOption {
	return func(reader *Reader) {
		reader.conn = conn
		reader.Reader = conn
	}
}

// WithReaderBytes sets the bytes for the Reader.
func WithReaderBytes(b []byte) ReaderOption {
	return func(reader *Reader) {
		reader.conn = nil
		reader.Reader = bytes.NewReader(b)
	}
}

// NewReader returns a new message reader.
func NewReaderWith(opts ...ReaderOption) *Reader {
	reader := &Reader{
		Reader:  nil,
		conn:    nil,
		peekBuf: make([]byte, 0),
	}
	for _, opt := range opts {
		opt(reader)
	}

	return reader
}

// SetReadDeadline sets the read deadline.
func (reader *Reader) SetReadDeadline(t time.Time) error {
	if reader.conn == nil {
		return ErrNotSupported
	}

	return reader.conn.SetReadDeadline(t)
}

// ReadBytes reads a byte array.
func (reader *Reader) ReadBytes(buf []byte) (int, error) {
	nBufSize := len(buf)
	nReadBuf := 0

	if 0 < len(reader.peekBuf) {
		nCopy := copy(buf, reader.peekBuf)
		reader.peekBuf = reader.peekBuf[nCopy:]
		nReadBuf = nCopy
	}

	if nBufSize <= nReadBuf {
		return nReadBuf, nil
	}

	nRead, err := io.ReadAtLeast(reader.Reader, buf[nReadBuf:], nBufSize-nReadBuf)
	if err != nil {
		return nReadBuf, err
	}

	nReadBuf += nRead

	return nReadBuf, err
}

// ReadBytes reads a byte.
func (reader *Reader) ReadByte() (byte, error) {
	b := make([]byte, 1)

	_, err := reader.ReadBytes(b)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

// PeekBytes reads bytes without removing them from the buffer.
func (reader *Reader) PeekBytes(n int) ([]byte, error) {
	buf := make([]byte, n)

	nRead, err := reader.ReadBytes(buf)
	if err != nil {
		return nil, err
	}

	if nRead != n {
		return nil, newShortMessageError(n, nRead)
	}

	reader.peekBuf = append(reader.peekBuf, buf...)

	return buf, nil
}

// PeekInt32 reads a 32-bit integer.
func (reader *Reader) PeekInt32() (int32, error) {
	int32Bytes, err := reader.PeekBytes(4)
	if err != nil {
		return 0, err
	}

	return util.BytesToInt32(int32Bytes), nil
}

// ReadInt32 reads a 32-bit integer.
func (reader *Reader) ReadInt32() (int32, error) {
	int32Bytes := make([]byte, 4)

	nRead, err := reader.ReadBytes(int32Bytes)
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

	nRead, err := reader.ReadBytes(int16Bytes)
	if err != nil {
		return 0, err
	}

	if nRead != 2 {
		return 0, newShortMessageError(2, nRead)
	}

	return util.BytesToInt16(int16Bytes), nil
}

// ReadBytesUntil reads bytes until the specified delimiter.
func (reader *Reader) ReadBytesUntil(delim byte) ([]byte, error) {
	buf := make([]byte, 0)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		buf = append(buf, b)
		if b == delim {
			break
		}
	}

	return buf, nil
}

// ReadString reads a string.
func (reader *Reader) ReadString() (string, error) {
	strBytes, err := reader.ReadBytesUntil(0x00)
	if err != nil {
		return "", err
	}

	return string(strBytes[:len(strBytes)-1]), nil
}
