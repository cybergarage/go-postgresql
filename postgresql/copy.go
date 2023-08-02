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

package postgresql

import (
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
)

// CopyStream represents a copy stream.
type CopyStream struct {
	*message.MessageReader
}

// NewCopyStreamWithReader returns a new copy stream with the specified reader.
func NewCopyStreamWithReader(reader *message.MessageReader) *CopyStream {
	return &CopyStream{
		MessageReader: reader,
	}
}

// Next returns true if the next message is available.
func (stream *CopyStream) Next() (bool, error) {
	reqType, err := stream.PeekType()
	if err != nil {
		return false, err
	}

	switch reqType { // nolint:exhaustive
	case message.CopyDataMessage:
	case message.CopyDoneMessage:
		_, err := message.NewCopyDoneWithReader(stream.MessageReader)
		if err != nil {
			return false, err
		}
	default:
		return false, message.NewErrInvalidMessage(reqType)
	}

	return true, nil
}

// CopyData returns a copy data message.
func (stream *CopyStream) CopyData() (*message.CopyData, error) {
	return message.NewCopyDataWithReader(stream.MessageReader)
}
