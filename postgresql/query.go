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
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

// PreparedQuery represents a prepared query.
type PreparedQuery = message.Parse

// PreparedQueryMap represents a prepared query map.
type PreparedQueryMap map[string]*PreparedQuery

// NewPreparedQueryMap returns a new prepared query map.
func NewPreparedQueryMap() PreparedQueryMap {
	return make(PreparedQueryMap)
}

// PreparedQuery returns a prepared query.
func (queries PreparedQueryMap) PreparedQuery(name string) (*PreparedQuery, error) {
	q, oK := queries[name]
	if !oK {
		return nil, query.NewErrPreparedQueryNotExist(name)
	}
	return q, nil
}

// SetPreparedQuery sets a prepared query.
func (queries PreparedQueryMap) SetPreparedQuery(query *PreparedQuery) error {
	queries[query.Name] = query
	return nil
}
