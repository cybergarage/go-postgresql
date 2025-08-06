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
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

// ErrorType represents a error response type.
type ErrorType byte

const (
	SeverityError         ErrorType = 'S'
	CodeError             ErrorType = 'C'
	MessageError          ErrorType = 'M'
	DetailError           ErrorType = 'D'
	HintError             ErrorType = 'H'
	PositionError         ErrorType = 'P'
	InternalPositionError ErrorType = 'p'
	InternalQueryError    ErrorType = 'q'
	WhereError            ErrorType = 'W'
	SchemaError           ErrorType = 's'
	TableError            ErrorType = 't'
	ColumnError           ErrorType = 'c'
	DataTypeNameError     ErrorType = 'd'
	ConstraintError       ErrorType = 'n'
	FileError             ErrorType = 'F'
	LineError             ErrorType = 'L'
	RoutineError          ErrorType = 'R'
)

// ErrorResponse represents an error response protocol.
type ErrorResponse struct {
	*ResponseMessage
}

// NewErrorResponse returns a new error response instance.
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		ResponseMessage: NewResponseMessageWith(ErrorResponseMessage),
	}
}

// NewErrorResponseWith returns a new error response instance with the specified error.
func NewErrorResponseWith(err error) (*ErrorResponse, error) {
	msg := NewErrorResponse()
	return msg, msg.AddError(err)
}

// AppendField appends an error field to the error response.
func (msg *ErrorResponse) AppendField(t ErrorType, v string) error {
	if err := msg.AppendByte(byte(t)); err != nil {
		return err
	}

	return msg.AppendString(v)
}

// AddCode adds an error code to the error response.
func (msg *ErrorResponse) AddCode(code int32) error {
	if err := msg.AppendByte(byte(CodeError)); err != nil {
		return err
	}

	return msg.AppendInt32(code)
}

// AddError adds an error message to the error response.
func (msg *ErrorResponse) AddError(err error) error {
	if err := msg.AppendByte(byte(MessageError)); err != nil {
		return err
	}

	return msg.AppendString(err.Error())
}

// Bytes returns the message bytes after adding a null terminator.
func (msg *ErrorResponse) Bytes() ([]byte, error) {
	if err := msg.AppendTerminator(); err != nil {
		return nil, err
	}

	return msg.ResponseMessage.Bytes()
}
