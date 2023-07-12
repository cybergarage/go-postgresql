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

package server

import (
	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// Authenticate authenticates the connection with the startup message.
func (server *Server) Authenticate(*postgresql.Conn, *message.Startup) (message.Response, error) {
	return message.NewAuthenticationOk()
}