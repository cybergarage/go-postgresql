// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package protocol

// BindComplete represents a bind complete protocol.
type BindComplete struct {
	*ResponseMessage
}

// NewBindComplete returns a new bind complete message instance.
func NewBindComplete() *BindComplete {
	return &BindComplete{
		ResponseMessage: NewResponseMessageWith(BindCompleteMessage),
	}
}
