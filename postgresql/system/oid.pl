#!/usr/bin/perl
# Copyright (C) 2022 Satoshi Konno. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

print<<HEADER;
// Copyright (C) 2022 Satoshi Konno. All rights reserved.
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
HEADER

my $pg_type_file = "res/pg_type.csv";
open(IN, $pg_type_file) or die "Failed to open $pg_type_file: $!";
$line_no = 0;
while(<IN>){
  $line_no++;
  if ($line_no <= 1) {
    next;
  }
  chomp($_);
  my @row = split(/,/, $_, -1);
  my $name = $row[0];
  my $val = $row[1];
  if ($name =~ /^_/){
    next;
  }
  # UpperCase
  # $name =~ s/([a-z])/\u$1/g;
  # CamelCase
  $name =~ s/_([a-z])/\u$1/g;
  $name =~ s/^([a-z])/\u$1/;
  # check nil rows
  if (length($name)<=0 || length($val)<=0) {
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

  print $name . " OID = " . $val . "\n";
  # END: yd8c6549bwf9
  }
  close(IN);

print<<FOTTER;
)
FOTTER