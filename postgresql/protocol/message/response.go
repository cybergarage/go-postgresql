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

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/16/protocol-overview.html

// Response represents a backend response.
type Response struct {
	Type Type
	*Writer
}

// NewResponse returns a new request message with the specified reader.
func NewResponse() *Response {
	return &Response{
		Type:   NoneMessage,
		Writer: NewWriter(),
	}
}

// SetType sets a message type.
func (msg *Response) SetType(t Type) {
	msg.Type = t
}

// Bytes returns the message bytes.
func (msg *Response) Bytes() ([]byte, error) {
	msgBytes, err := msg.Writer.Bytes()
	if err != nil {
		return nil, err
	}
	l := len(msgBytes)
	b := make([]byte, 0, l+1+4)
	if msg.Type != NoneMessage {
		b = append(b, byte(msg.Type))
	}
	b = append(b, Int32ToBytes(int32(l))...)
	return append(b, msgBytes...), nil
}
