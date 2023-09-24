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
	"errors"
	"fmt"
)

// ErrNotImplemented is returned when the operation is not implemented.
var ErrNotImplemented = errors.New("not implemented")

// ErrNotSupported is returned when the operation is not supported.
var ErrNotSupported = errors.New("not supported")

// ErrNotExist is returned when the specified object is not exist.
var ErrNotExist = errors.New("not exist")

// ErrExist is returned when the specified object is exist.
var ErrExist = errors.New("exist")

// ErrNotEqual is returned when the specified object is not equal.
var ErrNotEqual = errors.New("not equal")

// NewErrNotImplemented returns a new ErrNotImplemented error.
func NewErrNotImplemented(msg string) error {
	return fmt.Errorf("%s is %w", msg, ErrNotImplemented)
}

// NewErrNotSupported returns a new ErrNotSupported error.
func NewErrNotSupported(msg string) error {
	return fmt.Errorf("%s is %w", msg, ErrNotSupported)
}

// NewErrExist returns a new exist error.
func NewErrExist(v any) error {
	return fmt.Errorf("%v is %w", v, ErrExist)
}

// NewErrNotExist returns a new not exist error.
func NewErrNotExist(v any) error {
	return fmt.Errorf("%v is %w", v, ErrNotExist)
}

// NewErrNotEqual returns a new not equal error.
func NewErrNotEqual(v1, v2 any) error {
	return fmt.Errorf("%w (%v != %v) ", ErrNotEqual, v1, v2)
}

// NewErrDatabaseNotExist returns a new database not exist error.
func NewErrDatabaseNotExist(v string) error {
	return fmt.Errorf("database (%v) is %w", v, ErrNotExist)
}

// NewErrTableNotExist returns a new table not exist error.
func NewErrTableNotExist(v string) error {
	return fmt.Errorf("table (%v) is %w", v, ErrNotExist)
}

// NewErrDatabaseExist returns a new database exist error.
func NewErrDatabaseExist(v string) error {
	return fmt.Errorf("database (%v) is %w", v, ErrExist)
}

// NewErrTableExist returns a new table exist error.
func NewErrTableExist(v string) error {
	return fmt.Errorf("table (%v) is %w", v, ErrExist)
}

// NewErrColumnNotExist returns a new column not exist error.
func NewErrColumnNotExist(v any) error {
	return fmt.Errorf("column (%v) is %w", v, ErrNotExist)
}

// NewErrColumnValueNotExist returns a new column value not exist error.
func NewErrColumnValueNotExist(v any) error {
	return fmt.Errorf("column (%v) value is %w", v, ErrNotExist)
}

// NewErrGroupByColumnValueNotExist returns a new group by column not exist error.
func NewErrGroupByColumnValueNotExist(v any) error {
	return fmt.Errorf("group by column (%v) value is %w", v, ErrNotExist)
}

// NewErrColumnsNotEqual returns a new columns not equal error.
func NewErrColumnsNotEqual(v1, v2 int) error {
	return fmt.Errorf("the number of columns (%d) is %w to the number of schema columns (%d)", v1, ErrNotEqual, v2)
}

// NewErrPreparedStatementNotExist returns a new prepared statement not exist error.
func NewErrPreparedStatementNotExist(name string) error {
	return fmt.Errorf("prepared statement (%v) is %w", name, ErrNotExist)
}

// NewErrPreparedPortalNotExist returns a new prepared portal not exist error.
func NewErrPreparedPortalNotExist(name string) error {
	return fmt.Errorf("prepared portal (%v) is %w", name, ErrNotExist)
}

// NewErrPreparedStatementMultiStatement returns a new prepared statement multi statement error.
func NewErrMultiplePreparedStatementNotSupported(query string) error {
	return fmt.Errorf("multiple prepared statement (%v) is %w", query, ErrNotSupported)
}

// NewErrMultipleTableNotSupported returns a new prepared statement multi table error.
func NewErrMultipleTableNotSupported(query string) error {
	return fmt.Errorf("multiple table (%v) is %w", query, ErrNotSupported)
}
