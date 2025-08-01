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

package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/auth"
	pgx "github.com/jackc/pgx/v5"
)

const testDBNamePrefix = "pgtest"

type ServerTestFunc = func(*testing.T, *Server, string)

func RunServerTests(t *testing.T, server *Server) {
	t.Helper()

	log.EnableStdoutDebug(true)

	testFuncs := []struct {
		name string
		fn   ServerTestFunc
	}{
		{"authenticator", RunPasswordAuthenticatorTest},
		// {"tls", RunCertificateAuthenticatorTest},
		// {"copy", TestServerCopy},
	}

	for _, testFunc := range testFuncs {
		testDBName := fmt.Sprintf("%s%d", testDBNamePrefix, time.Now().UnixNano())
		t.Run(testFunc.name, func(t *testing.T) {
			// Create a test database

			client := postgresql.NewDefaultClient()

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

			err = client.Open()
			if err != nil {
				t.Error(err)
				return
			}

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

// RunPasswordAuthenticatorTest tests the authenticators.
func RunPasswordAuthenticatorTest(t *testing.T, server *Server, testDBName string) {
	t.Helper()

	const (
		username = "testuser"
		password = "testpassword"
	)

	cred := auth.NewCredential(
		auth.WithCredentialUsername(username),
		auth.WithCredentialPassword(password),
	)
	server.SetCredential(cred)
	server.SetCredentialStore(server)
	defer func() {
		server.SetCredentialStore(nil)
	}()

	client := postgresql.NewDefaultClient()
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
}

// RunCertificateAuthenticatorTest tests the TLS session.
// PostgreSQL: Documentation: 16: 34.19. SSL Support
// https://www.postgresql.org/docs/current/libpq-ssl.html
// PostgreSQL: Documentation: 16: 19.9. Secure TCP/IP Connections with SSL
// https://www.postgresql.org/docs/current/ssl-tcp.html#SSL-CERTIFICATE-CREATION
func RunCertificateAuthenticatorTest(t *testing.T, server *Server, testDBName string) {
	t.Helper()

	const (
		clientKey  = "../certs/client-key.pem"
		clientCert = "../certs/client-cert.pem"
		rootCert   = "../certs/cacert.pem"
	)

	auth.NewCertificateAuthenticator(auth.WithCommonNameRegexp("localhost"))

	client := postgresql.NewDefaultClient()
	client.SetClientKeyFile(clientKey)
	client.SetClientCertFile(clientCert)
	client.SetRootCertFile(rootCert)

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
}

// RunServerCopyTest tests the COPY command.
func RunServerCopyTest(t *testing.T, server *Server, testDBName string) {
	t.Helper()

	// Run tests

	client := postgresql.NewPgxClient()
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
		t.Context(),
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
