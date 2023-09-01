// Copyright (C) 2022 The go-postgresql Authors. All rights reserved.
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
	_ "embed"
)

// Catalogs is a map of catalog names and catalog bytes.
// PostgreSQL: Documentation: System Catalogs
// https://www.postgresql.org/docs/current/catalogs.html
var Catalogs = map[string][]byte{
	"PgType": pgType,
}

//go:embed res/pg_type.csv
var pgType []byte
