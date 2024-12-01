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
	return system.NewDataTypeFrom(objectID)
}

// NewObjectIDFrom returns a data type from the specified query data type.
func NewObjectIDFrom(t query.DataType) (ObjectID, error) {
	switch t {
	case query.BigIntData:
		return system.Int8, nil
	case query.BinaryData:
		return system.Bytea, nil
	case query.BitData:
		return system.Bit, nil
	case query.BlobData:
		return system.Bytea, nil
	case query.BooleanData:
		return system.Bool, nil
	case query.CharData:
		return system.Char, nil
	case query.CharacterData:
		return system.Varchar, nil
	// case query.ClobData:
	// 	return
	case query.DateData:
		return system.Date, nil
	// case query.DecimalData:
	// 	return
	case query.DoubleData:
		return system.Float8, nil
	case query.FloatData:
		return system.Float4, nil
	case query.IntData:
		return system.Int4, nil
	case query.IntegerData:
		return system.Int4, nil
	case query.LongBlobData:
		return system.Bytea, nil
	case query.LongTextData:
		return system.Text, nil
	case query.MediumBlobData:
		return system.Bytea, nil
	// case query.MediumIntData:
	// 	return
	case query.MediumTextData:
		return system.Text, nil
	// case query.NumericData:
	// 	return
	// case query.RealData:
	// 	return
	// case query.SetData:
	// 	return
	case query.SmallIntData:
		return system.Int2, nil
	case query.TextData:
		return system.Text, nil
	case query.TimeData:
		return system.Time, nil
	case query.TimeStampData:
		return system.Timestamp, nil
	case query.DateTimeData:
		return system.Timestamp, nil
	case query.TinyBlobData:
		return system.Bytea, nil
	// case query.TinyIntData:
	// 	return
	case query.TinyTextData:
		return system.Text, nil
	case query.VarBinaryData:
		return system.Bytea, nil
	case query.VarCharData:
		return system.Varchar, nil
	case query.VarCharacterData:
		return system.Varchar, nil
		// case query.YearData:
		// 	return
	}
	return 0, fmt.Errorf("data type (%s) %w", t, errors.ErrNotSupported)
}
