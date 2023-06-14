// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

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

// ErrorResponse represents an error response message.
type ErrorResponse struct {
	*Response
}

// NewErrorResponse returns a new error response instance.
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Response: NewResponseWith(ErrorResponseMessage),
	}
}

// AppendTerminator appends a null terminator.
func (msg *ErrorResponse) AppendField(t ErrorType, v string) error {
	if err := msg.AppendByte(byte(t)); err != nil {
		return err
	}
	if _, err := msg.AppendString(v); err != nil {
		return err
	}
	return nil
}

// Bytes returns the message bytes.
func (msg *ErrorResponse) Bytes() ([]byte, error) {
	if err := msg.AppendTerminator(); err != nil {
		return nil, err
	}
	return msg.Response.Bytes()
}
