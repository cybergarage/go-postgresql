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

// BackendKeyData represents a parameter status response message.
type BackendKeyData struct {
	*ResponseMessage
}

// NewBackendKeyData returns a parameter status response instance.
func NewBackendKeyData() *BackendKeyData {
	return &BackendKeyData{
		ResponseMessage: NewResponseMessageWith(BackendKeyDataMessage),
	}
}

// NewBackendKeyDataWith returns a parameter status response instance with the specified parameters.
func NewBackendKeyDataWith(processID int32, secretKey int32) (*BackendKeyData, error) {
	msg := &BackendKeyData{
		ResponseMessage: NewResponseMessageWith(BackendKeyDataMessage),
	}
	err := msg.AppendInt32(processID)
	if err != nil {
		return nil, err
	}
	err = msg.AppendInt32(secretKey)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
