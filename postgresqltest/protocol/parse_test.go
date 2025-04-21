// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package protocol

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-postgresql/postgresql/protocol"
	"github.com/cybergarage/go-postgresql/postgresql/query"
)

func TestParsePacket(t *testing.T) {
	type expected struct {
	}
	for _, test := range []struct {
		name        string
		chkDescribe bool
		expected
	}{
		{
			"data/sysbench-parse-001.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-002.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-003.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-004.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-005.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-006.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-007.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-008.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-009.hex",
			false,
			expected{},
		},
		{
			"data/sysbench-parse-010.hex",
			false,
			expected{},
		},
		{
			"data/go-pq-parse-001.hex",
			true,
			expected{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			testData, err := testEmbedPacketFiles.ReadFile(test.name)
			if err != nil {
				t.Error(err)
				return
			}
			testBytes, err := hexdump.NewBytesWithHexdumpBytes(testData)
			if err != nil {
				t.Error(err)
				return
			}

			reader := protocol.NewMessageReaderWith(bytes.NewReader(testBytes))

			pkt, err := protocol.NewParseWithReader(reader)
			if err != nil {
				t.Error(err)
				return
			}

			parser := query.NewParser()
			_, err = parser.ParseString(pkt.Query)
			if err != nil {
				t.Error(err)
			}

			if !test.chkDescribe {
				return
			}

			_, err = protocol.NewDescribeWithReader(reader)
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
}
