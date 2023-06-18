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
)

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/16/protocol-overview.html

// RequestMessage represents a frontend request.
type RequestMessage struct {
	*Reader
	Type   Type
	Length int32
}

// NewRequestMessageWith returns a new request message with the specified reader.
func NewRequestMessageWith(reader *bufio.Reader) *RequestMessage {
	return &RequestMessage{
		Reader: NewReaderWith(reader),
		Type:   0,
		Length: 0,
	}
}

// ReadType reads a message type.
func (msg *RequestMessage) ReadType() (Type, error) {
	var err error
	msg.Type, err = msg.Reader.ReadType()
	return msg.Type, err
}

// ReadLength reads a message length.
func (msg *RequestMessage) ReadLength() (int32, error) {
	var err error
	msg.Length, err = msg.Reader.ReadLength()
	return msg.Length, err
}

// ParseStartupMessage parses a startup message.
func (msg *RequestMessage) ParseStartupMessage() (*Startup, error) {
	return NewStartupWith(msg.Reader)
}

// ParseParseMessage parses a parse message.
func (msg *RequestMessage) ParseParseMessage() (*Parse, error) {
	return NewParseWith(msg.Reader)
}

// ParseBindMessage parses a bind message.
func (msg *RequestMessage) ParseBindMessage() (*Bind, error) {
	return NewBindWith(msg.Reader)
}

// ParseQueryMessage parses a bind message.
func (msg *RequestMessage) ParseQueryMessage() (*Query, error) {
	return NewQueryWith(msg.Reader)
}
