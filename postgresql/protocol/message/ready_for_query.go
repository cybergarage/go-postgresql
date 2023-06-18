// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// ReadyForQuery represents a ready for query message.
type ReadyForQuery struct {
	*ResponseMessage
}

// TransactionStatus represents a transaction status.
type TransactionStatus = byte

const (
	TransactionIdle   = 'I'
	TransactionBlock  = 'T'
	TransactionFailed = 'E'
)

// NewReadyForQuery returns a new ready for query message instance.
func NewReadyForQuery() *ReadyForQuery {
	return &ReadyForQuery{
		ResponseMessage: NewResponseMessageWith(ReadyForQueryMessage),
	}
}

// NewReadyForQueryWith returns a new error response instance with the specified error.
func NewReadyForQueryWith(s TransactionStatus) (*ReadyForQuery, error) {
	msg := NewReadyForQuery()
	return msg, msg.AppendByte(s)
}
