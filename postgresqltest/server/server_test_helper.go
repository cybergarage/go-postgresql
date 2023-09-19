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
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresqltest/client"
)

const testDBNamePrefix = "pgtest"

func RunServerTests(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	testDBName := fmt.Sprintf("%s%d", testDBNamePrefix, time.Now().UnixNano())

	client := client.NewPqClient()
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

	err = client.CreateDatabase(testDBName)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err := client.DropDatabase(testDBName)
		if err != nil {
			t.Error(err)
		}
	}()

}
