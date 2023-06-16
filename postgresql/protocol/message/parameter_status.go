// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// ParameterStatus represents an error response message.
type ParameterStatus struct {
	*StringResponse
}

// NewParameterStatus returns a new error response instance.
func NewParameterStatus() *ParameterStatus {
	return &ParameterStatus{
		StringResponse: NewStringResponseWith(ParameterStatusMessage),
	}
}
