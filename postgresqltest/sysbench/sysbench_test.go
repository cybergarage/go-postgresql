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
	"os"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresqltest"
	"github.com/cybergarage/go-postgresql/postgresqltest/server"
	"github.com/cybergarage/go-sqltest/sqltest/sysbench"
)

func TestSysbench(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	// Working directory

	wkdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Working directory: %s", wkdir)

	// Check cert files if exists
	// Running sysbench on RDS MySQL, RDS MariaDB, and Amazon Aurora MySQL via SSL/TLS
	// https://aws.amazon.com/jp/blogs/database/running-sysbench-on-rds-postgresql-rds-mariadb-and-amazon-aurora-postgresql-via-ssl-tls/

	certFiles := []string{
		"client-key.pem",
		"client-cert.pem",
		"cacert.pem",
	}

	for _, certFile := range certFiles {
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			t.Errorf("Cert file not found: %s", certFile)
			return
		}
	}

	// Setup server

	server := server.NewServer()
	server.SetTLSConfig(nil)
	err = server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	// Setup client

	client := postgresqltest.NewClient()

	err = client.Open()
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

	testDBName := sysbench.GenerateTempDBName()

	err =se(testDBName)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err := client.DropDatabase(testDBName)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// Setup sysbench configuation

	cfg := NewDefaultConfig()
	cfg.SetPort(server.Port())
	cfg.SetDB(testDBName)
	cfg.SetDBDebug(true)
	cfg.SetTableSize(100)

	cmds := []string{
		sysbench.OltpReadWrite,
	}

	for _, cmd := range cmds {
		t.Run(cmd, func(t *testing.T) {
			sysbench.RunCommand(t, cmd, cfg.Config)
		})
	}
}
