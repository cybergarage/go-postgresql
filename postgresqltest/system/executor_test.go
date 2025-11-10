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

package system

import (
	"testing"

	"github.com/cybergarage/go-postgresql/postgresql/system/fn"
	"github.com/cybergarage/go-sqlparser/sql/net"
)

func TestSystemFunction(t *testing.T) {
	functionNames := []string{
		fn.CurrentDatabaseFunctionName,
		fn.CurrentCatalogFunctionName,
	}
	con := net.NewConnWith(nil)
	con.SetDatabase("testdb")
	for _, functionName := range functionNames {
		t.Run(functionName, func(t *testing.T) {
			executor, err := fn.NewExecutorForName(
				functionName,
				fn.WithExecutorConn(con),
			)
			if err != nil {
				t.Errorf("Failed to create executor for function '%s': %v", functionName, err)
			}
			_, err = executor.Execute(nil)
			if err != nil {
				t.Errorf("Failed to execute function '%s': %v", functionName, err)
			}
		})
	}
}
