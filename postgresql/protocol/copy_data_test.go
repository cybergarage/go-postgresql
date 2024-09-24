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

package protocol

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"testing"
)

func TestCopyData(t *testing.T) {
	tests := []struct {
		data     string
		expected []string
	}{
		{
			data:     "640000000b3109310930090a",
			expected: []string{"1", "1", "0", ""},
		},
		{
			data:     "640000000b3209310930090a",
			expected: []string{"2", "1", "0", ""},
		},
		{
			data:     "640000000f313730373909310930090a",
			expected: []string{"17079", "1", "0", ""},
		},
	}

	for _, test := range tests {
		byteData, err := hex.DecodeString(test.data)
		if err != nil {
			t.Error(err)
			return
		}

		reader := NewMessageReaderWith(bufio.NewReader(bytes.NewReader(byteData)))
		copyData, err := NewCopyDataWithReader(reader)
		if err != nil {
			t.Error(err)
			return
		}

		if len(test.expected) != len(copyData.Data) {
			t.Errorf("expected %d, got %d", len(test.expected), len(copyData.Data))
			return
		}

		for n, expected := range test.expected {
			if expected != copyData.Data[n] {
				t.Errorf("expected %s, got %s", expected, copyData.Data[n])
				return
			}
		}
	}
}
