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

// Message represents a message of PostgreSQL packet.
// See : PostgreSQL Packets
// https://www.postgresql.org/docs/16/protocol-overview.html

// RequestMessage represents a frontend request.
type RequestMessage struct {
	*Message
}

// NewRequestMessageWithReader returns a new request message with the specified reader.
func NewRequestMessageWithReader(reader *MessageReader) (*RequestMessage, error) {
	msg, err := NewMessageWithReader(reader)
	if err != nil {
		return nil, err
	}
	return &RequestMessage{
		Message: msg,
	}, nil
}
