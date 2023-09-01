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

// OID represents a PostgreSQL object ID.
type OID = int32

const (
HEADER

my $pg_type_file = "res/pg_type.csv";
open(IN, $pg_type_file) or die "Failed to open $pg_type_file: $!";
my $line_no = 0;
my $name_idx = -1;
my $oid_idx = -1;
while(<IN>){
  $line_no++;
  my @row = split(/,/, $_, -1);
  if ($line_no <= 1) {
    for (my $i = 0; $i < scalar(@row); $i++) {
      if ($row[$i] eq "typname") {
        $name_idx = $i;
      }
      if ($row[$i] eq "oid") {
        $oid_idx = $i;
      }
    }
    if ($name_idx < 0 || $oid_idx < 0) {
      die "Failed to find typname or oid column";
    }
    next;
  }
  chomp($_);
  my $name = $row[$name_idx];
  my $oid = $row[$oid_idx];
  if ($name =~ /^_/){
    next;
  }
  # UpperCase
  # $name =~ s/([a-z])/\u$1/g;
  # CamelCase
  $name =~ s/_([a-z])/\u$1/g;
  $name =~ s/^([a-z])/\u$1/;
  # check nil rows
  if (length($name)<=0 || length($oid)<=0) {
    next;
  }
  # stylecheck
  my @patterns = ('Sql', 'Xml', 'Json', 'Uuid', 'Db', 'Ts', 'Acl');
  my @replacements = ('SQL', 'XML', 'JSON', 'UUID', 'DB', 'TS', 'ACL');
  for (my $i = 0; $i < scalar(@patterns); $i++) {
    if ($name =~ /$patterns[$i]/) {
      $name =~ s/$patterns[$i]/$replacements[$i]/g;
    }
  }

  print $name . " OID = " . $oid . "\n";
  }
  close(IN);

print<<FOTTER;
)
FOTTER