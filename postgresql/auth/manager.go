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

package auth

import (
	"errors"

	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// AuthManager represent an authenticator manager.
type AuthManager struct {
	authenticators []Authenticator
}

// NewAuthManager returns a new authenticator manager.
func NewAuthManager() *AuthManager {
	manager := &AuthManager{
		authenticators: make([]Authenticator, 0),
	}
	return manager
}

// AddAuthenticator adds a new authenticator.
func (mgr *AuthManager) AddAuthenticator(authenticator Authenticator) {
	mgr.authenticators = append(mgr.authenticators, authenticator)
}

// ClearAuthenticators clears all authenticators.
func (mgr *AuthManager) ClearAuthenticators() {
	mgr.authenticators = make([]Authenticator, 0)
}

// Authenticate authenticates the connection with the startup message.
func (mgr *AuthManager) Authenticate(conn Conn, startup *message.Startup) (bool, error) {
	if len(mgr.authenticators) == 0 {
		return true, nil
	}
	var authErr error
	for _, authenticator := range mgr.authenticators {
		ok, err := authenticator.Authenticate(conn, startup)
		if ok {
			return true, nil
		}
		if err != nil {
			authErr = errors.Join(authErr, err)
		}
	}
	return false, authErr
}
