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

package auth

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// CleartextPasswordAuthenticator represents an authenticator for the cleartext password.
type CleartextPasswordAuthenticator struct {
	username string
	password string
}

// NewCleartextPasswordAuthenticator returns a new authenticator with the specified username and password.
func NewCleartextPasswordAuthenticatorWith(username string, password string) *CleartextPasswordAuthenticator {
	authenticator := &CleartextPasswordAuthenticator{
		username: username,
		password: password,
	}
	return authenticator
}

func (authenticator *CleartextPasswordAuthenticator) Authenticate(conn Conn, startupMessage *message.Startup) (bool, error) {
	return true, nil
}