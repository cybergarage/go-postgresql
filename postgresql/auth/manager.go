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
)

// AuthManager represent an authenticator manager.
type AuthManager interface {
	// AddAuthenticator adds a new authenticator.
	AddAuthenticator(authenticator Authenticator)
	// ClearAuthenticators clears all authenticators.
	ClearAuthenticators()
	// SetCredentialStore sets the credential store.
	SetCredentialStore(store auth.CredentialStore)
	// SetCertificateAuthenticator sets the certificate authenticator.
	SetCertificateAuthenticator(auth auth.CertificateAuthenticator)
	// Authenticate authenticates the connection with the startup protocol.
	Authenticate(conn Conn) (bool, error)
}
