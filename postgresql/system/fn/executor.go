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
	"github.com/cybergarage/go-sqlparser/sql/fn"
)

// Executor represents a function executor.
type Executor = fn.Executor

// NewExecutorForName returns a function executor with the specified name.
func NewExecutorForName(name string, opts ...any) (Executor, error) {
	return NewExecutorWith(opts...), nil
}
