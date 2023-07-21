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

package query

import (
	"github.com/cybergarage/go-sqlparser/sql/query"
)

const (
	// Aliases for types defined in github.com/cybergarage/go-postgresql/postgresql/query.
	EQ  = query.EQ
	NEQ = query.NEQ
	LT  = query.LT
	LE  = query.LE
	GT  = query.GT
	GE  = query.GE
	IN  = query.IN
	NIN = query.NIN
)

type (
	// Aliases for types defined in github.com/cybergarage/go-postgresql/postgresql/query.
	BindParam      = query.BindParam
	CreateDatabase = query.CreateDatabase
	CreateTable    = query.CreateTable
	CreateIndex    = query.CreateIndex
	DropDatabase   = query.DropDatabase
	DropTable      = query.DropTable
	Select         = query.Select
	Insert         = query.Insert
	Update         = query.Update
	Delete         = query.Delete
	Column         = query.Column
	Table          = query.Table
	Condition      = query.Condition
	Schema         = query.Schema
	Expr           = query.Expr
	CmpExpr        = query.CmpExpr
	AndExpr        = query.AndExpr
	OrExpr         = query.OrExpr
)