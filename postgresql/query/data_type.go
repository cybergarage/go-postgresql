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
	"fmt"

	"github.com/cybergarage/go-postgresql/postgresql/errors"
	"github.com/cybergarage/go-postgresql/postgresql/system"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// PostgreSQL: Documentation: 16: Chapter 8. Data s
// https://www.postgresql.org/docs/16/datatype.html

// DataType represents a data type.
type DataType = system.DataType

// ObjectID represents a object identifier.
type ObjectID = system.ObjectID

// NewDataTypeFrom returns a data type from the specified query data type.
func NewDataTypeFrom(t query.DataType) (*DataType, error) {
	objectID, err := NewObjectIDFrom(t)
	if err != nil {
		return nil, err
	}
	return system.NewDataTypeFromObjectID(objectID)
}

// NewObjectIDFrom returns a data type from the specified query data type.
func NewObjectIDFrom(t query.DataType) (ObjectID, error) {
	switch t {
	case query.BigIntType:
		return system.Int8, nil
	case query.BinaryType:
		return system.Bytea, nil
	case query.BitType:
		return system.Bit, nil
	case query.BlobType:
		return system.Bytea, nil
	case query.BooleanType:
		return system.Bool, nil
	case query.CharType:
		return system.Char, nil
	case query.CharacterType:
		return system.Varchar, nil
	// case query.ClobData:
	// 	return
	case query.DateType:
		return system.Date, nil
	// case query.DecimalData:
	// 	return
	case query.DoubleType:
		return system.Float8, nil
	case query.FloatType:
		return system.Float4, nil
	case query.IntType:
		return system.Int4, nil
	case query.IntegerType:
		return system.Int4, nil
	case query.LongBlobType:
		return system.Bytea, nil
	case query.LongTextType:
		return system.Text, nil
	case query.MediumBlobType:
		return system.Bytea, nil
	// case query.MediumIntData:
	// 	return
	case query.MediumTextType:
		return system.Text, nil
	// case query.NumericData:
	// 	return
	// case query.RealData:
	// 	return
	// case query.SetData:
	// 	return
	case query.SmallIntType:
		return system.Int2, nil
	case query.TextType:
		return system.Text, nil
	case query.TimeType:
		return system.Time, nil
	case query.TimeStampType:
		return system.Timestamp, nil
	case query.DateTimeType:
		return system.Timestamp, nil
	case query.TinyBlobType:
		return system.Bytea, nil
	// case query.TinyIntData:
	// 	return
	case query.TinyTextType:
		return system.Text, nil
	case query.VarBinaryType:
		return system.Bytea, nil
	case query.VarCharType:
		return system.Varchar, nil
	case query.VarCharacterType:
		return system.Varchar, nil
		// case query.YearData:
		// 	return
	}
	return 0, fmt.Errorf("data type (%s) %w", t, errors.ErrNotSupported)
}
