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

const (
	// SystemDatabaseName represents a system database name.
	SystemDatabaseName = "postgres"
	// SystemSchemaName represents a system schema name.
	SystemSchemaName = "pg_catalog"
	// SystemInformationSchema represents a system information schema name.
	SystemInformationSchema = "information_schema"
)

// SystemSchemaNames represents system schema names.
var SystemSchemaNames = []string{
	SystemSchemaName,
	SystemInformationSchema,
}

const (
	// InformationSchemaColumns represents a system information schema cloumns tables name.
	InformationSchemaColumns = SystemInformationSchema + "." + "columns"
)
