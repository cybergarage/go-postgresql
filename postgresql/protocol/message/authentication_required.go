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

// NewAuthenticationCleartextPassword returns a new AuthenticationCleartextPassword message.
func NewAuthenticationCleartextPassword() (*ResponseMessage, error) {
	msg := NewResponseMessageWith(AuthenticationCleartextPasswordMessage)
	return msg, msg.AppendInt32(AuthenticationCleartextPasswordRequired)
}

// NewAuthenticationMD5Password returns a new AuthenticationMD5Password message.
func NewAuthenticationMD5Password(salt []byte) (*ResponseMessage, error) {
	msg := NewResponseMessageWith(AuthenticationMD5PasswordMessage)
	msg.AppendInt32(AuthenticationMD5PasswordRequired)
	msg.AppendBytes(salt)
	return msg, nil
}
