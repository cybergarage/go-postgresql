// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

package message

// BackendKeyData represents a parameter status response message.
type BackendKeyData struct {
	*ResponseMessage
}

// NewBackendKeyData returns a parameter status response instance.
func NewBackendKeyData() *BackendKeyData {
	return &BackendKeyData{
		ResponseMessage: NewResponseMessageWith(BackendKeyDataMessage),
	}
}

// NewBackendKeyDataWith returns a parameter status response instance with the specified parameters.
func NewBackendKeyDataWith(processID int32, secretKey int32) (*BackendKeyData, error) {
	msg := &BackendKeyData{
		ResponseMessage: NewResponseMessageWith(BackendKeyDataMessage),
	}
	err := msg.AppendInt32(12)
	if err != nil {
		return nil, err
	}
	err = msg.AppendInt32(processID)
	if err != nil {
		return nil, err
	}
	err = msg.AppendInt32(secretKey)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
