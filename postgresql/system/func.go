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

package system

import "github.com/cybergarage/go-sqlparser/sql/fn"

// PostgreSQL: Documentation: 16: 9.21. Aggregate Functions
// https://www.postgresql.org/docs/16/functions-aggregate.html

// GetFunctionDataType returns the data type of the specified function.
func GetFunctionDataType(fx fn.Function, dt *DataType) (*DataType, error) {
	switch fx.Name() {
	case fn.MaxFunctionName:
		return dt, nil
	case fn.MinFunctionName:
		return dt, nil
	case fn.SumFunctionName:
		switch dt.OID() {
		case Int2:
			return dataTypes[Int8], nil
		case Int4:
			return dataTypes[Int8], nil
		}
	// NOTE: bigint for any integer-type argument instead of numeric
	case fn.AvgFunctionName:
		return dataTypes[Float8], nil
	case fn.CountFunctionName:
		return dataTypes[Int8], nil
	case fn.AbsFunctionName:
		return dataTypes[Float8], nil
	case fn.FloorFunctionName:
		return dataTypes[Float8], nil
	case fn.CeilFunctionName:
		return dataTypes[Float8], nil
	}
	return dt, nil
}
