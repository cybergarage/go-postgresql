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

package protocol

const (
	defaultAddr = ""
	defaultPort = 5432
)

// config stores server configuration parammeters.
type config struct {
	productName    string
	productVersion string
	addr           string
	port           int
	*tlsConfig
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() Config {
	config := &config{
		productName:    "",
		productVersion: "",
		addr:           defaultAddr,
		port:           defaultPort,
		tlsConfig:      NewTLSConfig(),
	}
	return config
}

// SetAddress sets a listen address to the configuration.
func (config *config) SetAddress(addr string) {
	config.addr = addr
}

// SetPort sets a listen port to the configuration.
func (config *config) SetPort(port int) {
	config.port = port
}

// Address returns a listen address from the configuration.
func (config *config) Address() string {
	return config.addr
}

// Port returns a listen port from the configuration.
func (config *config) Port() int {
	return config.port
}

// SetProuctName sets a product name to the configuration.
func (config *config) SetProductName(v string) {
	config.productName = v
}

// SetProductVersion sets a product version to the configuration.
func (config *config) SetProductVersion(v string) {
	config.productVersion = v
}

// ProductName returns the product name from the configuration.
func (config *config) ProductName() string {
	return config.productName
}

// ProductVersion returns the product version from the configuration.
func (config *config) ProductVersion() string {
	return config.productVersion
}
