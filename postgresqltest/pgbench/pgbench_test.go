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
	_ "embed"
	"os/exec"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresqltest/server"
)

func BenchmarkPgBench(b *testing.B) {
	log.EnableStdoutDebug(true)

	server := server.NewServer()

	err := server.Start()
	if err != nil {
		b.Error(err)
		return
	}

	scripts := []string{
		"./pgbench-init",
		"./pgbench-run",
	}

	for _, script := range scripts {
		cmd := exec.Command(script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			b.Skip(err)
			return
		}
		b.Log(string(output))
	}

	err = server.Stop()
	if err != nil {
		b.Error(err)
		return
	}
}
