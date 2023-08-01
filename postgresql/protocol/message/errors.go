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

package message

import (
	"errors"
	"fmt"
)

// ErrInvalid is returned when the message is invalid.
var ErrInvalid = errors.New("invalid")

// ErrNotSupported is returned when the message is not supported.
var ErrNotSupported = errors.New("not supported")

// ErrNotExist is returned when the specified object is not exist.
var ErrNotExist = errors.New("not exist")

// ErrExist is returned when the specified object is exist.
var ErrExist = errors.New("exist")

func newShortMessageError(expected int, actual int) error {
	return fmt.Errorf("%w short message : %d < %d", ErrInvalid, actual, expected)
}

func newColumnTypeNotSuppotedError(v any) error {
	return fmt.Errorf("column value type: %T is %w", v, ErrNotSupported)
}

func newInvalidLengthError(v int) error {
	return fmt.Errorf("%d is %w length", v, ErrInvalid)
}

// NewErrMessageNotSuppoted returns a new message not supported error.
func NewErrMessageNotSuppoted(t Type) error {
	return fmt.Errorf("message type (%c) is %w", t, ErrNotSupported)
}

// NewErrExist returns a new exist error.
func NewErrExist(v any) error {
	return fmt.Errorf("%v is %w", v, ErrExist)
}

// NewErrNotExist returns a new not exist error.
func NewErrNotExist(v any) error {
	return fmt.Errorf("%v is %w", v, ErrNotExist)
}

// NewErrInvalidMessage eturns a new message not supported error.
func NewErrInvalidMessage(t Type) error {
	return fmt.Errorf("message type (%c) is %w", t, ErrInvalid)
}
