// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// ParseComplete represents a parser complete message.
type ParseComplete struct {
	*ResponseMessage
}

// NewParseComplete returns a parser complete instance.
func NewParseComplete() *ParseComplete {
	return &ParseComplete{
		ResponseMessage: NewResponseMessageWith(ParseCompleteMessage),
	}
}
