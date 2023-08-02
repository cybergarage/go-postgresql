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

// Message represents a message of PostgreSQL packet.
type Message struct {
	*MessageReader
	Type   Type
	Length int32
}

// NewMessageWithReader returns a new message with the specified reader.
func NewMessageWithReader(reader *MessageReader) (*Message, error) {
	t, err := reader.ReadType()
	if err != nil {
		return nil, err
	}
	l, err := reader.ReadLength()
	if err != nil {
		return nil, err
	}
	return &Message{
		MessageReader: reader,
		Type:          t,
		Length:        l,
	}, nil
}

// MessageType returns a message type.
func (msg *Message) MessageType() Type {
	return msg.Type
}

// MessageLength returns a message length.
func (msg *Message) MessageLength() int32 {
	return msg.Length
}