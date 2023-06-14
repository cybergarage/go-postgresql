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

package message

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// Startup represents a startup message.
type Startup struct {
	MajorVersion  int
	MinorVersion  int
	MessageLength int32
	Parameters    map[string]string
}

// NewStartup returns a new startup message.
func NewStartupWith(reader *Reader) (*Startup, error) {
	readLen := 0
	msgLen, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	readLen += 4

	ver, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	readLen += 4

	majorVer := ver >> 16
	minorVer := ver & 0xFFFF

	params := make(map[string]string)

	for readLen < int(msgLen) {
		k, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		kl := len(k)
		readLen += kl + 1
		if kl == 0 {
			break
		}
		v, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		readLen += len(v) + 1
		params[k] = v
	}

	return &Startup{
		MessageLength: msgLen,
		MajorVersion:  int(majorVer),
		MinorVersion:  int(minorVer),
		Parameters:    params,
	}, nil
}
