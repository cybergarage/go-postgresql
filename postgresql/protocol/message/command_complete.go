// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// CommandComplete represents a command complete message.
type CommandComplete struct {
	*ResponseMessage
}

// NewCommandComplete returns a new command complete message instance.
func NewCommandComplete() *CommandComplete {
	return &CommandComplete{
		ResponseMessage: NewResponseMessageWith(CommandCompleteMessage),
	}
}
