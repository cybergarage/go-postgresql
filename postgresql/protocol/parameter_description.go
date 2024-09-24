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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// ParameterDescription represents a parameter description response protocol.
type ParameterDescription struct {
	*ResponseMessage
}

// NewParameterDescription returns a parameter description response instance.
func NewParameterDescription() *ParameterDescription {
	return &ParameterDescription{
		ResponseMessage: NewResponseMessageWith(ParameterDescriptionMessage),
	}
}

// NewParameterDescriptionWith returns a parameter description response instance with the specified parameters.
func NewParameterDescriptionWith(objectIDs ...ObjectID) (*ParameterDescription, error) {
	msg := NewParameterDescription()
	err := msg.AppendInt16(int16(len(objectIDs)))
	if err != nil {
		return nil, err
	}
	for _, objectID := range objectIDs {
		err := msg.AppendInt32(objectID)
		if err != nil {
			return nil, err
		}
	}
	return msg, nil
}
