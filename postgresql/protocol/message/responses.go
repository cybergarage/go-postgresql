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

// Responses represents a list of response.
type Responses []Response

// NewResponses returns a new responses.
func NewResponsesWith(responses ...Response) Responses {
	return Responses(responses)
}

// NewCommandCompleteResponsesWith returns a new responses with the specified message.
func NewCommandCompleteResponsesWith(msg string) (Responses, error) {
	res, err := NewCommandCompleteWith(msg)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewInsertCompleteResponsesWith returns a new responses with the specified message.
func NewInsertCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewInsertCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewUpdateCompleteResponsesWith returns a new responses with the specified message.
func NewUpdateCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewUpdateCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewSelectCompleteResponsesWith returns a new responses with the specified message.
func NewSelectCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewSelectCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}

// NewDeleteCompleteResponsesWith returns a new responses with the specified message.
func NewDeleteCompleteResponsesWith(n int) (Responses, error) {
	res, err := NewDeleteCompleteWith(n)
	if err != nil {
		return nil, err
	}
	return Responses{res}, nil
}
