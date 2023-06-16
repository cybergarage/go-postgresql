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

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/16/protocol-overview.html

// ResponseMessage represents a backend response message interface.
type ResponseMessage interface {
	// Type returns the message type.
	Type() Type
	// Bytes returns the message bytes.
	Bytes() ([]byte, error)
}

// Response represents a backend response instance.
type Response struct {
	typ Type
	*Writer
}

// NewResponse returns a new request message instance.
func NewResponse() *Response {
	return NewResponseWith(NoneMessage)
}

// NewResponseWith returns a new response message with the specified message type.
func NewResponseWith(t Type) *Response {
	return &Response{
		typ:    t,
		Writer: NewWriter(),
	}
}

// SetType sets a message type.
func (msg *Response) SetType(t Type) {
	msg.typ = t
}

// Type returns the message type.
func (msg *Response) Type() Type {
	return msg.typ
}

// Bytes returns the message bytes.
func (msg *Response) Bytes() ([]byte, error) {
	msgBytes, err := msg.Writer.Bytes()
	if err != nil {
		return nil, err
	}
	l := len(msgBytes)
	b := make([]byte, 0, l+1+4)
	if msg.typ != NoneMessage {
		b = append(b, byte(msg.typ))
	}
	b = append(b, Int32ToBytes(int32(l))...)
	return append(b, msgBytes...), nil
}

// StringResponse represents a backend string response instance.
type StringResponse struct {
	*Response
}

// NewStringResponseWith returns a new string response message with the specified message type.
func NewStringResponseWith(t Type) *StringResponse {
	return &StringResponse{
		Response: NewResponseWith(t),
	}
}

// Bytes returns the message bytes after adding a null terminator.
func (msg *StringResponse) Bytes() ([]byte, error) {
	if err := msg.AppendTerminator(); err != nil {
		return nil, err
	}
	return msg.Response.Bytes()
}
