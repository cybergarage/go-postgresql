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

// BytesToUint64 converts the specified byte array to an integer.
func BytesToUint64(b []byte) uint64 {
	v := uint64(b[0])<<56 |
		uint64(b[1])<<48 |
		uint64(b[2])<<40 |
		uint64(b[3])<<32 |
		uint64(b[4])<<24 |
		uint64(b[5])<<16 |
		uint64(b[6])<<8 |
		uint64(b[7])
	return v
}

// Uint64ToBytes converts the specified integer to a byte array.
func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	return b
}

// BytesToUint32 converts the specified byte array to an integer.
func BytesToUint32(b []byte) uint32 {
	v := uint32(b[0])<<24 |
		uint32(b[1])<<16 |
		uint32(b[2])<<8 |
		uint32(b[3])
	return v
}

// Uint32ToBytes converts the specified integer to a byte array.
func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
	return b
}

// BytesToUint16 converts the specified byte array to an integer.
func BytesToUint16(b []byte) uint16 {
	v := uint16(b[0])<<8 | uint16(b[1])
	return v
}

// Uint16ToBytes converts the specified integer to a byte array.
func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return b
}

// BytesToUint8 converts the specified byte array to an integer.
func BytesToUint8(b []byte) uint8 {
	return b[0]
}

// Uint8ToBytes converts the specified integer to a byte array.
func Uint8ToBytes(v uint8) []byte {
	b := make([]byte, 1)
	b[0] = v
	return b
}
