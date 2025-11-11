// Copyright (C) 2022 The go-postgresql Authors. All rights reserved.
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

package fn

import (
	"fmt"
	"strings"

	"github.com/cybergarage/go-sqlparser/sql/fn"
)

// PostgreSQL: Documentation: 18: 9.27. System Information Functions and Operators
// https://www.postgresql.org/docs/current/functions-info.html
// 9.27.1. Session Information Functions
// https://www.postgresql.org/docs/current/functions-info.html#FUNCTIONS-INFO-SESSION

const (
	CurrentDatabaseFunctionName = "current_database"
	CurrentCatalogFunctionName  = "current_catalog"
	CurrentSchemaFunctionName   = "current_schema"
	CurrentSchemasFunctionName  = "current_schemas"
	CurrentUserFunctionName     = "current_user"
	CurrentRoleFunctionName     = "current_role"
	SessionUserFunctionName     = "session_user"
	UserFunctionName            = "user"
)

// SessionFunctionNames returns the all names of session functions.
func SessionFunctionNames() []string {
	return []string{
		CurrentDatabaseFunctionName,
		CurrentCatalogFunctionName,
		CurrentSchemaFunctionName,
		CurrentSchemasFunctionName,
		CurrentUserFunctionName,
		CurrentRoleFunctionName,
		SessionUserFunctionName,
		UserFunctionName,
	}
}

// IsSessionFunctionName returns true if the specified name is a session function name.
func IsSessionFunctionName(name string) bool {
	for _, fnName := range SessionFunctionNames() {
		if strings.EqualFold(fnName, name) {
			return true
		}
	}
	return false
}

// execImpl represents a base math function.
type sessionFunction struct {
	*execImpl
}

// NewSessionExecutor returns a new session function executor.
func NewSessionExecutor(name string, opts ...any) (Executor, error) {
	ex := &sessionFunction{
		execImpl: nil,
	}
	opts = append(opts,
		fn.WithExecutorName(strings.ToLower(name)),
		fn.WithExecutorType(fn.SessionFunction),
		fn.WithExecutorFunc(ex.execute),
	)
	ex.execImpl = newExecutorWith(opts...)
	return ex, nil
}

// Execute returns the executed value with the specified arguments.
func (ex *sessionFunction) execute(args ...any) (any, error) {
	conn := ex.Conn()
	if conn == nil {
		return nil, fmt.Errorf("%s: %w connection", ex.Name(), fn.ErrInvalid)
	}
	switch ex.Name() {
	case CurrentDatabaseFunctionName, CurrentCatalogFunctionName:
		return conn.Database(), nil
	case CurrentSchemasFunctionName:
		return conn.Schemas(), nil
	case CurrentSchemaFunctionName:
		schemas := conn.Schemas()
		if 0 < len(schemas) {
			return schemas[0], nil
		}
		return "", nil
	case CurrentUserFunctionName, CurrentRoleFunctionName, SessionUserFunctionName, UserFunctionName:
		return conn.User(), nil
	}
	return nil, fn.NewErrNotSupportedFunction(ex.Name())
}
