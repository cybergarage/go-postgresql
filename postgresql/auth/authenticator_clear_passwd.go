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

package auth

import (
	"github.com/cybergarage/go-authenticator/auth"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
)

// ClearTextPasswordAuthenticator represents an authenticator for the cleartext password.
type ClearTextPasswordAuthenticator struct {
	auth.CredentialStore
	username string
	password string
}

// NewClearTextPasswordAuthenticator returns a new authenticator with the specified username and password.
func NewClearTextPasswordAuthenticator(username string, password string) *ClearTextPasswordAuthenticator {
	authenticator := &ClearTextPasswordAuthenticator{
		CredentialStore: nil,
		username:        username,
		password:        password,
	}
	return authenticator
}

// Authenticate authenticates the specified connection.
func (authenticator *ClearTextPasswordAuthenticator) Authenticate(conn Conn) (bool, error) {
	startupMessage, ok := conn.StartupMessage()
	if !ok {
		return false, nil
	}

	clientUsername, ok := startupMessage.User()
	if !ok {
		return false, nil
	}
	if clientUsername != authenticator.username {
		return false, nil
	}
	authMsg, err := protocol.NewAuthenticationCleartextPassword()
	if err != nil {
		return false, err
	}
	err = conn.ResponseMessage(authMsg)
	if err != nil {
		return false, err
	}
	msg, err := protocol.NewPasswordWithReader(conn.MessageReader())
	if err != nil {
		return false, err
	}
	if msg.Password != authenticator.password {
		return false, nil
	}
	return true, nil
}
