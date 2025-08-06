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

package bytes

import (
	"math"
	"testing"
)

func TestFloat32Encode(t *testing.T) {
	ts := []float32{
		math.SmallestNonzeroFloat32,
		math.MaxFloat32 / 2,
		math.MaxFloat32,
	}

	for _, tv := range ts {
		b := Float32ToBytes(tv)

		v := BytesToFloat32(b)
		if tv != v {
			t.Errorf("Failed to convert (%f != %f)", tv, v)
		}
	}
}

func TestFloat64Encode(t *testing.T) {
	ts := []float64{
		math.SmallestNonzeroFloat64,
		math.MaxFloat64 / 2,
		math.MaxFloat64,
	}

	for _, tv := range ts {
		b := Float64ToBytes(tv)

		v := BytesToFloat64(b)
		if tv != v {
			t.Errorf("Failed to convert (%f != %f)", tv, v)
		}
	}
}
