#!/usr/bin/perl
# Copyright (C) 2022 The go-postgresql Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

print<<HEADER;
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
	case Bytea:
		return BinaryFormat
	}
	return TextFormat
}

func init() { //nolint: maintidx
HEADER

my $pg_type_file = "res/pg_type.csv";
open(IN, $pg_type_file) or die "Failed to open $pg_type_file: $!";
my $line_no = 0;
my $name_idx = -1;
my $oid_idx = -1;
my $typ_len_idx = -1;
while(<IN>){
  $line_no++;
  chomp($_);
  my @row = split(/,/, $_, -1);
  if ($line_no <= 1) {
    for (my $i = 0; $i < scalar(@row); $i++) {
      if ($row[$i] eq "typname") {
        $name_idx = $i;
      }
      if ($row[$i] eq "oid") {
        $oid_idx = $i;
      }
      if ($row[$i] eq "typlen") {
        $typ_len_idx = $i;
      }      
    }
    if ($name_idx < 0) {
      die "Failed to find typname column";
    }
    if ($oid_idx < 0) {
      die "Failed to find oid column";
    }
    if ($typ_len_idx < 0) {
      die "Failed to find typlen column";
    }
    next;
  }
  my $name = $row[$name_idx];
  my $oid = $row[$oid_idx];
  my $size = $row[$typ_len_idx];

  if ($name =~ /^_/){
    next;
  }
  # check nil rows
  if (length($name)<=0 || length($oid)<=0 || length($size)<=0) {
    next;
  }
  # UpperCase
  # $name =~ s/([a-z])/\u$1/g;
  # CamelCase
  $name =~ s/_([a-z])/\u$1/g;
  $name =~ s/^([a-z])/\u$1/;
  # stylecheck
  my @patterns = ('Sql', 'Xml', 'Json', 'Uuid', 'Db', 'Ts', 'Acl');
  my @replacements = ('SQL', 'XML', 'JSON', 'UUID', 'DB', 'TS', 'ACL');
  for (my $i = 0; $i < scalar(@patterns); $i++) {
    if ($name =~ /$patterns[$i]/) {
      $name =~ s/$patterns[$i]/$replacements[$i]/g;
    }
  }
  print "\tdataTypes[" . $oid . "] = newDataType(\"" . $name . "\", " . $oid . ", " . $size . ")\n";  
}
close(IN);

print<<FOTTER;
}
FOTTER