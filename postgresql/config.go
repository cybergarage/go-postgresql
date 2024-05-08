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

package postgresql

import (
	"crypto/tls"
)

const (
	defaultAddr = ""
	defaultPort = DefaultPort
)

// Config stores server configuration parammeters.
type Config struct {
	addr string
	port int
	tls  *tls.Config
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		addr: defaultAddr,
		port: defaultPort,
		tls:  nil,
	}
	return config
}

// SetAddress sets a listen address to the configuration.
func (config *Config) SetAddress(addr string) {
	config.addr = addr
}

// SetPort sets a listen port to the configuration.
func (config *Config) SetPort(port int) {
	config.port = port
}

// SetTLSConfig sets a TLS configuration to the configuration.
func (config *Config) SetTLSConfig(conf *tls.Config) {
	config.tls = conf
}

// Address returns a listen address from the configuration.
func (config *Config) Address() string {
	return config.addr
}

// Port returns a listen port from the configuration.
func (config *Config) Port() int {
	return config.port
}

// TLSConfig returns a TLS configuration from the configuration.
func (config *Config) TLSConfig() *tls.Config {
	return config.tls
}
