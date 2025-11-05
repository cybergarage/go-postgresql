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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// 54.2.7. Asynchronous Operations
// https://www.postgresql.org/docs/current/protocol-flow.html#PROTOCOL-ASYNC

const (
	ApplicationName            = "application_name"
	ScramIterations            = "scram_iterations"
	ClientEncoding             = "client_encoding"
	SearchPath                 = "search_path"
	DateStyle                  = "DateStyle"
	ServerEncoding             = "server_encoding"
	DefaultTransactionReadOnly = "default_transaction_read_only"
	ServerVersion              = "server_version"
	InHotStandby               = "in_hot_standby"
	SessionAuth                = "session_authorization"
	IntegerDatetimes           = "integer_datetimes"
	StandardConformingStrings  = "standard_conforming_strings"
	IsSuperuser                = "is_superuser"
	IntervalStyle              = "IntervalStyle"
	TimeZone                   = "TimeZone"
)

const (
	EncodingUTF8 = "UTF8"
)

const (
	DateStyleISO = "ISO, MDY.S"
)

// ParameterStatus represents a parameter status response protocol.
type ParameterStatus struct {
	*ResponseMessage
}

// NewParameterStatus returns a parameter status response instance.
func NewParameterStatus() *ParameterStatus {
	return &ParameterStatus{
		ResponseMessage: NewResponseMessageWith(ParameterStatusMessage),
	}
}

// NewParameterStatusWith returns a parameter status response instance with the specified parameter status.
func NewParameterStatusWith(name string, value string) (*ParameterStatus, error) {
	msg := NewParameterStatus()
	err := msg.AppendParameters(name, value)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// NewParameterStatusesWith returns parameter status response instances with the specified parameter statuses.
func NewParameterStatusesWith(m map[string]string) (Responses, error) {
	msgs := Responses{}
	for k, v := range m {
		msg, err := NewParameterStatusWith(k, v)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// AppendParameters appends the specified parameters.
func (msg *ParameterStatus) AppendParameters(s ...string) error {
	for _, v := range s {
		err := msg.AppendString(v)
		if err != nil {
			return err
		}
	}
	return nil
}
