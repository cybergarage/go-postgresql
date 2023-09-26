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
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	// Create a buffer with some data
	buf := bytes.NewBuffer([]byte{0x01, 0x02, 0x03, 0x04})

	// Create a new reader with the buffer
	reader := NewReaderWith(buf)

	// Test PeekInt32
	expectedInt32 := int32(0x01020304)
	actualInt32, err := reader.PeekInt32()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt32 != expectedInt32 {
		t.Errorf("Expected %v, but got %v", expectedInt32, actualInt32)
	}

	// Test ReadInt32
	actualInt32, err = reader.ReadInt32()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt32 != expectedInt32 {
		t.Errorf("Expected %v, but got %v", expectedInt32, actualInt32)
	}

	// Test ReadInt16
	expectedInt16 := int16(0x0102)
	actualInt16, err := reader.ReadInt16()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt16 != expectedInt16 {
		t.Errorf("Expected %v, but got %v", expectedInt16, actualInt16)
	}

	// Test ReadBytesUntil
	expectedBytes := []byte{0x01, 0x02, 0x03}
	actualBytes, err := reader.ReadBytesUntil(0x04)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Errorf("Expected %v, but got %v", expectedBytes, actualBytes)
	}

	// Test ReadString
	expectedString := "\x01\x02\x03"
	actualString, err := reader.ReadString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}
}
