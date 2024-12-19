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

// CertAuthenticator represents an authenticator for TLS certificates.
type CertAuthenticator struct {
	Authenticator auth.CertificateAuthenticator
	commonName    string
}

// CertAuthenticatorOption represents an authenticator option.
type CertAuthenticatorOption = func(*CertAuthenticator)

// NewCertificateAuthenticator returns a new certificate authenticator.
func NewCertificateAuthenticator(opts ...CertAuthenticatorOption) *CertAuthenticator {
	authenticator := &CertAuthenticator{
		Authenticator: nil,
		commonName:    "",
	}
	for _, opt := range opts {
		opt(authenticator)
	}

	return authenticator
}

// WithCommonName returns an authenticator option to set the common name.
func WithCommonName(name string) func(*CertAuthenticator) {
	return func(conn *CertAuthenticator) {
		conn.commonName = name
	}
}

// Authenticate authenticates the specified connection.
func (authenticator *CertAuthenticator) Authenticate(conn Conn) (bool, error) {
	if authenticator.Authenticator == nil {
		return true, nil
	}
	if !conn.IsTLSConnection() {
		return false, nil
	}
	return authenticator.Authenticator.VerifyCertificate(conn.TLSConn())
}
