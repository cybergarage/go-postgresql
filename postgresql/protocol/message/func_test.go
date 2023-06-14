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

package message

import (
	"math"
	"testing"
)

func TestInt32Convert(t *testing.T) {
	ts := []int32{
		math.MinInt32,
		math.MinInt32 / 2,
		-1,
		0,
		1,
		math.MaxInt32 / 2,
		math.MaxInt32,
	}

	for _, tv := range ts {
		b := Int32ToBytes(tv)
		v := Int32BytesToInt(b)
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}