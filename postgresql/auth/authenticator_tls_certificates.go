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

// TLSCertificatedAuthenticator represents an authenticator for the cleartext password.
type TLSCertificatedAuthenticator struct {
	username string
	password string
}

// NewTLSCertificatedAuthenticator returns a new authenticator with the specified username and password.
func NewTLSCertificatedAuthenticatorWith(username string, password string) *TLSCertificatedAuthenticator {
	authenticator := &TLSCertificatedAuthenticator{
		username: username,
		password: password,
	}
	return authenticator
}

// Authenticate authenticates the specified connection.
func (authenticator *TLSCertificatedAuthenticator) Authenticate(conn Conn) (bool, error) {
	_, ok := conn.TLSConnectionState()
	if !ok {
		return false, nil
	}
	return true, nil
}
