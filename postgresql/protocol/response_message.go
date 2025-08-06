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

// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

import (
	util "github.com/cybergarage/go-postgresql/postgresql/encoding/bytes"
)

// ResponseMessage represents a backend response instance.
type ResponseMessage struct {
	typ Type
	*Writer
}

// NewResponseMessage returns a new request message instance.
func NewResponseMessage() *ResponseMessage {
	return NewResponseMessageWith(NoneMessage)
}

// NewResponseMessageWith returns a new response message with the specified message type.
func NewResponseMessageWith(t Type) *ResponseMessage {
	return &ResponseMessage{
		typ:    t,
		Writer: NewWriter(),
	}
}

// SetType sets a message type.
func (msg *ResponseMessage) SetType(t Type) {
	msg.typ = t
}

// Type returns the message type.
func (msg *ResponseMessage) Type() Type {
	return msg.typ
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *ResponseMessage) Bytes() ([]byte, error) {
	msgBytes, err := msg.Writer.Bytes()
	if err != nil {
		return nil, err
	}

	l := len(msgBytes)

	b := make([]byte, 0, 1+4+l)
	if msg.typ != NoneMessage {
		b = append(b, byte(msg.typ))
	}

	b = append(b, util.Int32ToBytes(int32(l+4))...)

	return append(b, msgBytes...), nil
}
