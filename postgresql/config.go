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

import (
	"crypto/tls"
)

// TLSConfig represents a TLS configuration.
type TLSConfig interface {
	// SetTLSEnabled sets a TLS enabled flag.
	SetTLSEnabled(enabled bool)
	// IsEnabled returns true if the TLS is enabled.
	IsTLSEnabled() bool
	// SetClientAuthType sets a client authentication type.
	SetClientAuthType(authType tls.ClientAuthType)
	// SetServerKeyFile sets a SSL server key file.
	SetServerKeyFile(file string) error
	// SetServerCertFile sets a SSL server certificate file.
	SetServerCertFile(file string) error
	// SetRootCertFile sets a SSL root certificates.
	SetRootCertFiles(files ...string) error
	// SetServerKey sets a SSL server key.
	SetServerKey(key []byte)
	// SetServerCert sets a SSL server certificate.
	SetServerCert(cert []byte)
	// SetRootCerts sets a SSL root certificates.
	SetRootCerts(certs ...[]byte)
	// SetTLSConfig sets a TLS configuration.
	SetTLSConfig(tlsConfig *tls.Config)
	// TLSConfig returns a TLS configuration from the configuration.
	TLSConfig() (*tls.Config, error)
}

// Config represents a server configuration.
type Config interface {
	TLSConfig
	// SetAddress sets a listen address to the configuration.
	SetAddress(addr string)
	// SetPort sets a listen port to the configuration.
	SetPort(port int)
	// Address returns a listen address from the configuration.
	Address() string
	// Port returns a listen port from the configuration.
	Port() int
}
