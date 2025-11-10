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
	"github.com/cybergarage/go-sqlparser/sql/net"
)

// ExecutorOption is a function type for configuring an executor.
type ExecutorOption func(*execImpl)

// execImpl represents a base math function.
type execImpl struct {
	name string
	args []string
	conn net.Conn
}

// WithExecutorName sets the name for the executor.
func WithExecutorName(name string) ExecutorOption {
	return func(ex *execImpl) {
		ex.name = name
	}
}

// WithExecutorArguments sets the arguments for the executor.
func WithExecutorArguments(args []string) ExecutorOption {
	return func(ex *execImpl) {
		ex.args = args
	}
}

// WithExecutorConn sets the connection for the executor.
func WithExecutorConn(conn net.Conn) ExecutorOption {
	return func(ex *execImpl) {
		ex.conn = conn
	}
}

// NewExecutorWith returns a new function executor with options.
func NewExecutorWith(opts ...ExecutorOption) Executor {
	return newExecutorWith(opts...)
}

func newExecutorWith(opts ...ExecutorOption) *execImpl {
	ex := &execImpl{
		name: "",
		args: []string{},
		conn: nil,
	}
	for _, opt := range opts {
		opt(ex)
	}
	return ex
}

// Name returns the name of the function.
func (ex *execImpl) Name() string {
	return ex.name
}

// Arguments returns the arguments of the executor.
func (ex *execImpl) Arguments() []string {
	return ex.args
}
