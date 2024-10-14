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
	"errors"
	"fmt"
	"io"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
	sql "github.com/cybergarage/go-sqlparser/sql/query"
)

// NewCopyInResponsesFrom returns a new copy in response from the specified query.
func NewCopyInResponsesFrom(q query.Copy, schema *sql.Schema) (protocol.Responses, error) {
	// PostgreSQL: Documentation: 16: COPY
	// https://www.postgresql.org/docs/16/sql-copy.html

	copyColums := q.Columns()
	if 0 < len(copyColums) {
		for _, copyColumn := range q.Columns() {
			_, err := schema.ColumnByName(copyColumn.Name())
			if err != nil {
				return nil, err
			}
		}
	} else {
		copyColums = schema.Columns()
	}

	// Support only text format
	res := protocol.NewCopyInResponseWith(protocol.TextCopy)
	for n := 0; n < len(copyColums); n++ {
		res.AppendFormatCode(protocol.TextFormat)
	}

	return protocol.NewResponsesWith(res), nil
}

// NewCopyQueryFrom returns a new copy query from the specified query.
func NewCopyQueryFrom(schema *query.Schema, copyColumns sql.ColumnList, copyData *protocol.CopyData) (query.Insert, error) {
	// COPY FROM will raise an error if any line of the input file contains
	// more or fewer columns than are expected.
	copyColumData := copyData.Data
	if len(copyColumData) != len(copyColumns) {
		return nil, query.NewErrColumnsNotEqual(len(copyColumData), len(copyColumns))
	}

	columns := make(sql.ColumnList, len(copyColumns))
	for idx, copyColumn := range copyColumns {
		copyColumnData := copyColumData[idx]
		if len(copyColumnData) == 0 {
			copyColumnData = "NULL"
		}
		columns[idx] = sql.NewColumnWithOptions(
			sql.WithColumnName(copyColumn.Name()),
			sql.WithColumnLiteral(sql.NewLiteralWith(copyColumnData)),
		)
	}
	return sql.NewInsertWith(sql.NewTableWith(schema.Name()), columns), nil
}

// NewCopyCompleteResponsesFrom returns a new copy complete response from the specified query.
func NewCopyCompleteResponsesFrom(q query.Copy, stream *CopyStream, conn Conn, schema *sql.Schema, queryExecutor QueryExecutor) (protocol.Responses, error) {
	copyData := func(schema *query.Schema, colums sql.ColumnList, copyData *protocol.CopyData) error {
		q, err := NewCopyQueryFrom(schema, colums, copyData)
		if err != nil {
			return err
		}
		log.Tracef("%s", q.String())
		_, err = queryExecutor.Insert(conn, q)
		return err
	}

	copyColums := q.Columns()
	if len(copyColums) == 0 {
		copyColums = schema.Columns()
	}

	nCopy := 0
	nFail := 0
	cpData, err := stream.Next()
	for {
		if err != nil {
			break
		}
		if err := copyData(schema, copyColums, cpData); err != nil {
			nFail++
			log.Errorf("%s (%d/%d) (%s)", q.String(), nCopy, nFail, err)
		} else {
			nCopy++
		}
		cpData, err = stream.Next()
	}

	if !errors.Is(err, io.EOF) {
		log.Errorf("%s (%d/%d) (%s)", q.String(), nCopy, nFail, err.Error())
		return nil, err
	}

	if 0 < nFail {
		return nil, fmt.Errorf("%s (%d/%d)", q.String(), nCopy, nFail)
	}

	return protocol.NewCopyCompleteResponsesWith(nCopy)
}
