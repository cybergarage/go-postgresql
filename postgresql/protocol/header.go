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
	"encoding/hex"
	"fmt"
)

const (
	// HeaderSize is the static header size of PostgreSQL packet.
	HeaderSize = (1 + 4)
)

// Header represents a standard header of PostgreSQL packet.
type Header struct {
	msgType   MessageType
	msgLength uint32
}

// NewHeader returns a new header instance.
func NewHeader() *Header {
	header := &Header{
		msgType:   0,
		msgLength: 0,
	}
	return header
}

// NewHeaderWithBytes returns a new header instance of the specified bytes.
func NewHeaderWithBytes(msg []byte) (*Header, error) {
	header := NewHeader()
	return header, header.ParseBytes(msg)
}

// Type returns the message type.
func (header *Header) Type() MessageType {
	return header.msgType
}

// Length returns the message length .
func (header *Header) Length() uint32 {
	return header.msgLength
}

// ParseBytes parses the specified bytes.
func (header *Header) ParseBytes(frame []byte) error {
	if len(frame) < HeaderSize {
		return fmt.Errorf(errorShortHeaderLength, len(frame), hex.EncodeToString(frame))
	}

	header.msgType = MessageType(frame[0])
	header.msgLength = uint32(frame[1])<<24 | uint32(frame[2])<<16 | uint32(frame[3])<<8 | uint32(frame[4])

	return nil
}
