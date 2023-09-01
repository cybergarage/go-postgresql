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

// OID represents a PostgreSQL object ID.
type OID int32

const (
	Bool                               OID = 16
	Bytea                              OID = 17
	Char                               OID = 18
	Name                               OID = 19
	Int8                               OID = 20
	Int2                               OID = 21
	Int2vector                         OID = 22
	Int4                               OID = 23
	Regproc                            OID = 24
	Text                               OID = 25
	Oid                                OID = 26
	Tid                                OID = 27
	Xid                                OID = 28
	Cid                                OID = 29
	Oidvector                          OID = 30
	PgType                             OID = 71
	PgAttribute                        OID = 75
	PgProc                             OID = 81
	PgClass                            OID = 83
	JSON                               OID = 114
	XML                                OID = 142
	PgNodeTree                         OID = 194
	PgNdistinct                        OID = 3361
	PgDependencies                     OID = 3402
	PgMcvList                          OID = 5017
	PgDdlCommand                       OID = 32
	Xid8                               OID = 5069
	Point                              OID = 600
	Lseg                               OID = 601
	Path                               OID = 602
	Box                                OID = 603
	Polygon                            OID = 604
	Line                               OID = 628
	Float4                             OID = 700
	Float8                             OID = 701
	Unknown                            OID = 705
	Circle                             OID = 718
	Money                              OID = 790
	Macaddr                            OID = 829
	Inet                               OID = 869
	Cidr                               OID = 650
	Macaddr8                           OID = 774
	ACLitem                            OID = 1033
	Bpchar                             OID = 1042
	Varchar                            OID = 1043
	Date                               OID = 1082
	Time                               OID = 1083
	Timestamp                          OID = 1114
	Timestamptz                        OID = 1184
	Interval                           OID = 1186
	Timetz                             OID = 1266
	Bit                                OID = 1560
	Varbit                             OID = 1562
	Numeric                            OID = 1700
	Refcursor                          OID = 1790
	Regprocedure                       OID = 2202
	Regoper                            OID = 2203
	Regoperator                        OID = 2204
	Regclass                           OID = 2205
	Regcollation                       OID = 4191
	Regtype                            OID = 2206
	Regrole                            OID = 4096
	Regnamespace                       OID = 4089
	UUID                               OID = 2950
	PgLsn                              OID = 3220
	TSvector                           OID = 3614
	Gtsvector                          OID = 3642
	TSquery                            OID = 3615
	Regconfig                          OID = 3734
	Regdictionary                      OID = 3769
	JSONb                              OID = 3802
	JSONpath                           OID = 4072
	TxidSnapshot                       OID = 2970
	PgSnapshot                         OID = 5038
	Int4range                          OID = 3904
	Numrange                           OID = 3906
	TSrange                            OID = 3908
	TStzrange                          OID = 3910
	Daterange                          OID = 3912
	Int8range                          OID = 3926
	Int4multirange                     OID = 4451
	Nummultirange                      OID = 4532
	TSmultirange                       OID = 4533
	TStzmultirange                     OID = 4534
	Datemultirange                     OID = 4535
	Int8multirange                     OID = 4536
	Record                             OID = 2249
	Cstring                            OID = 2275
	Any                                OID = 2276
	Anyarray                           OID = 2277
	Void                               OID = 2278
	Trigger                            OID = 2279
	EventTrigger                       OID = 3838
	LanguageHandler                    OID = 2280
	Internal                           OID = 2281
	Anyelement                         OID = 2283
	Anynonarray                        OID = 2776
	Anyenum                            OID = 3500
	FdwHandler                         OID = 3115
	IndexAmHandler                     OID = 325
	TSmHandler                         OID = 3310
	TableAmHandler                     OID = 269
	Anyrange                           OID = 3831
	Anycompatible                      OID = 5077
	Anycompatiblearray                 OID = 5078
	Anycompatiblenonarray              OID = 5079
	Anycompatiblerange                 OID = 5080
	Anymultirange                      OID = 4537
	Anycompatiblemultirange            OID = 4538
	PgBrinBloomSummary                 OID = 4600
	PgBrinMinmaxMultiSummary           OID = 4601
	PgAttrdef                          OID = 12001
	PgConstraint                       OID = 12003
	PgInherits                         OID = 12005
	PgIndex                            OID = 12007
	PgOperator                         OID = 12009
	PgOpfamily                         OID = 12011
	PgOpclass                          OID = 12013
	PgAm                               OID = 12015
	PgAmop                             OID = 12017
	PgAmproc                           OID = 12019
	PgLanguage                         OID = 12021
	PgLargeobjectMetadata              OID = 12023
	PgLargeobject                      OID = 12025
	PgAggregate                        OID = 12027
	PgStatistic                        OID = 12029
	PgStatisticExt                     OID = 12031
	PgStatisticExtData                 OID = 12033
	PgRewrite                          OID = 12035
	PgTrigger                          OID = 12037
	PgEventTrigger                     OID = 12039
	PgDescription                      OID = 12041
	PgCast                             OID = 12043
	PgEnum                             OID = 12045
	PgNamespace                        OID = 12047
	PgConversion                       OID = 12049
	PgDepend                           OID = 12051
	PgDatabase                         OID = 1248
	PgDBRoleSetting                    OID = 12054
	PgTablespace                       OID = 12056
	PgAuthid                           OID = 2842
	PgAuthMembers                      OID = 2843
	PgShdepend                         OID = 12060
	PgShdescription                    OID = 12062
	PgTSConfig                         OID = 12064
	PgTSConfigMap                      OID = 12066
	PgTSDict                           OID = 12068
	PgTSParser                         OID = 12070
	PgTSTemplate                       OID = 12072
	PgExtension                        OID = 12074
	PgForeignDataWrapper               OID = 12076
	PgForeignServer                    OID = 12078
	PgUserMapping                      OID = 12080
	PgForeignTable                     OID = 12082
	PgPolicy                           OID = 12084
	PgReplicationOrigin                OID = 12086
	PgDefaultACL                       OID = 12088
	PgInitPrivs                        OID = 12090
	PgSeclabel                         OID = 12092
	PgShseclabel                       OID = 4066
	PgCollation                        OID = 12095
	PgPartitionedTable                 OID = 12097
	PgRange                            OID = 12099
	PgTransform                        OID = 12101
	PgSequence                         OID = 12103
	PgPublication                      OID = 12105
	PgPublicationRel                   OID = 12107
	PgSubscription                     OID = 6101
	PgSubscriptionRel                  OID = 12110
	PgRoles                            OID = 12219
	PgShadow                           OID = 12224
	PgGroup                            OID = 12229
	PgUser                             OID = 12233
	PgPolicies                         OID = 12237
	PgRules                            OID = 12242
	PgViews                            OID = 12247
	PgTables                           OID = 12252
	PgMatviews                         OID = 12257
	PgIndexes                          OID = 12262
	PgSequences                        OID = 12267
	PgStats                            OID = 12272
	PgStatsExt                         OID = 12277
	PgStatsExtExprs                    OID = 12282
	PgPublicationTables                OID = 12287
	PgLocks                            OID = 12292
	PgCursors                          OID = 12296
	PgAvailableExtensions              OID = 12300
	PgAvailableExtensionVersions       OID = 12304
	PgPreparedXacts                    OID = 12309
	PgPreparedStatements               OID = 12314
	PgSeclabels                        OID = 12318
	PgSettings                         OID = 12323
	PgFileSettings                     OID = 12329
	PgHbaFileRules                     OID = 12333
	PgTimezoneAbbrevs                  OID = 12337
	PgTimezoneNames                    OID = 12341
	PgConfig                           OID = 12345
	PgShmemAllocations                 OID = 12349
	PgBackendMemoryContexts            OID = 12353
	PgStatAllTables                    OID = 12357
	PgStatXactAllTables                OID = 12362
	PgStatSysTables                    OID = 12367
	PgStatXactSysTables                OID = 12372
	PgStatUserTables                   OID = 12376
	PgStatXactUserTables               OID = 12381
	PgStatioAllTables                  OID = 12385
	PgStatioSysTables                  OID = 12390
	PgStatioUserTables                 OID = 12394
	PgStatAllIndexes                   OID = 12398
	PgStatSysIndexes                   OID = 12403
	PgStatUserIndexes                  OID = 12407
	PgStatioAllIndexes                 OID = 12411
	PgStatioSysIndexes                 OID = 12416
	PgStatioUserIndexes                OID = 12420
	PgStatioAllSequences               OID = 12424
	PgStatioSysSequences               OID = 12429
	PgStatioUserSequences              OID = 12433
	PgStatActivity                     OID = 12437
	PgStatReplication                  OID = 12442
	PgStatSlru                         OID = 12447
	PgStatWalReceiver                  OID = 12451
	PgStatSubscription                 OID = 12455
	PgStatSsl                          OID = 12460
	PgStatGssapi                       OID = 12464
	PgReplicationSlots                 OID = 12468
	PgStatReplicationSlots             OID = 12473
	PgStatDatabase                     OID = 12477
	PgStatDatabaseConflicts            OID = 12482
	PgStatUserFunctions                OID = 12486
	PgStatXactUserFunctions            OID = 12491
	PgStatArchiver                     OID = 12496
	PgStatBgwriter                     OID = 12500
	PgStatWal                          OID = 12504
	PgStatProgressAnalyze              OID = 12508
	PgStatProgressVacuum               OID = 12513
	PgStatProgressCluster              OID = 12518
	PgStatProgressCreateIndex          OID = 12523
	PgStatProgressBasebackup           OID = 12528
	PgStatProgressCopy                 OID = 12533
	PgUserMappings                     OID = 12538
	PgReplicationOriginStatus          OID = 12543
	CardinalNumber                     OID = 13690
	CharacterData                      OID = 13693
	SQLIdentifier                      OID = 13695
	InformationSchemaCatalogName       OID = 13698
	TimeStamp                          OID = 13701
	YesOrNo                            OID = 13703
	ApplicableRoles                    OID = 13707
	AdministrableRoleAuthorizations    OID = 13712
	Attributes                         OID = 13716
	CharacterSets                      OID = 13721
	CheckConstraintRoutineUsage        OID = 13726
	CheckConstraints                   OID = 13731
	Collations                         OID = 13736
	CollationCharacterSetApplicability OID = 13741
	ColumnColumnUsage                  OID = 13746
	ColumnDomainUsage                  OID = 13751
	ColumnPrivileges                   OID = 13756
	ColumnUdtUsage                     OID = 13761
	Columns                            OID = 13766
	ConstraintColumnUsage              OID = 13771
	ConstraintTableUsage               OID = 13776
	DomainConstraints                  OID = 13781
	DomainUdtUsage                     OID = 13786
	Domains                            OID = 13791
	EnabledRoles                       OID = 13796
	KeyColumnUsage                     OID = 13800
	Parameters                         OID = 13805
	ReferentialConstraints             OID = 13810
	RoleColumnGrants                   OID = 13815
	RoutineColumnUsage                 OID = 13819
	RoutinePrivileges                  OID = 13824
	RoleRoutineGrants                  OID = 13829
	RoutineRoutineUsage                OID = 13833
	RoutineSequenceUsage               OID = 13838
	RoutineTableUsage                  OID = 13843
	Routines                           OID = 13848
	Schemata                           OID = 13853
	Sequences                          OID = 13857
	SQLFeatures                        OID = 13862
	SQLImplementationInfo              OID = 13867
	SQLParts                           OID = 13872
	SQLSizing                          OID = 13877
	TableConstraints                   OID = 13882
	TablePrivileges                    OID = 13887
	RoleTableGrants                    OID = 13892
	Tables                             OID = 13896
	Transforms                         OID = 13901
	TriggeredUpdateColumns             OID = 13906
	Triggers                           OID = 13911
	UdtPrivileges                      OID = 13916
	RoleUdtGrants                      OID = 13921
	UsagePrivileges                    OID = 13925
	RoleUsageGrants                    OID = 13930
	UserDefinedTypes                   OID = 13934
	ViewColumnUsage                    OID = 13939
	ViewRoutineUsage                   OID = 13944
	ViewTableUsage                     OID = 13949
	Views                              OID = 13954
	DataTypePrivileges                 OID = 13959
	ElementTypes                       OID = 13964
	ColumnOptions                      OID = 13974
	ForeignDataWrapperOptions          OID = 13982
	ForeignDataWrappers                OID = 13986
	ForeignServerOptions               OID = 13995
	ForeignServers                     OID = 13999
	ForeignTableOptions                OID = 14008
	ForeignTables                      OID = 14012
	UserMappingOptions                 OID = 14021
	UserMappings                       OID = 14026
)
