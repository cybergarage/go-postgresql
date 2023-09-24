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

const (
	InformationSchemaColumnsTableCatalog           = "table_catalog"
	InformationSchemaColumnsTableSchema            = "table_schema"
	InformationSchemaColumnsTableName              = "table_name"
	InformationSchemaColumnsColumnName             = "column_name"
	InformationSchemaColumnsOrdinalPosition        = "ordinal_position"
	InformationSchemaColumnsColumnDefault          = "column_default"
	InformationSchemaColumnsIsNullable             = "is_nullable"
	InformationSchemaColumnsDataType               = "data_type"
	InformationSchemaColumnsCharacterMaximumLength = "character_maximum_length"
	InformationSchemaColumnsCharacterOctetLength   = "character_octet_length"
	InformationSchemaColumnsNumericPrecision       = "numeric_precision"
	InformationSchemaColumnsNumericPrecisionRadix  = "numeric_precision_radix"
	InformationSchemaColumnsNumericScale           = "numeric_scale"
	InformationSchemaColumnsDatetimePrecision      = "datetime_precision"
	InformationSchemaColumnsIntervalType           = "interval_type"
	InformationSchemaColumnsIntervalPrecision      = "interval_precision"
	InformationSchemaColumnsCharacterSetCatalog    = "character_set_catalog"
	InformationSchemaColumnsCharacterSetSchema     = "character_set_schema"
	InformationSchemaColumnsCharacterSetName       = "character_set_name"
	InformationSchemaColumnsCollationCatalog       = "collation_catalog"
	InformationSchemaColumnsCollationSchema        = "collation_schema"
	InformationSchemaColumnsCollationName          = "collation_name"
	InformationSchemaColumnsDomainCatalog          = "domain_catalog"
	InformationSchemaColumnsDomainSchema           = "domain_schema"
	InformationSchemaColumnsDomainName             = "domain_name"
	InformationSchemaColumnsUdtCatalog             = "udt_catalog"
	InformationSchemaColumnsUdtSchema              = "udt_schema"
	InformationSchemaColumnsUdtName                = "udt_name"
	InformationSchemaColumnsScopeCatalog           = "scope_catalog"
	InformationSchemaColumnsScopeSchema            = "scope_schema"
	InformationSchemaColumnsScopeName              = "scope_name"
	InformationSchemaColumnsMaximumCardinality     = "maximum_cardinality"
	InformationSchemaColumnsDtdIdentifier          = "dtd_identifier"
	InformationSchemaColumnsIsSelfReferencing      = "is_self_referencing"
	InformationSchemaColumnsIsIdentity             = "is_identity"
	InformationSchemaColumnsIdentityGeneration     = "identity_generation"
	InformationSchemaColumnsIdentityStart          = "identity_start"
	InformationSchemaColumnsIdentityIncrement      = "identity_increment"
	InformationSchemaColumnsIdentityMaximum        = "identity_maximum"
	InformationSchemaColumnsIdentityMinimum        = "identity_minimum"
	InformationSchemaColumnsIdentityCycle          = "identity_cycle"
	InformationSchemaColumnsIsGenerated            = "is_generated"
	InformationSchemaColumnsGenerationExpression   = "generation_expression"
	InformationSchemaColumnsIsUpdatable            = "is_updatable"
)
