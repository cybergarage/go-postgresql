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

package postgresql

type PasswordAuthenticator interface {
	// AuthenticateUser authenticates the user with the given credentials.
	AuthenticateUser(conn *Conn, username string, password string) (bool, error)
}

// AuthServer represents a authenticator.
type AuthServer struct {
	passwdAuthHandler PasswordAuthenticator
}

// NewAuthenticator returns a new authenticator.
func NewAuthenticator() *AuthServer {
	authenticator := &AuthServer{
		passwdAuthHandler: nil,
	}
	return authenticator
}

// SetPasswordAuthenticator sets a new password authenticator.
func (server *AuthServer) SetPasswordAuthenticator(handler PasswordAuthenticator) {
	server.passwdAuthHandler = handler
}

// PasswordAuthenticator returns the password authenticator.
func (server *AuthServer) PasswordAuthenticator() PasswordAuthenticator {
	return server.passwdAuthHandler
}
