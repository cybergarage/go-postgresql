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

package protocol

import (
	"github.com/cybergarage/go-authenticator/auth/tls"
)

// Config represents a server configuration.
type Config interface {
	tls.CertConfig

	// SetProuctName sets a product name to the configuration.
	SetProductName(v string)
	// SetProductVersion sets a product version to the configuration.
	SetProductVersion(v string)
	// ProductName returns the product name from the configuration.
	ProductName() string
	// ProductVersion returns the product version from the configuration.
	ProductVersion() string
	// SetAddress sets a listen address to the configuration.
	SetAddress(addr string)
	// SetPort sets a listen port to the configuration.
	SetPort(port int)
	// Address returns a listen address from the configuration.
	Address() string
	// Port returns a listen port from the configuration.
	Port() int
}
