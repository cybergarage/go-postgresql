// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

package message

// RowDescription represents a row description field.
type RowField struct {
	Name         string
	TableID      int32
	Number       int16
	DataTypeID   int32
	DataTypeSize int16
	TypeModifier int32
	FormatCode   int16
}

// NewRowField returns a new row description field.
func NewRowFieldWith(name string, n int) *RowField {
	return &RowField{
		Name:         name,
		TableID:      0,
		Number:       int16(n),
		DataTypeID:   0,
		TypeModifier: 0,
		FormatCode:   0,
	}
}

// WirteBytes appends a row field elements.
func (field *RowField) WirteBytes(w *Writer) error {
	if err := w.AppendString(field.Name); err != nil {
		return err
	}
	if err := w.AppendInt32(field.TableID); err != nil {
		return err
	}
	if err := w.AppendInt16(field.Number); err != nil {
		return err
	}
	if err := w.AppendInt32(field.DataTypeID); err != nil {
		return err
	}
	if err := w.AppendInt16(field.DataTypeSize); err != nil {
		return err
	}
	if err := w.AppendInt32(field.TypeModifier); err != nil {
		return err
	}
	if err := w.AppendInt16(field.FormatCode); err != nil {
		return err
	}
	return nil
}
