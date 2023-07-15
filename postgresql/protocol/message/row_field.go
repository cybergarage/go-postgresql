// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// RowDescription represents a row description field.
type RowField struct {
	Name                 string
	TableOID             int32
	TableAttributeNumber int16
	DataTypeOID          int32
	DataTypeSize         int16
	TypeModifier         int32
	FormatCode           int16
}

// NewRowField returns a new row description field.
func NewRowFieldWith(name string) *RowField {
	return &RowField{
		Name: name,
	}
}
