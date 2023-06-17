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

package postgresql

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// Authenticator represents a frontend message authenticator.
type Authenticator interface {
	// Authenticate authenticates the connection.
	Authenticate(*Conn, *message.Startup) (message.Response, error)
}

// StatusExecutor represents a backend status message executor.
type StatusExecutor interface {
	// ParameterStatus returns the parameter status.
	ParameterStatus(*Conn) (message.Response, error)
	// BackendKeyData returns the backend key data.
	BackendKeyData(*Conn) (message.Response, error)
}

type QueryExecutor interface {
	// Parse returns the parse response.
	Parse(*Conn, *message.Parse) (message.Response, error)
	// Bind returns the bind response.
	Bind(*Conn, *message.Bind) (message.Response, error)
}

// Executor represents a frontend message executor.
type Executor interface {
	Authenticator
	StatusExecutor
	QueryExecutor
}
