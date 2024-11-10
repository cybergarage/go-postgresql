// Copyright (C) 2020 The go-postgresql Authors. All rights reserved.
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
	"crypto/x509"
	"os"
)

// tlsConfig represents a TLS configuration.
type tlsConfig struct {
	ClientAuthType tls.ClientAuthType
	ServerCert     []byte
	ServerKey      []byte
	RootCerts      [][]byte
	enabled        bool
	tlsConfig      *tls.Config
}

// NewTLSConf returns a new TLS configuration.
func NewTLSConf() *tlsConfig {
	return &tlsConfig{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		ServerCert:     []byte{},
		ServerKey:      []byte{},
		RootCerts:      [][]byte{},
		tlsConfig:      nil,
		enabled:        false,
	}
}

// SetTLSEnabled sets a TLS enabled flag.
func (config *tlsConfig) SetTLSEnabled(enabled bool) {
	config.enabled = enabled
}

// IsEnabled returns true if the TLS is enabled.
func (config *tlsConfig) IsTLSEnabled() bool {
	return config.enabled
}

// SetClientAuthType sets a client authentication type.
func (config *tlsConfig) SetClientAuthType(authType tls.ClientAuthType) {
	config.ClientAuthType = authType
	config.tlsConfig = nil
	config.SetTLSEnabled(true)
}

// SetServerKeyFile sets a SSL server key file.
func (config *tlsConfig) SetServerKeyFile(file string) error {
	key, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	config.SetServerKey(key)
	return nil
}

// SetServerCertFile sets a SSL server certificate file.
func (config *tlsConfig) SetServerCertFile(file string) error {
	cert, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	config.SetServerCert(cert)
	return nil
}

// SetRootCertFile sets a SSL root certificates.
func (config *tlsConfig) SetRootCertFiles(files ...string) error {
	certs := make([][]byte, len(files))
	for n, file := range files {
		cert, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		certs[n] = cert
	}
	config.SetRootCerts(certs...)
	return nil
}

// SetServerKey sets a SSL server key.
func (config *tlsConfig) SetServerKey(key []byte) {
	config.ServerKey = key
	config.tlsConfig = nil
	config.SetTLSEnabled(true)
}

// SetServerCert sets a SSL server certificate.
func (config *tlsConfig) SetServerCert(cert []byte) {
	config.ServerCert = cert
	config.tlsConfig = nil
	config.SetTLSEnabled(true)
}

// SetRootCerts sets a SSL root certificates.
func (config *tlsConfig) SetRootCerts(certs ...[]byte) {
	config.RootCerts = certs
	config.tlsConfig = nil
	config.SetTLSEnabled(true)
}

// SetTLSConfig sets a TLS configuration.
func (config *tlsConfig) SetTLSConfig(tlsConfig *tls.Config) {
	config.tlsConfig = tlsConfig
	config.SetTLSEnabled(true)
}

// TLSConfig returns a TLS configuration from the configuration.
func (config *tlsConfig) TLSConfig() (*tls.Config, error) {
	if !config.IsTLSEnabled() {
		return nil, nil
	}
	if config.tlsConfig != nil {
		return config.tlsConfig, nil
	}
	serverCert, err := tls.X509KeyPair(config.ServerCert, config.ServerKey)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	for _, rootCert := range config.RootCerts {
		certPool.AppendCertsFromPEM(rootCert)
	}
	config.tlsConfig = &tls.Config{ // nolint: exhaustruct
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		RootCAs:      certPool,
		ClientAuth:   config.ClientAuthType,
	}
	return config.tlsConfig, nil
}
