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

// TLSConf represents a TLS configuration.
type TLSConf struct {
	ClientAuthType tls.ClientAuthType
	ServerCertFile string
	ServerKeyFile  string
	RootCertFiles  []string
	tlsConfig      *tls.Config
}

// NewTLSConf returns a new TLS configuration.
func NewTLSConf() *TLSConf {
	return &TLSConf{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		ServerCertFile: "",
		ServerKeyFile:  "",
		RootCertFiles:  []string{},
		tlsConfig:      nil,
	}
}

// SetClientAuthType sets a client authentication type.
func (config *TLSConf) SetClientAuthType(authType tls.ClientAuthType) {
	config.ClientAuthType = authType
}

// SetServerKeyFile sets a SSL server key file.
func (config *TLSConf) SetServerKeyFile(file string) {
	config.ServerKeyFile = file
}

// SetServerCertFile sets a SSL server certificate file.
func (config *TLSConf) SetServerCertFile(file string) {
	config.ServerCertFile = file
}

// SetRootCertFile sets a SSL root certificates.
func (config *TLSConf) SetRootCertFiles(files ...string) {
	config.RootCertFiles = files
}

// TLSConfig returns a TLS configuration from the configuration.
func (config *TLSConf) TLSConfig() (*tls.Config, error) {
	if config.tlsConfig != nil {
		return config.tlsConfig, nil
	}
	if len(config.ServerCertFile) == 0 || len(config.ServerKeyFile) == 0 {
		return nil, nil
	}
	serverCert, err := tls.LoadX509KeyPair(config.ServerCertFile, config.ServerKeyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	for _, rootCertFile := range config.RootCertFiles {
		rootCert, err := os.ReadFile(rootCertFile)
		if err != nil {
			return nil, err
		}
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
