// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// CloseComplete represents a close complete message.
type CloseComplete struct {
	*ResponseMessage
}

// NewCloseComplete returns a new close complete message instance.
func NewCloseComplete() *CloseComplete {
	return &CloseComplete{
		ResponseMessage: NewResponseMessageWith(CloseCompleteMessage),
	}
}
