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

package data

// PostgreSQL: Documentation: 16: Chapter 8. Data Types
// https://www.postgresql.org/docs/16/datatype.html

// DataType represents a PostgreSQL data type using the object ID.
type DataType int16

const (
	BoolType            DataType = 16
	ByteaType           DataType = 17
	CharType            DataType = 18
	NameType            DataType = 19
	Int8Type            DataType = 20
	Int2Type            DataType = 21
	Int2VectorType      DataType = 22
	Int4Type            DataType = 23
	RegProcType         DataType = 24
	TextType            DataType = 25
	OidType             DataType = 26
	TidType             DataType = 27
	XidType             DataType = 28
	CidType             DataType = 29
	XMLType             DataType = 142
	PointType           DataType = 600
	LsegType            DataType = 601
	PathType            DataType = 602
	BoxType             DataType = 603
	PolygonType         DataType = 604
	LineType            DataType = 628
	LineArrayType       DataType = 629
	CircleType          DataType = 718
	CircleArrayType     DataType = 719
	MacaddrType         DataType = 829
	Macaddr8Type        DataType = 774
	InetType            DataType = 869
	InetArrayType       DataType = 1040
	CidrType            DataType = 650
	CidrArrayType       DataType = 651
	Float4Type          DataType = 700
	Float8Type          DataType = 701
	UnknownType         DataType = 705
	AbstimeType         DataType = 702
	ReltimeType         DataType = 703
	TintervalType       DataType = 704
	PolygonArrayType    DataType = 628
	OidvectorType       DataType = 30
	BpcharType          DataType = 1042
	VarcharType         DataType = 1043
	DateType            DataType = 1082
	TimeType            DataType = 1083
	TimestampType       DataType = 1114
	TimestampTzType     DataType = 1184
	IntervalType        DataType = 1186
	TimeTzType          DataType = 1266
	BitType             DataType = 1560
	VarbitType          DataType = 1562
	NumericType         DataType = 1700
	RefcursorType       DataType = 1790
	RegprocedureType    DataType = 2202
	RegoperType         DataType = 2203
	RegoperatorType     DataType = 2204
	RegclassType        DataType = 2205
	RegtypeType         DataType = 2206
	RecordType          DataType = 2249
	CstringType         DataType = 2275
	AnyType             DataType = 2276
	AnyarrayType        DataType = 2277
	VoidType            DataType = 2278
	TriggerType         DataType = 2279
	LanguageHandlerType DataType = 2280
	InternalType        DataType = 2281
	OpaqueType          DataType = 2282
	AnyelementType      DataType = 2283
	AnynonarrayType     DataType = 2776
	GeometryType        DataType = 3614
	GeograpyType        DataType = 4326
)
