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
	"fmt"

	"github.com/cybergarage/go-postgresql/postgresql/protocol"
)

// nullErrorHandler represents a base error handler.
type nullErrorHandler struct {
}

// NewNullErrorHandler returns a new null ErrorHandler.
func NewNullErrorHandler() ErrorHandler {
	return &nullErrorHandler{}
}

// ParserError handles a parser error.
func (executor *nullErrorHandler) ParserError(conn Conn, q string, err error) (protocol.Responses, error) {
	resErr := fmt.Errorf("parser error : %w", err)
	res, err := protocol.NewErrorResponseWith(resErr)
	if err != nil {
		return nil, err
	}
	return protocol.NewResponsesWith(res), nil
}
