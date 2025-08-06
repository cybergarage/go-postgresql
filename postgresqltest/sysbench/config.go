// Copyright (C) 2025 The go-postgresql Authors. All rights reserved.
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

package sysbench

import (
	"strconv"

	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-sqltest/sqltest/sysbench"
)

const (
	// https://github.com/akopytov/sysbench
	PgSQLHost     = "pgsql-host"
	PgSQLPort     = "pgsql-port"
	PgSQLUser     = "pgsql-user"
	PgSQLPassword = "pgsql-password"
	PgSQLDB       = "pgsql-db"
)

// Config represents a sysbench config.
type Config struct {
	*sysbench.Config
}

// NewDefaultConfig returns a new default config.
func NewDefaultConfig() *Config {
	cfg := &Config{
		Config: sysbench.NewDefaultConfig(),
	}
	cfg.SetDBDriver("pgsql")
	cfg.SetHost("127.0.0.1")
	cfg.SetPort(postgresql.DefaultPort)
	cfg.SetUser(sysbench.User())
	cfg.SetPassword(sysbench.Password())

	return cfg
}

// SetHost sets the host.
func (cfg *Config) SetHost(host string) {
	cfg.Set(PgSQLHost, host)
}

// SetPort sets the port.
func (cfg *Config) SetPort(port int) {
	cfg.Set(PgSQLPort, strconv.Itoa(port))
}

// SetUser sets the user.
func (cfg *Config) SetUser(user string) {
	cfg.Set(PgSQLUser, user)
}

// SetPassword sets the password.
func (cfg *Config) SetPassword(password string) {
	cfg.Set(PgSQLPassword, password)
}

// SetDB sets the db.
func (cfg *Config) SetDB(db string) {
	cfg.Set(PgSQLDB, db)
}
