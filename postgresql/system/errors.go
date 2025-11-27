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

import (
	"fmt"

	"github.com/cybergarage/go-postgresql/postgresql/errors"
)

// ErrNotFound represents a not found error.
var ErrNotFound = errors.ErrNotFound

// ErrNotImplemented represents a not implemented error.
var ErrNotImplemented = errors.ErrNotImplemented

// ErrNotSupported represents a not supported error.
var ErrNotSupported = errors.ErrNotSupported

// ErrInvalid represents an invalid error.
var ErrInvalid = errors.ErrInvalid

func newDataTypeNotFound(oid ObjectID) error {
	return fmt.Errorf("data type (%d) is %w", oid, ErrNotFound)
}
