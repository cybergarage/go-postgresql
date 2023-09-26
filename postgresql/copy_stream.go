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
	"errors"
	"io"

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
func (stream *CopyStream) Next() (*message.CopyData, error) {
	t, err := stream.MessageReader.PeekType()
	if err != nil {
		return nil, err
	}

	skipCopyDone := func(reader *message.MessageReader) error {
		_, err := message.NewCopyDoneWithReader(reader)
		return err
	}

	switch t { // nolint:exhaustive
	case message.CopyDataMessage:
		copyData, copyErr := message.NewCopyDataWithReader(stream.MessageReader)
		if copyErr == nil {
			return copyData, nil
		}
		if !errors.Is(copyErr, io.EOF) {
			return nil, copyErr
		}
		ok, peekErr := stream.MessageReader.IsPeekType(message.CopyDoneMessage)
		if peekErr != nil {
			return nil, peekErr
		}
		if !ok {
			return nil, copyErr
		}
		if skipErr := skipCopyDone(stream.MessageReader); skipErr != nil {
			return nil, skipErr
		}
		return nil, copyErr
	case message.CopyDoneMessage:
		err := skipCopyDone(stream.MessageReader)
		if err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	return nil, io.EOF
}
