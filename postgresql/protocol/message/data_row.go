// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// DataRow represents a data row message.
type DataRow struct {
	*ResponseMessage
	ColumnValues []any
}

// NewDataRow returns a new data row message instance.
func NewDataRow() *DataRow {
	return &DataRow{
		ResponseMessage: NewResponseMessageWith(DataRowMessage),
	}
}

// AppendColumnValue appends a column value to the data row message.
func (msg *DataRow) AppendColumnValue(v any) error {
	msg.ColumnValues = append(msg.ColumnValues, v)
	return nil
}

// Bytes appends a length of the message content bytes, and returns the message bytes.
func (msg *DataRow) Bytes() ([]byte, error) {
	err := msg.AppendInt16(int16(len(msg.ColumnValues)))
	if err != nil {
		return nil, err
	}
	for _, v := range msg.ColumnValues {
		switch v := v.(type) {
		case []byte:
			if err := msg.AppendInt32(int32(len(v))); err != nil {
				return nil, err
			}
			if err := msg.AppendBytes(v); err != nil {
				return nil, err
			}
		case string:
			if err := msg.AppendInt32(int32(len(v))); err != nil {
				return nil, err
			}
			if err := msg.AppendString(v); err != nil {
				return nil, err
			}
		case int32:
			if err := msg.AppendInt32(4); err != nil {
				return nil, err
			}
			if err := msg.AppendInt32(v); err != nil {
				return nil, err
			}
		default:
			return nil, newColumnTypeNotSuppotedError(v)
		}
	}
	return msg.ResponseMessage.Bytes()
}
