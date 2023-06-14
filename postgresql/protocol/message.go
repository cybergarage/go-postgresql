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

	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/16/protocol-overview.html

// Message represents an operation message.
type Message struct {
	*message.Reader
	Type   MessageType
	Length int
}

// NewMessageWith returns a new Message with the specified reader.
func NewMessageWith(reader *bufio.Reader) *Message {
	return &Message{
		Reader: message.NewReaderWith(reader),
		Type:   0,
		Length: 0,
	}
}

// ReadType reads a message type.
func (msg *Message) ReadType() (MessageType, error) {
	t, err := msg.ReadByte()
	if err != nil {
		return 0, err
	}
	msg.Type = MessageType(t)
	return msg.Type, nil
}

// ReadLength reads a message length.
func (msg *Message) ReadLength() (int, error) {
	l, err := msg.ReadInt32()
	if err != nil {
		return 0, err
	}
	msg.Length = l
	return msg.Length, nil
}

// ParseStartupMessage parses a startup message.
func (msg *Message) ParseStartupMessage() (*message.Startup, error) {
	return message.NewStartupWith(msg.Reader)
}
