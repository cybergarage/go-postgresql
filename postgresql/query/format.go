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
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// PostgreSQL: Documentation: 16: Chapter 8. Data Types
// https://www.postgresql.org/docs/16/datatype.html

// FormatCode represents a format code.
type FormatCode int16

const (
	TextFormat   FormatCode = 0
	BinaryFormat FormatCode = 1
)

// FormatCodeFrom returns a format code from the specified data type.
func FormatCodeFrom(t query.DataType) FormatCode {
	switch t {
	case query.LongTextData:
		return TextFormat
	case query.MediumTextData:
		return TextFormat
	case query.TextData:
		return TextFormat
	case query.TinyTextData:
		return TextFormat
	case query.VarCharData:
		return TextFormat
	}
	return BinaryFormat
}
