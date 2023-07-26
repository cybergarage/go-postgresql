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
)

// Float32ToBytes converts a float32 value to a byte array.
func Float32ToBytes(v float32) []byte {
	return Uint32ToBytes(math.Float32bits(v))
}

// Float64ToBytes converts a float64 value to a byte array.
func Float64ToBytes(v float64) []byte {
	return Uint64ToBytes(math.Float64bits(v))
}

// BytesToFloat32 converts a byte array to a float32 value.
func BytesToFloat32(b []byte) float32 {
	return math.Float32frombits(BytesToUint32(b))
}

// BytesToFloat64 converts a byte array to a float64 value.
func BytesToFloat64(b []byte) float64 {
	return math.Float64frombits(BytesToUint64(b))
}
