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

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresqltest/client"
)

const testDBNamePrefix = "pgtest"

type ServerTestFunc = func(*testing.T, *client.PqClient)

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

	testFuncs := []struct {
		name string
		fn   ServerTestFunc
	}{
		{"copy", TestServerCopy},
	}

	for _, testFunc := range testFuncs {
		t.Run(testFunc.name, func(t *testing.T) {
			testFunc.fn(t, client)
		})
	}
}

// TestServerCopy tests the COPY command.
func TestServerCopy(t *testing.T, client *client.PqClient) {
	_, err := client.Query("CREATE TABLE cptest (ctext TEXT PRIMARY KEY, cint INT, cfloat FLOAT, cdouble DOUBLE);")
	if err != nil {
		t.Error(err)
		return
	}

	db := client.DB()
	txn, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}

	stmt, err := txn.Prepare(pq.CopyIn("cptest", "ctext", "cint", "cfloat", "cdouble"))
	if err != nil {
		t.Error(err)
		return
	}

	records := [][]interface{}{
		{"text1", 1, 1.1, 1.11},
		{"text2", 2, 2.2, 2.22},
		{"text3", 3, 3.3, 3.33},
	}

	for _, record := range records {
		_, err = stmt.Exec(record...)
		if err != nil {
			t.Error(err)
			return
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		t.Error(err)
		return
	}

	err = stmt.Close()
	if err != nil {
		t.Error(err)
		return
	}

	err = txn.Commit()
	if err != nil {
		t.Error(err)
		return
	}
}
