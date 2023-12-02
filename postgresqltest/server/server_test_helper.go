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

package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/auth"
	"github.com/cybergarage/go-postgresql/postgresqltest/client"
	"github.com/jackc/pgx/v5"
)

const testDBNamePrefix = "pgtest"

type ServerTestFunc = func(*testing.T, *Server, string)

func RunServerTests(t *testing.T, server *Server) {
	t.Helper()

	log.SetStdoutDebugEnbled(true)

	testFuncs := []struct {
		name string
		fn   ServerTestFunc
	}{
		{"authenticator", RunAuthenticatorTest},
		//		{"copy", TestServerCopy},
	}

	for _, testFunc := range testFuncs {
		testDBName := fmt.Sprintf("%s%d", testDBNamePrefix, time.Now().UnixNano())
		t.Run(testFunc.name, func(t *testing.T) {
			// Create a test database

			client := client.NewDefaultClient()

			client.SetDatabase("postgres")
			err := client.Open()
			if err != nil {
				t.Error(err)
				return
			}

			err = client.CreateDatabase(testDBName)
			if err != nil {
				t.Error(err)
				return
			}

			err = client.Close()
			if err != nil {
				t.Error(err)
			}

			// Run tests

			testFunc.fn(t, server, testDBName)

			err = client.DropDatabase(testDBName)
			if err != nil {
				t.Error(err)
			}

			err = client.Close()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// RunAuthenticatorTest tests the authenticator.
func RunAuthenticatorTest(t *testing.T, server *Server, testDBName string) {
	t.Helper()

	username := "testuser"
	password := "testpassword"

	authenticators := []auth.Authenticator{
		auth.NewCleartextPasswordAuthenticatorWith(username, password),
	}

	for _, authenticator := range authenticators {
		server.AddAuthenticator(authenticator)

		client := client.NewDefaultClient()
		client.SetUser(username)
		client.SetPassword(password)
		client.SetDatabase(testDBName)
		err := client.Open()
		if err != nil {
			t.Error(err)
			return
		}

		err = client.Ping()
		if err != nil {
			t.Error(err)
		}

		err = client.Close()
		if err != nil {
			t.Error(err)
		}

		server.ClearAuthenticators()
	}
}

// RunServerCopyTest tests the COPY command.
func RunServerCopyTest(t *testing.T, server *Server, testDBName string) {
	t.Helper()

	// Run tests

	client := client.NewPgxClient()
	client.SetDatabase(testDBName)
	err := client.Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err := client.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	defer func() {
		err := client.DropDatabase(testDBName)
		if err != nil {
			t.Error(err)
		}
	}()

	rows, err := client.Query("CREATE TABLE cptest (ctext TEXT PRIMARY KEY, cint INT, cfloat FLOAT);")
	if err != nil {
		t.Error(err)
		return
	}

	if rows.Err() != nil {
		t.Error(rows.Err())
		rows.Close()
		return
	}
	rows.Close()

	conn := client.Conn()

	copyRows := [][]any{
		{"text1", 1, 1.1},
		{"text2", 2, 2.2},
		{"text3", 3, 3.3},
	}

	copyCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"cptest"},
		[]string{"ctext", "cint", "cfloat"},
		pgx.CopyFromRows(copyRows),
	)

	if err != nil {
		t.Error(err)
		return
	}

	if copyCount != int64(len(copyRows)) {
		t.Errorf("copyCount (%d) != len(rows) (%d)", copyCount, len(copyRows))
		return
	}
}
