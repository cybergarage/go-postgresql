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

type DataType struct {
	name string
	oid  OID
	size int
}

// DataType represents a PostgreSQL data type.
var dataTypes = map[OID]*DataType{}

// GetDataType returns a data type for the specified OID.
func GetDataType(oid OID) (*DataType, error) {
	dt, ok := dataTypes[oid]
	if !ok {
		return nil, newDataTypeNotFound(oid)
	}
	return dt, nil
}

func newDataType(name string, oid OID, size int) *DataType {
	dt := &DataType{
		name: name,
		oid:  oid,
		size: size,
	}
	return dt
}

// Name returns the data type name.
func (dt *DataType) Name() string {
	return dt.name
}

// OID returns the data type OID.
func (dt *DataType) OID() OID {
	return dt.oid
}

// Size returns the data type size.
func (dt *DataType) Size() int {
	return dt.size
}

// FormatCodeFrom returns a format code from the specified data type.
func (dt *DataType) FormatCode() FormatCode {
	switch dt.oid {
	case Name:
		return TextFormat
	case Text:
		return TextFormat
	case Varchar:
		return TextFormat
	case Int2:
		return TextFormat
	case Int4:
		return TextFormat
	case Int8:
		return TextFormat
	}
	return BinaryFormat
}

func init() { //nolint: gochecknoinits, maintidx
	dataTypes[16] = newDataType("Bool", 16, 1)
	dataTypes[17] = newDataType("Bytea", 17, -1)
	dataTypes[18] = newDataType("Char", 18, 1)
	dataTypes[19] = newDataType("Name", 19, 64)
	dataTypes[20] = newDataType("Int8", 20, 8)
	dataTypes[21] = newDataType("Int2", 21, 2)
	dataTypes[22] = newDataType("Int2vector", 22, -1)
	dataTypes[23] = newDataType("Int4", 23, 4)
	dataTypes[24] = newDataType("Regproc", 24, 4)
	dataTypes[25] = newDataType("Text", 25, -1)
	dataTypes[26] = newDataType("Oid", 26, 4)
	dataTypes[27] = newDataType("Tid", 27, 6)
	dataTypes[28] = newDataType("Xid", 28, 4)
	dataTypes[29] = newDataType("Cid", 29, 4)
	dataTypes[30] = newDataType("Oidvector", 30, -1)
	dataTypes[71] = newDataType("PgType", 71, -1)
	dataTypes[75] = newDataType("PgAttribute", 75, -1)
	dataTypes[81] = newDataType("PgProc", 81, -1)
	dataTypes[83] = newDataType("PgClass", 83, -1)
	dataTypes[114] = newDataType("JSON", 114, -1)
	dataTypes[142] = newDataType("XML", 142, -1)
	dataTypes[194] = newDataType("PgNodeTree", 194, -1)
	dataTypes[3361] = newDataType("PgNdistinct", 3361, -1)
	dataTypes[3402] = newDataType("PgDependencies", 3402, -1)
	dataTypes[5017] = newDataType("PgMcvList", 5017, -1)
	dataTypes[32] = newDataType("PgDdlCommand", 32, 8)
	dataTypes[5069] = newDataType("Xid8", 5069, 8)
	dataTypes[600] = newDataType("Point", 600, 16)
	dataTypes[601] = newDataType("Lseg", 601, 32)
	dataTypes[602] = newDataType("Path", 602, -1)
	dataTypes[603] = newDataType("Box", 603, 32)
	dataTypes[604] = newDataType("Polygon", 604, -1)
	dataTypes[628] = newDataType("Line", 628, 24)
	dataTypes[700] = newDataType("Float4", 700, 4)
	dataTypes[701] = newDataType("Float8", 701, 8)
	dataTypes[705] = newDataType("Unknown", 705, -2)
	dataTypes[718] = newDataType("Circle", 718, 24)
	dataTypes[790] = newDataType("Money", 790, 8)
	dataTypes[829] = newDataType("Macaddr", 829, 6)
	dataTypes[869] = newDataType("Inet", 869, -1)
	dataTypes[650] = newDataType("Cidr", 650, -1)
	dataTypes[774] = newDataType("Macaddr8", 774, 8)
	dataTypes[1033] = newDataType("ACLitem", 1033, 12)
	dataTypes[1042] = newDataType("Bpchar", 1042, -1)
	dataTypes[1043] = newDataType("Varchar", 1043, -1)
	dataTypes[1082] = newDataType("Date", 1082, 4)
	dataTypes[1083] = newDataType("Time", 1083, 8)
	dataTypes[1114] = newDataType("Timestamp", 1114, 8)
	dataTypes[1184] = newDataType("Timestamptz", 1184, 8)
	dataTypes[1186] = newDataType("Interval", 1186, 16)
	dataTypes[1266] = newDataType("Timetz", 1266, 12)
	dataTypes[1560] = newDataType("Bit", 1560, -1)
	dataTypes[1562] = newDataType("Varbit", 1562, -1)
	dataTypes[1700] = newDataType("Numeric", 1700, -1)
	dataTypes[1790] = newDataType("Refcursor", 1790, -1)
	dataTypes[2202] = newDataType("Regprocedure", 2202, 4)
	dataTypes[2203] = newDataType("Regoper", 2203, 4)
	dataTypes[2204] = newDataType("Regoperator", 2204, 4)
	dataTypes[2205] = newDataType("Regclass", 2205, 4)
	dataTypes[4191] = newDataType("Regcollation", 4191, 4)
	dataTypes[2206] = newDataType("Regtype", 2206, 4)
	dataTypes[4096] = newDataType("Regrole", 4096, 4)
	dataTypes[4089] = newDataType("Regnamespace", 4089, 4)
	dataTypes[2950] = newDataType("UUID", 2950, 16)
	dataTypes[3220] = newDataType("PgLsn", 3220, 8)
	dataTypes[3614] = newDataType("TSvector", 3614, -1)
	dataTypes[3642] = newDataType("Gtsvector", 3642, -1)
	dataTypes[3615] = newDataType("TSquery", 3615, -1)
	dataTypes[3734] = newDataType("Regconfig", 3734, 4)
	dataTypes[3769] = newDataType("Regdictionary", 3769, 4)
	dataTypes[3802] = newDataType("JSONb", 3802, -1)
	dataTypes[4072] = newDataType("JSONpath", 4072, -1)
	dataTypes[2970] = newDataType("TxidSnapshot", 2970, -1)
	dataTypes[5038] = newDataType("PgSnapshot", 5038, -1)
	dataTypes[3904] = newDataType("Int4range", 3904, -1)
	dataTypes[3906] = newDataType("Numrange", 3906, -1)
	dataTypes[3908] = newDataType("TSrange", 3908, -1)
	dataTypes[3910] = newDataType("TStzrange", 3910, -1)
	dataTypes[3912] = newDataType("Daterange", 3912, -1)
	dataTypes[3926] = newDataType("Int8range", 3926, -1)
	dataTypes[4451] = newDataType("Int4multirange", 4451, -1)
	dataTypes[4532] = newDataType("Nummultirange", 4532, -1)
	dataTypes[4533] = newDataType("TSmultirange", 4533, -1)
	dataTypes[4534] = newDataType("TStzmultirange", 4534, -1)
	dataTypes[4535] = newDataType("Datemultirange", 4535, -1)
	dataTypes[4536] = newDataType("Int8multirange", 4536, -1)
	dataTypes[2249] = newDataType("Record", 2249, -1)
	dataTypes[2275] = newDataType("Cstring", 2275, -2)
	dataTypes[2276] = newDataType("Any", 2276, 4)
	dataTypes[2277] = newDataType("Anyarray", 2277, -1)
	dataTypes[2278] = newDataType("Void", 2278, 4)
	dataTypes[2279] = newDataType("Trigger", 2279, 4)
	dataTypes[3838] = newDataType("EventTrigger", 3838, 4)
	dataTypes[2280] = newDataType("LanguageHandler", 2280, 4)
	dataTypes[2281] = newDataType("Internal", 2281, 8)
	dataTypes[2283] = newDataType("Anyelement", 2283, 4)
	dataTypes[2776] = newDataType("Anynonarray", 2776, 4)
	dataTypes[3500] = newDataType("Anyenum", 3500, 4)
	dataTypes[3115] = newDataType("FdwHandler", 3115, 4)
	dataTypes[325] = newDataType("IndexAmHandler", 325, 4)
	dataTypes[3310] = newDataType("TSmHandler", 3310, 4)
	dataTypes[269] = newDataType("TableAmHandler", 269, 4)
	dataTypes[3831] = newDataType("Anyrange", 3831, -1)
	dataTypes[5077] = newDataType("Anycompatible", 5077, 4)
	dataTypes[5078] = newDataType("Anycompatiblearray", 5078, -1)
	dataTypes[5079] = newDataType("Anycompatiblenonarray", 5079, 4)
	dataTypes[5080] = newDataType("Anycompatiblerange", 5080, -1)
	dataTypes[4537] = newDataType("Anymultirange", 4537, -1)
	dataTypes[4538] = newDataType("Anycompatiblemultirange", 4538, -1)
	dataTypes[4600] = newDataType("PgBrinBloomSummary", 4600, -1)
	dataTypes[4601] = newDataType("PgBrinMinmaxMultiSummary", 4601, -1)
	dataTypes[12001] = newDataType("PgAttrdef", 12001, -1)
	dataTypes[12003] = newDataType("PgConstraint", 12003, -1)
	dataTypes[12005] = newDataType("PgInherits", 12005, -1)
	dataTypes[12007] = newDataType("PgIndex", 12007, -1)
	dataTypes[12009] = newDataType("PgOperator", 12009, -1)
	dataTypes[12011] = newDataType("PgOpfamily", 12011, -1)
	dataTypes[12013] = newDataType("PgOpclass", 12013, -1)
	dataTypes[12015] = newDataType("PgAm", 12015, -1)
	dataTypes[12017] = newDataType("PgAmop", 12017, -1)
	dataTypes[12019] = newDataType("PgAmproc", 12019, -1)
	dataTypes[12021] = newDataType("PgLanguage", 12021, -1)
	dataTypes[12023] = newDataType("PgLargeobjectMetadata", 12023, -1)
	dataTypes[12025] = newDataType("PgLargeobject", 12025, -1)
	dataTypes[12027] = newDataType("PgAggregate", 12027, -1)
	dataTypes[12029] = newDataType("PgStatistic", 12029, -1)
	dataTypes[12031] = newDataType("PgStatisticExt", 12031, -1)
	dataTypes[12033] = newDataType("PgStatisticExtData", 12033, -1)
	dataTypes[12035] = newDataType("PgRewrite", 12035, -1)
	dataTypes[12037] = newDataType("PgTrigger", 12037, -1)
	dataTypes[12039] = newDataType("PgEventTrigger", 12039, -1)
	dataTypes[12041] = newDataType("PgDescription", 12041, -1)
	dataTypes[12043] = newDataType("PgCast", 12043, -1)
	dataTypes[12045] = newDataType("PgEnum", 12045, -1)
	dataTypes[12047] = newDataType("PgNamespace", 12047, -1)
	dataTypes[12049] = newDataType("PgConversion", 12049, -1)
	dataTypes[12051] = newDataType("PgDepend", 12051, -1)
	dataTypes[1248] = newDataType("PgDatabase", 1248, -1)
	dataTypes[12054] = newDataType("PgDBRoleSetting", 12054, -1)
	dataTypes[12056] = newDataType("PgTablespace", 12056, -1)
	dataTypes[2842] = newDataType("PgAuthid", 2842, -1)
	dataTypes[2843] = newDataType("PgAuthMembers", 2843, -1)
	dataTypes[12060] = newDataType("PgShdepend", 12060, -1)
	dataTypes[12062] = newDataType("PgShdescription", 12062, -1)
	dataTypes[12064] = newDataType("PgTSConfig", 12064, -1)
	dataTypes[12066] = newDataType("PgTSConfigMap", 12066, -1)
	dataTypes[12068] = newDataType("PgTSDict", 12068, -1)
	dataTypes[12070] = newDataType("PgTSParser", 12070, -1)
	dataTypes[12072] = newDataType("PgTSTemplate", 12072, -1)
	dataTypes[12074] = newDataType("PgExtension", 12074, -1)
	dataTypes[12076] = newDataType("PgForeignDataWrapper", 12076, -1)
	dataTypes[12078] = newDataType("PgForeignServer", 12078, -1)
	dataTypes[12080] = newDataType("PgUserMapping", 12080, -1)
	dataTypes[12082] = newDataType("PgForeignTable", 12082, -1)
	dataTypes[12084] = newDataType("PgPolicy", 12084, -1)
	dataTypes[12086] = newDataType("PgReplicationOrigin", 12086, -1)
	dataTypes[12088] = newDataType("PgDefaultACL", 12088, -1)
	dataTypes[12090] = newDataType("PgInitPrivs", 12090, -1)
	dataTypes[12092] = newDataType("PgSeclabel", 12092, -1)
	dataTypes[4066] = newDataType("PgShseclabel", 4066, -1)
	dataTypes[12095] = newDataType("PgCollation", 12095, -1)
	dataTypes[12097] = newDataType("PgPartitionedTable", 12097, -1)
	dataTypes[12099] = newDataType("PgRange", 12099, -1)
	dataTypes[12101] = newDataType("PgTransform", 12101, -1)
	dataTypes[12103] = newDataType("PgSequence", 12103, -1)
	dataTypes[12105] = newDataType("PgPublication", 12105, -1)
	dataTypes[12107] = newDataType("PgPublicationRel", 12107, -1)
	dataTypes[6101] = newDataType("PgSubscription", 6101, -1)
	dataTypes[12110] = newDataType("PgSubscriptionRel", 12110, -1)
	dataTypes[12219] = newDataType("PgRoles", 12219, -1)
	dataTypes[12224] = newDataType("PgShadow", 12224, -1)
	dataTypes[12229] = newDataType("PgGroup", 12229, -1)
	dataTypes[12233] = newDataType("PgUser", 12233, -1)
	dataTypes[12237] = newDataType("PgPolicies", 12237, -1)
	dataTypes[12242] = newDataType("PgRules", 12242, -1)
	dataTypes[12247] = newDataType("PgViews", 12247, -1)
	dataTypes[12252] = newDataType("PgTables", 12252, -1)
	dataTypes[12257] = newDataType("PgMatviews", 12257, -1)
	dataTypes[12262] = newDataType("PgIndexes", 12262, -1)
	dataTypes[12267] = newDataType("PgSequences", 12267, -1)
	dataTypes[12272] = newDataType("PgStats", 12272, -1)
	dataTypes[12277] = newDataType("PgStatsExt", 12277, -1)
	dataTypes[12282] = newDataType("PgStatsExtExprs", 12282, -1)
	dataTypes[12287] = newDataType("PgPublicationTables", 12287, -1)
	dataTypes[12292] = newDataType("PgLocks", 12292, -1)
	dataTypes[12296] = newDataType("PgCursors", 12296, -1)
	dataTypes[12300] = newDataType("PgAvailableExtensions", 12300, -1)
	dataTypes[12304] = newDataType("PgAvailableExtensionVersions", 12304, -1)
	dataTypes[12309] = newDataType("PgPreparedXacts", 12309, -1)
	dataTypes[12314] = newDataType("PgPreparedStatements", 12314, -1)
	dataTypes[12318] = newDataType("PgSeclabels", 12318, -1)
	dataTypes[12323] = newDataType("PgSettings", 12323, -1)
	dataTypes[12329] = newDataType("PgFileSettings", 12329, -1)
	dataTypes[12333] = newDataType("PgHbaFileRules", 12333, -1)
	dataTypes[12337] = newDataType("PgTimezoneAbbrevs", 12337, -1)
	dataTypes[12341] = newDataType("PgTimezoneNames", 12341, -1)
	dataTypes[12345] = newDataType("PgConfig", 12345, -1)
	dataTypes[12349] = newDataType("PgShmemAllocations", 12349, -1)
	dataTypes[12353] = newDataType("PgBackendMemoryContexts", 12353, -1)
	dataTypes[12357] = newDataType("PgStatAllTables", 12357, -1)
	dataTypes[12362] = newDataType("PgStatXactAllTables", 12362, -1)
	dataTypes[12367] = newDataType("PgStatSysTables", 12367, -1)
	dataTypes[12372] = newDataType("PgStatXactSysTables", 12372, -1)
	dataTypes[12376] = newDataType("PgStatUserTables", 12376, -1)
	dataTypes[12381] = newDataType("PgStatXactUserTables", 12381, -1)
	dataTypes[12385] = newDataType("PgStatioAllTables", 12385, -1)
	dataTypes[12390] = newDataType("PgStatioSysTables", 12390, -1)
	dataTypes[12394] = newDataType("PgStatioUserTables", 12394, -1)
	dataTypes[12398] = newDataType("PgStatAllIndexes", 12398, -1)
	dataTypes[12403] = newDataType("PgStatSysIndexes", 12403, -1)
	dataTypes[12407] = newDataType("PgStatUserIndexes", 12407, -1)
	dataTypes[12411] = newDataType("PgStatioAllIndexes", 12411, -1)
	dataTypes[12416] = newDataType("PgStatioSysIndexes", 12416, -1)
	dataTypes[12420] = newDataType("PgStatioUserIndexes", 12420, -1)
	dataTypes[12424] = newDataType("PgStatioAllSequences", 12424, -1)
	dataTypes[12429] = newDataType("PgStatioSysSequences", 12429, -1)
	dataTypes[12433] = newDataType("PgStatioUserSequences", 12433, -1)
	dataTypes[12437] = newDataType("PgStatActivity", 12437, -1)
	dataTypes[12442] = newDataType("PgStatReplication", 12442, -1)
	dataTypes[12447] = newDataType("PgStatSlru", 12447, -1)
	dataTypes[12451] = newDataType("PgStatWalReceiver", 12451, -1)
	dataTypes[12455] = newDataType("PgStatSubscription", 12455, -1)
	dataTypes[12460] = newDataType("PgStatSsl", 12460, -1)
	dataTypes[12464] = newDataType("PgStatGssapi", 12464, -1)
	dataTypes[12468] = newDataType("PgReplicationSlots", 12468, -1)
	dataTypes[12473] = newDataType("PgStatReplicationSlots", 12473, -1)
	dataTypes[12477] = newDataType("PgStatDatabase", 12477, -1)
	dataTypes[12482] = newDataType("PgStatDatabaseConflicts", 12482, -1)
	dataTypes[12486] = newDataType("PgStatUserFunctions", 12486, -1)
	dataTypes[12491] = newDataType("PgStatXactUserFunctions", 12491, -1)
	dataTypes[12496] = newDataType("PgStatArchiver", 12496, -1)
	dataTypes[12500] = newDataType("PgStatBgwriter", 12500, -1)
	dataTypes[12504] = newDataType("PgStatWal", 12504, -1)
	dataTypes[12508] = newDataType("PgStatProgressAnalyze", 12508, -1)
	dataTypes[12513] = newDataType("PgStatProgressVacuum", 12513, -1)
	dataTypes[12518] = newDataType("PgStatProgressCluster", 12518, -1)
	dataTypes[12523] = newDataType("PgStatProgressCreateIndex", 12523, -1)
	dataTypes[12528] = newDataType("PgStatProgressBasebackup", 12528, -1)
	dataTypes[12533] = newDataType("PgStatProgressCopy", 12533, -1)
	dataTypes[12538] = newDataType("PgUserMappings", 12538, -1)
	dataTypes[12543] = newDataType("PgReplicationOriginStatus", 12543, -1)
	dataTypes[13690] = newDataType("CardinalNumber", 13690, 4)
	dataTypes[13693] = newDataType("CharacterData", 13693, -1)
	dataTypes[13695] = newDataType("SQLIdentifier", 13695, 64)
	dataTypes[13698] = newDataType("InformationSchemaCatalogName", 13698, -1)
	dataTypes[13701] = newDataType("TimeStamp", 13701, 8)
	dataTypes[13703] = newDataType("YesOrNo", 13703, -1)
	dataTypes[13707] = newDataType("ApplicableRoles", 13707, -1)
	dataTypes[13712] = newDataType("AdministrableRoleAuthorizations", 13712, -1)
	dataTypes[13716] = newDataType("Attributes", 13716, -1)
	dataTypes[13721] = newDataType("CharacterSets", 13721, -1)
	dataTypes[13726] = newDataType("CheckConstraintRoutineUsage", 13726, -1)
	dataTypes[13731] = newDataType("CheckConstraints", 13731, -1)
	dataTypes[13736] = newDataType("Collations", 13736, -1)
	dataTypes[13741] = newDataType("CollationCharacterSetApplicability", 13741, -1)
	dataTypes[13746] = newDataType("ColumnColumnUsage", 13746, -1)
	dataTypes[13751] = newDataType("ColumnDomainUsage", 13751, -1)
	dataTypes[13756] = newDataType("ColumnPrivileges", 13756, -1)
	dataTypes[13761] = newDataType("ColumnUdtUsage", 13761, -1)
	dataTypes[13766] = newDataType("Columns", 13766, -1)
	dataTypes[13771] = newDataType("ConstraintColumnUsage", 13771, -1)
	dataTypes[13776] = newDataType("ConstraintTableUsage", 13776, -1)
	dataTypes[13781] = newDataType("DomainConstraints", 13781, -1)
	dataTypes[13786] = newDataType("DomainUdtUsage", 13786, -1)
	dataTypes[13791] = newDataType("Domains", 13791, -1)
	dataTypes[13796] = newDataType("EnabledRoles", 13796, -1)
	dataTypes[13800] = newDataType("KeyColumnUsage", 13800, -1)
	dataTypes[13805] = newDataType("Parameters", 13805, -1)
	dataTypes[13810] = newDataType("ReferentialConstraints", 13810, -1)
	dataTypes[13815] = newDataType("RoleColumnGrants", 13815, -1)
	dataTypes[13819] = newDataType("RoutineColumnUsage", 13819, -1)
	dataTypes[13824] = newDataType("RoutinePrivileges", 13824, -1)
	dataTypes[13829] = newDataType("RoleRoutineGrants", 13829, -1)
	dataTypes[13833] = newDataType("RoutineRoutineUsage", 13833, -1)
	dataTypes[13838] = newDataType("RoutineSequenceUsage", 13838, -1)
	dataTypes[13843] = newDataType("RoutineTableUsage", 13843, -1)
	dataTypes[13848] = newDataType("Routines", 13848, -1)
	dataTypes[13853] = newDataType("Schemata", 13853, -1)
	dataTypes[13857] = newDataType("Sequences", 13857, -1)
	dataTypes[13862] = newDataType("SQLFeatures", 13862, -1)
	dataTypes[13867] = newDataType("SQLImplementationInfo", 13867, -1)
	dataTypes[13872] = newDataType("SQLParts", 13872, -1)
	dataTypes[13877] = newDataType("SQLSizing", 13877, -1)
	dataTypes[13882] = newDataType("TableConstraints", 13882, -1)
	dataTypes[13887] = newDataType("TablePrivileges", 13887, -1)
	dataTypes[13892] = newDataType("RoleTableGrants", 13892, -1)
	dataTypes[13896] = newDataType("Tables", 13896, -1)
	dataTypes[13901] = newDataType("Transforms", 13901, -1)
	dataTypes[13906] = newDataType("TriggeredUpdateColumns", 13906, -1)
	dataTypes[13911] = newDataType("Triggers", 13911, -1)
	dataTypes[13916] = newDataType("UdtPrivileges", 13916, -1)
	dataTypes[13921] = newDataType("RoleUdtGrants", 13921, -1)
	dataTypes[13925] = newDataType("UsagePrivileges", 13925, -1)
	dataTypes[13930] = newDataType("RoleUsageGrants", 13930, -1)
	dataTypes[13934] = newDataType("UserDefinedTypes", 13934, -1)
	dataTypes[13939] = newDataType("ViewColumnUsage", 13939, -1)
	dataTypes[13944] = newDataType("ViewRoutineUsage", 13944, -1)
	dataTypes[13949] = newDataType("ViewTableUsage", 13949, -1)
	dataTypes[13954] = newDataType("Views", 13954, -1)
	dataTypes[13959] = newDataType("DataTypePrivileges", 13959, -1)
	dataTypes[13964] = newDataType("ElementTypes", 13964, -1)
	dataTypes[13974] = newDataType("ColumnOptions", 13974, -1)
	dataTypes[13982] = newDataType("ForeignDataWrapperOptions", 13982, -1)
	dataTypes[13986] = newDataType("ForeignDataWrappers", 13986, -1)
	dataTypes[13995] = newDataType("ForeignServerOptions", 13995, -1)
	dataTypes[13999] = newDataType("ForeignServers", 13999, -1)
	dataTypes[14008] = newDataType("ForeignTableOptions", 14008, -1)
	dataTypes[14012] = newDataType("ForeignTables", 14012, -1)
	dataTypes[14021] = newDataType("UserMappingOptions", 14021, -1)
	dataTypes[14026] = newDataType("UserMappings", 14026, -1)
}
