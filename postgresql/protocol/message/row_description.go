// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// RowDescription represents a row description message.
type RowDescription struct {
	*ResponseMessage
	fileds []*RowField
}

// NewRowDescription returns a new row description message instance.
func NewRowDescription() *RowDescription {
	return &RowDescription{
		ResponseMessage: NewResponseMessageWith(RowDescriptionMessage),
		fileds:          []*RowField{},
	}
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *RowDescription) Bytes() ([]byte, error) {
	msg.AppendInt16(int16(len(msg.fileds)))
	for _, field := range msg.fileds {
		err := field.WirteBytes(msg.Writer)
		if err != nil {
			return nil, err
		}
	}
	return msg.ResponseMessage.Bytes()
}
