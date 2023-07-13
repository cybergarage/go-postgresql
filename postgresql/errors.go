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
)

// ErrNotImplemented is returned when the operation is not implemented.
var ErrNotImplemented = errors.New("not implemented")

// ErrNotSupported is returned when the operation is not supported.
var ErrNotSupported = errors.New("not supported")

// ErrNotExist is returned when the specified object is not exist.
var ErrNotExist = errors.New("not exist")

// ErrExist is returned when the specified object is exist.
var ErrExist = errors.New("exist")

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
