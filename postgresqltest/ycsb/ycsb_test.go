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

package ycsb

import (
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresqltest/client"
	"github.com/cybergarage/go-postgresql/postgresqltest/server"
)

func TestYCSB(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	// Setup client

	client := client.NewDefaultClient()
	client.SetDatabase(ycsbDatabaseName)
	err = client.CreateDatabase(ycsbDatabaseName)
	if err != nil {
		t.Error(err)
		return
	}

	// Setup for YCSB benchmark

	err = client.Open()
	if err != nil {
		t.Error(err)
		return
	}

	for _, setupQuery := range setUpQueries {
		rs, err := client.Query(setupQuery)
		if err != nil {
			t.Error(err)
		}
		defer rs.Close()
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// Tries to execute ycsb command

	workloads := []string{
		"workloada",
	}

	for _, workload := range workloads {
		t.Run(workload, func(t *testing.T) {
			RunYCSBWorkload(t, workload)
		})
	}
}
