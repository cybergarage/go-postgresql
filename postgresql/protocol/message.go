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

package protocol

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/9.3/protocol-overview.html

// MessageType represents a message type.
type MessageType byte

// Message represents an operation message.
type Message interface {
	// Type returns the message type.
	Type() MessageType
	// Length returns the payload size.
	Length() uint32
}

// NewMessageWithBytes returns a parsed message of the specified bytes.
func NewMessageWithBytes(msg []byte) (Message, error) {
	header, err := NewHeaderWithBytes(msg)
	if err != nil {
		return nil, err
	}
	return header, nil
}
