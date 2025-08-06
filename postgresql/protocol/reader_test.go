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
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	// Create a buffer with some data
	buf := []byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x00}

	// Test PeekInt32 and ReadInt32
	reader := NewReaderWith(WithReaderBytes(buf))
	expectedInt32 := int32(0x61626364)
	actualInt32, err := reader.PeekInt32()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt32 != expectedInt32 {
		t.Errorf("Expected %v, but got %v", expectedInt32, actualInt32)
	}
	actualInt32, err = reader.ReadInt32()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt32 != expectedInt32 {
		t.Errorf("Expected %v, but got %v", expectedInt32, actualInt32)
	}
	expectedInt32 = int32(0x65666768)
	actualInt32, err = reader.ReadInt32()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt32 != expectedInt32 {
		t.Errorf("Expected %v, but got %v", expectedInt32, actualInt32)
	}

	// Test ReadInt16
	reader = NewReaderWith(WithReaderBytes(buf))
	expectedInt16 := int16(0x6162)
	actualInt16, err := reader.ReadInt16()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt16 != expectedInt16 {
		t.Errorf("Expected %v, but got %v", expectedInt16, actualInt16)
	}

	// Test ReadBytesUntil
	reader = NewReaderWith(WithReaderBytes(buf))
	expectedBytes := []byte{0x61, 0x62, 0x63, 0x64}
	actualBytes, err := reader.ReadBytesUntil(0x64)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Errorf("Expected %v, but got %v", expectedBytes, actualBytes)
	}

	// Test ReadString
	reader = NewReaderWith(WithReaderBytes(buf))
	expectedString := "\x61\x62\x63\x64\x65\x66\x67\x68"
	actualString, err := reader.ReadString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}
}
