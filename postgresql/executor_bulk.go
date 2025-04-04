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

package postgresql

import (
	"github.com/cybergarage/go-postgresql/postgresql/errors"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// nullBulkExecutor represents a base bulk message executor.
type nullBulkExecutor struct {
}

// NewNullBulkExecutor returns a new null BulkQueryExecutor.
func NewNullBulkExecutor() BulkQueryExecutor {
	return &nullBulkExecutor{}
}

// Copy handles a COPY query.
func (executor *nullBulkExecutor) Copy(Conn, query.Copy) (protocol.Responses, error) {
	return nil, errors.NewErrNotImplemented("COPY")
}

// Copy handles a COPY DATA protocol.
func (executor *nullBulkExecutor) CopyData(Conn, query.Copy, *CopyStream) (protocol.Responses, error) {
	return nil, errors.NewErrNotImplemented("COPY DATA")
}
