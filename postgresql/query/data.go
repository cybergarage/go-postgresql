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
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// PostgreSQL: Documentation: 16: Chapter 8. Data Types
// https://www.postgresql.org/docs/16/datatype.html

// DataType represents a PostgreSQL data type using the object ID.
type DataType = message.DataType

// DataTypeFrom returns a data type from the specified query data type.
func DataTypeFrom(t query.DataType) DataType {
	switch t {
	case query.BigIntData:
		return message.Int8Type
	case query.BinaryData:
		return message.ByteaType
	case query.BitData:
		return message.BitType
	case query.BlobData:
		return message.ByteaType
	case query.BooleanData:
		return message.BoolType
	case query.CharData:
		return message.CharType
	case query.CharacterData:
		return message.VarcharType
	// case query.ClobData:
	// 	return
	case query.DateData:
		return message.DateType
	// case query.DecimalData:
	// 	return
	case query.DoubleData:
		return message.Float8Type
	case query.FloatData:
		return message.Float4Type
	case query.IntData:
		return message.Int4Type
	case query.IntegerData:
		return message.Int4Type
	case query.LongBlobData:
		return message.ByteaType
	case query.LongTextData:
		return message.TextType
	case query.MediumBlobData:
		return message.ByteaType
	// case query.MediumIntData:
	// 	return
	case query.MediumTextData:
		return message.TextType
	// case query.NumericData:
	// 	return
	// case query.RealData:
	// 	return
	// case query.SetData:
	// 	return
	case query.SmallIntData:
		return message.Int2Type
	case query.TextData:
		return message.TextType
	case query.TimeData:
		return message.TimeType
	case query.TimeStampData:
		return message.TimestampType
	case query.TinyBlobData:
		return message.ByteaType
	// case query.TinyIntData:
	// 	return
	case query.TinyTextData:
		return message.TextType
	case query.VarBinaryData:
		return message.ByteaType
	case query.VarCharData:
		return message.VarcharType
	case query.VarCharacterData:
		return message.VarcharType
		// case query.YearData:
		// 	return
	}
	return 0
}
