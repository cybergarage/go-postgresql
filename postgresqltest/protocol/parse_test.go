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
)

func TestParsePacket(t *testing.T) {

	type expected struct {
	}
	for _, test := range []struct {
		name string
		expected
	}{
		{
			"data/sysbench-parse-001.hex",
			expected{},
		},
		{
			"data/sysbench-parse-002.hex",
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
			_ = bytes.NewReader(testBytes)

			// pkt, err := protocol.NewParseWithReader(reader)
			// if err != nil {
			// 	t.Error(err)
			// }

			// // Compare the packet bytes

			// msgBytes, err := pkt.Bytes()
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }

			// if !bytes.Equal(msgBytes, testBytes) {
			// 	t.Errorf("expected %v, got %v", testBytes, msgBytes)
			// }
		})
	}
}
