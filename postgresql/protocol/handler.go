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

package protocol

// StartupHandler represents a start-up message handler.
type StartupHandler interface {
	// ParameterStatuses returns the parameter statuses.
	ParameterStatuses(Conn) (Responses, error)
	// BackendKeyData returns the backend key data.
	BackendKeyData(Conn) (Response, error)
}

// SimpleQueryHandler defines a executor interface for simple query operations.
type SimpleQueryHandler interface {
	// Query handles a query
	Query(Conn, *Query) (Responses, error)
}

// ExtendedQueryHandler defines a executor interface for extended query operations.
type ExtendedQueryHandler interface {
	// Prepare handles a parse
	Parse(Conn, *Parse) (Responses, error)
	// Bind handles a bind
	Bind(Conn, *Bind) (Responses, error)
	// Describe handles a describe
	Describe(Conn, *Describe) (Responses, error)
	// Execute handles a execute
	Execute(Conn, *Execute) (Responses, error)
	// Close handles a close
	Close(Conn, *Close) (Responses, error)
	// Sync handles a sync
	Sync(Conn, *Sync) (Responses, error)
	// Flush handles a flush
	Flush(Conn, *Flush) (Responses, error)
}

// QueryHandler represents a query handler.
type QueryHandler interface {
	SimpleQueryHandler
	ExtendedQueryHandler
}

// MessageHandler represents a message handler.
type MessageHandler interface {
	StartupHandler
	QueryHandler
}
