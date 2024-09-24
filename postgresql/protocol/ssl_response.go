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
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

const (
	SSLEnabled  = 'S'
	SSLDisabled = 'N'
)

// SSLResponse represents a SSLResponse protocol.
type SSLResponse struct {
	Enabled byte
}

// NewSSLResponse returns a new SSLResponse protocol.
func NewSSLResponseWith(enabled rune) *SSLResponse {
	return &SSLResponse{
		Enabled: byte(enabled),
	}
}

// Type returns the message type.
func (msg *SSLResponse) Type() Type {
	return SSLResponseMessage
}

// Bytes returns a byte array of the protocol.
func (msg *SSLResponse) Bytes() ([]byte, error) {
	return []byte{msg.Enabled}, nil
}
