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

// Responses represents a list of response.
type Responses []Response

// NewResponses returns a new empty responses.
func NewResponses() Responses {
	return []Response{}
}

// NewResponsesWith returns a new responses with the specified responses.
func NewResponsesWith(responses ...Response) Responses {
	return Responses(responses)
}

// NewCommandCompleteResponsesWith returns a new responses with the specified protocol.
func NewCommandCompleteResponsesWith(msg string) (Responses, error) {
	res, err := NewCommandCompleteWith(msg)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewInsertCompleteResponsesWith returns a new responses with the specified protocol.
func NewInsertCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewInsertCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewUpdateCompleteResponsesWith returns a new responses with the specified protocol.
func NewUpdateCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewUpdateCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewSelectCompleteResponsesWith returns a new responses with the specified protocol.
func NewSelectCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewSelectCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewDeleteCompleteResponsesWith returns a new responses with the specified protocol.
func NewDeleteCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewDeleteCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewCopyCompleteResponsesWith returns a new responses with the specified protocol.
func NewCopyCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewCopyCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewEmptyCompleteResponses returns a new responses with the specified protocol.
func NewEmptyCompleteResponses() (Responses, error) {
	res, err := NewEmptyComplete()
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// Append appends the specified response to this responses.
func (responses Responses) Append(res Response) Responses {
	return append(responses, res)
}

// HasErrorResponse returns true whether this responses has an error response.
func (responses Responses) HasErrorResponse() bool {
	for _, response := range responses {
		if response.Type() == ErrorResponseMessage {
			return true
		}
	}
	return false
}
