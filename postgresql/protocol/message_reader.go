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
	"net"
	"time"
)

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/16/protocol-overview.html

// MessageReader represents a message reader.
type MessageReader struct {
	*Reader
	Type   Type
	Length int32
}

// MessageReaderOption is a function that modifies the MessageReader.
type MessageReaderOption func(*MessageReader)

// WithMessageReadeConn sets the connection for the MessageReader.
func WithMessageReadeConn(conn net.Conn) MessageReaderOption {
	return func(reader *MessageReader) {
		reader.Reader = NewReaderWith(WithReaderConn(conn))
	}
}

// WithMessageReadeBytes sets the bytes for the MessageReader.
func WithMessageReadeBytes(b []byte) MessageReaderOption {
	return func(reader *MessageReader) {
		reader.Reader = NewReaderWith(WithReaderBytes(b))
	}
}

// NewMessageReader returns a new message reader.
func NewMessageReaderWith(opts ...MessageReaderOption) *MessageReader {
	reader := &MessageReader{
		Reader: nil,
		Type:   0,
		Length: 0,
	}
	for _, opt := range opts {
		opt(reader)
	}
	return reader
}

// PeekType peeks a message type.
func (reader *MessageReader) PeekType() (Type, error) {
	bytes, err := reader.Reader.PeekBytes(1)
	if err != nil {
		return 0, err
	}
	return Type(bytes[0]), nil
}

// PeekTypeNonBlocking peeks a message type without blocking.
func (reader *MessageReader) PeekTypeNonBlocking() (Type, error) {
	err := reader.Reader.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	if err != nil {
		return 0, err
	}
	defer func() {
		reader.Reader.conn.SetReadDeadline(time.Time{})
	}()
	reqType, err := reader.PeekType()
	if err != nil {
		return 0, err
	}
	return reqType, nil
}

// IsPeekType returns true whether the peeked message type is the specified type.
func (reader *MessageReader) IsPeekType(t Type) (bool, error) {
	pt, err := reader.PeekType()
	if err != nil {
		return false, err
	}
	return pt == t, nil
}

// ReadType reads a message type.
func (reader *MessageReader) ReadType() (Type, error) {
	t, err := reader.Reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return Type(t), nil
}

// ReadLength reads a message length.
func (reader *MessageReader) ReadLength() (int32, error) {
	l, err := reader.Reader.ReadInt32()
	if err != nil {
		return 0, err
	}
	return l, nil
}
