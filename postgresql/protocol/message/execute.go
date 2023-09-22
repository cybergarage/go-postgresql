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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// Execute represents an execute message.
type Execute struct {
	*RequestMessage
	PortalName string
	MaxRows    int32
}

// NewExecute returns a new execute message.
func NewExecuteWithReader(reader *MessageReader) (*Execute, error) {
	msg, err := NewRequestMessageWithReader(reader)
	if err != nil {
		return nil, err
	}

	// The name of the portal to execute (an empty string selects the unnamed portal).
	portal, err := reader.ReadString()
	if err != nil {
		return nil, err
	}

	// Maximum number of rows to return, if portal contains a query that returns rows (ignored otherwise). Zero denotes “no limit”.
	maxRows, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	return &Execute{
		RequestMessage: msg,
		PortalName:     portal,
		MaxRows:        maxRows,
	}, nil
}
