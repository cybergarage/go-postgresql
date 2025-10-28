// Copyright (C) 2025 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package benchbase contains integration tests that drive BenchBase (CMU DB) workloads
// against the embedded PostgreSQL-compatible server implementation.
package benchbase

import (
	"os"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresqltest"
	"github.com/cybergarage/go-postgresql/postgresqltest/server"
	"github.com/cybergarage/go-sqltest/sqltest/benchbase"
)

// TestBenchBase runs one or more BenchBase workloads against a temporary database.
// It mirrors the structure of the sysbench integration test.
//
// Skip Conditions:
// - BenchBase CLI/tooling not installed.
// - Missing workload runner API (adjust TODO sections as needed).
func TestBenchBase(t *testing.T) {
	// Enable verbose debug logging to observe benchmark progress.
	log.EnableStdoutDebug(true)

	// Skip early if BenchBase tooling is not available on this system.
	// if !benchbase.IsInstalled() {
	// 	t.Skip("BenchBase is not installed; skipping test")
	// 	return
	// }

	// Log working directory to help debug path-related issues or relative resource lookups.
	wkdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Working directory: %s", wkdir)

	// Start an embedded server instance. TLS is disabled here; add configuration if needed.
	srv := server.NewServer()
	srv.SetTLSConfig(nil)
	if err := srv.Start(); err != nil {
		t.Error(err)
		return
	}
	defer srv.Stop()

	// Open a client connection to manage test database lifecycle.
	client := postgresqltest.NewClient()
	if err := client.Open(); err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			t.Error(cerr)
		}
	}()

	// Create a temporary database for isolation. Fall back to a static name if helper is absent.
	// testDBName := benchbase.GenerateTempDBName()
	// if testDBName == "" {
	// 	testDBName = "benchbase_tempdb"
	// }
	// if err := client.CreateDatabase(testDBName); err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// defer func() {
	// 	if derr := client.DropDatabase(testDBName); derr != nil {
	// 		t.Error(derr)
	// 	}
	// }()

	// // Build a BenchBase config tailored for the spawned server.
	// cfg := NewDefaultConfig()
	// cfg.SetHost("127.0.0.1")
	// cfg.SetPort(srv.Port())
	// cfg.SetDB(testDBName)
	// // If SetUser / SetPassword are required, ensure they're configured:
	// cfg.SetUser("postgres")
	// cfg.SetPassword("postgres")

	// List of workloads to execute; expand as needed.
	// Common BenchBase workloads include: tpcc, tatp, smallbank, ycsb, epinions, etc.
	workloads := []string{
		"tpcc",
	}

	// Each workload is run as a subtest for isolated reporting in go test output.
	for _, wl := range workloads {
		wl := wl // shadow for closure capture
		t.Run(wl, func(t *testing.T) {
			benchbase.RunWorkload(t, wl /*cfg.Config*/)
		})
	}
}
