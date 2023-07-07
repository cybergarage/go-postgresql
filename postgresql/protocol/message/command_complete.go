// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

import (
	"fmt"
)

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

// NewInsertCompleteWith returns a new command complete message for insert query.
func NewInsertCompleteWith(n int) (*CommandComplete, error) {
	msg := &CommandComplete{
		ResponseMessage: NewResponseMessageWith(CommandCompleteMessage),
	}
	return msg, msg.AppendString(fmt.Sprintf("INSERT 0 %d", n))
}

// NewUpdateCompleteWith returns a new command complete message for update query.
func NewUpdateCompleteWith(n int) (*CommandComplete, error) {
	msg := &CommandComplete{
		ResponseMessage: NewResponseMessageWith(CommandCompleteMessage),
	}
	return msg, msg.AppendString(fmt.Sprintf("UPDATE %d", n))
}

// NewDeleteCompleteWith returns a new command complete message for delete query.
func NewDeleteCompleteWith(n int) (*CommandComplete, error) {
	msg := &CommandComplete{
		ResponseMessage: NewResponseMessageWith(CommandCompleteMessage),
	}
	return msg, msg.AppendString(fmt.Sprintf("DELETE %d", n))
}
