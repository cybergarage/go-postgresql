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

const (
	defaultHost = "127.0.0.1"
	defaultPort = DefaultPort
)

// Config stores server configuration parammeters.
type Config struct {
	host string
	port int
	db   string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		host: defaultHost,
		port: defaultPort,
	}
	return config
}

// SetHost sets a host address to the configuration.
func (config *Config) SetHost(host string) {
	config.host = host
}

// SetPort sets a listen port to the configuration.
func (config *Config) SetPort(port int) {
	config.port = port
}

// SetDatabase sets a database to the configuration.
func (config *Config) SetDatabase(db string) {
	config.db = db
}

// Host returns a host address from the configuration.
func (config *Config) Host() string {
	return config.host
}

// Port returns a port from the configuration.
func (config *Config) Port() int {
	return config.port
}

// Database returns a database from the configuration.
func (config *Config) Database() string {
	return config.db
}
