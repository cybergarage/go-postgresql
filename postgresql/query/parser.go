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

package query

import (
	"github.com/cybergarage/go-sqlparser/sql"
)

// Parse returns a SQL parser.
type Parser struct {
	sql.Parser
}

// NewParser returns a new parser.
func NewParser() *Parser {
	return &Parser{
		Parser: sql.NewParser(),
	}
}

// ParseString parses the specified query string and returns statements.
func (parser *Parser) ParseString(query string) ([]*Statement, error) {
	stmts, err := parser.Parser.ParseString(query)
	if err != nil {
		return nil, err
	}
	pgStmts := make([]*Statement, len(stmts))
	for n, stmt := range stmts {
		pgStmts[n] = NewStatementWith(stmt)
	}
	return pgStmts, nil
}
