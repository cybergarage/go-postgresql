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
}

// NewDataRow returns a new data row message instance.
func NewDataRow() *DataRow {
	return &DataRow{
		ResponseMessage: NewResponseMessageWith(DataRowMessage),
	}
}
