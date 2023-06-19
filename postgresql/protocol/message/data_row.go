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
