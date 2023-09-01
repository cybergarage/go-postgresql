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
	"github.com/cybergarage/go-postgresql/postgresql/system"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// PostgreSQL: Documentation: 16: Chapter 8. Data s
// https://www.postgresql.org/docs/16/datatype.html

// DataType represents a data type.
type DataType = system.DataType

// OID represents a object identifier.
type OID = system.OID

// DataFrom returns a data type from the specified query data type.
func DataFrom(t query.DataType) OID {
	switch t {
	case query.BigIntData:
		return system.Int8
	case query.BinaryData:
		return system.Bytea
	case query.BitData:
		return system.Bit
	case query.BlobData:
		return system.Bytea
	case query.BooleanData:
		return system.Bool
	case query.CharData:
		return system.Char
	case query.CharacterData:
		return system.Varchar
	// case query.ClobData:
	// 	return
	case query.DateData:
		return system.Date
	// case query.DecimalData:
	// 	return
	case query.DoubleData:
		return system.Float8
	case query.FloatData:
		return system.Float4
	case query.IntData:
		return system.Int4
	case query.IntegerData:
		return system.Int4
	case query.LongBlobData:
		return system.Bytea
	case query.LongTextData:
		return system.Text
	case query.MediumBlobData:
		return system.Bytea
	// case query.MediumIntData:
	// 	return
	case query.MediumTextData:
		return system.Text
	// case query.NumericData:
	// 	return
	// case query.RealData:
	// 	return
	// case query.SetData:
	// 	return
	case query.SmallIntData:
		return system.Int2
	case query.TextData:
		return system.Text
	case query.TimeData:
		return system.Time
	case query.TimeStampData:
		return system.Timestamp
	case query.TinyBlobData:
		return system.Bytea
	// case query.TinyIntData:
	// 	return
	case query.TinyTextData:
		return system.Text
	case query.VarBinaryData:
		return system.Bytea
	case query.VarCharData:
		return system.Varchar
	case query.VarCharacterData:
		return system.Varchar
		// case query.YearData:
		// 	return
	}
	return 0
}
