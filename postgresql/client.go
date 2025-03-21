// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgresql

import (
	"github.com/cybergarage/go-sqltest/sqltest"
)

// Client represents a PostgreSQL client interface.
type Client = sqltest.Client

// PqClient represents a PostgreSQL client.
type PqClient = sqltest.PqClient

// PgxClient represents a PostgreSQL client.
type PgxClient = sqltest.PgxClient

// NewDefaultClient returns a new default PostgreSQL client.
func NewDefaultClient() Client {
	return sqltest.NewPostgresClient()
}

// NewPgClient returns a new pq client.
func NewPqClient() *sqltest.PqClient {
	return sqltest.NewPqClient()
}

// NewPgxClient returns a new Pgx client.
func NewPgxClient() *sqltest.PgxClient {
	return sqltest.NewPgxClient()
}
