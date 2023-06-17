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

// NewAuthenticationOk returns a new AuthenticationOk message.
func NewAuthenticationOk() (*ResponseMessage, error) {
	msg := NewResponseMessageWith(AuthenticationOkMessage)
	err := msg.AppendInt32(8)
	if err != nil {
		return nil, err
	}
	return msg, msg.AppendInt32(0)
}
