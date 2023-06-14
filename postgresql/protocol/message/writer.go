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
	"bufio"
	"bytes"
)

// Writer represents a message writer.
type Writer struct {
	*bytes.Buffer
	*bufio.Writer
}

// NewWriter returns a new message writer.
func NewWriter() *Writer {
	buffer := &bytes.Buffer{}
	return &Writer{
		Buffer: buffer,
		Writer: bufio.NewWriter(buffer),
	}
}

// Bytes returns the message bytes.
func (writer *Writer) Bytes() ([]byte, error) {
	err := writer.Flush()
	if err != nil {
		return nil, err
	}
	return writer.Buffer.Bytes(), nil
}
