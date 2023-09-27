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
	"github.com/cybergarage/go-postgresql/postgresql/system"
)

var (
	pgbenchGetPartitionQuery = "select o.n, p.partstrat, pg_catalog.count(i.inhparent) from pg_catalog.pg_class"
)

// IsMatchQuery returns true if the query is matched with the prefix.
func IsMatchQuery(q string, prefix string) bool {
	return q[:len(prefix)] == prefix
}

// IsPgbenchGetPartitionQuery returns true if the query is pgbenchGetPartitionQuery.
func IsPgbenchGetPartitionQuery(q string) bool {
	return IsMatchQuery(q, pgbenchGetPartitionQuery)
}

// NewGetPartitionResponseForPgbench returns a new response for pgbenchGetPartitionQuery.
func NewGetPartitionResponseForPgbench() (message.Responses, error) {
	// PostgreSQL: Documentation: 16: 55.7. Message Formats
	// https://www.postgresql.org/docs/16/protocol-message-formats.html
	// pgbench - void GetTableInfo(PGconn *con, bool scale_given)
	// https://github.com/postgres/postgres/blob/master/src/bin/pgbench/pgbench.c
	// 00000000  54 00 00 00 4e 00 03 6e  00 00 00 00 00 00 00 00   T...N..n ........
	// 00000010  00 00 17 00 04 ff ff ff  ff 00 00 70 61 72 74 73   ........ ...parts
	// 00000020  74 72 61 74 00 00 00 0d  16 00 02 00 00 00 12 00   trat.... ........
	// 00000030  01 ff ff ff ff 00 00 63  6f 75 6e 74 00 00 00 00   .......c ount....
	// 00000040  00 00 00 00 00 00 14 00  08 ff ff ff ff 00 00 44   ........ .......D
	// 00000050  00 00 00 14 00 03 00 00  00 01 32 ff ff ff ff 00   ........ ..2.....
	// 00000060  00 00 01 30 43 00 00 00  0d 53 45 4c 45 43 54 20   ...0C... .SELECT
	// 00000070  31 00 5a 00 00 00 05 49                            1.Z....I

	rowDesc := message.NewRowDescription()
	dataRow := message.NewDataRow()

	resFieldNames := []string{
		"n",
		"partstrat",
		"count",
	}

	for n, fieldName := range resFieldNames {
		switch n {
		case 0:
			dt, _ := system.GetDataType(system.Int4)
			rowField := message.NewRowFieldWith(fieldName,
				message.WithRowFieldNumber(int16(n+1)),
				message.WithRowFieldDataType(dt),
			)
			dataRow.AppendData(rowField, 2)
		case 1:
			dt, _ := system.GetDataType(system.Char)
			rowField := message.NewRowFieldWith(fieldName,
				message.WithRowFieldNumber(int16(n+1)),
				message.WithRowFieldDataType(dt),
			)
			dataRow.AppendData(rowField, nil)
		case 2:
			dt, _ := system.GetDataType(system.Int8)
			rowField := message.NewRowFieldWith(fieldName,
				message.WithRowFieldNumber(int16(n+1)),
				message.WithRowFieldDataType(dt),
			)
			dataRow.AppendData(rowField, 0)
		}
	}

	res := message.NewResponses()
	res = res.Append(rowDesc)
	res = res.Append(dataRow)

	cmpRes, err := message.NewSelectCompleteWith(1)
	if err != nil {
		return nil, err
	}
	res = res.Append(cmpRes)

	return res, nil
}
