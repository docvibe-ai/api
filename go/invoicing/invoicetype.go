package invoicing

//go:generate go tool go-enum $GOFILE

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

type InvoiceType string //#enum,jsonschema

const (
	InvoiceTypeNull     InvoiceType = "" //#null
	InvoiceTypeIncoming InvoiceType = "INCOMING_INVOICE"
	InvoiceTypeOutgoing InvoiceType = "OUTGOING_INVOICE"
)

// Valid indicates if t is any of the valid values for InvoiceType
func (t InvoiceType) Valid() bool {
	switch t {
	case
		InvoiceTypeNull,
		InvoiceTypeIncoming,
		InvoiceTypeOutgoing:
		return true
	}
	return false
}

// Validate returns an error if t is none of the valid values for InvoiceType
func (t InvoiceType) Validate() error {
	if !t.Valid() {
		return fmt.Errorf("invalid value %#v for type invoicing.InvoiceType", t)
	}
	return nil
}

// Enums returns all valid values for InvoiceType
func (InvoiceType) Enums() []InvoiceType {
	return []InvoiceType{
		InvoiceTypeNull,
		InvoiceTypeIncoming,
		InvoiceTypeOutgoing,
	}
}

// EnumStrings returns all valid values for InvoiceType as strings
func (InvoiceType) EnumStrings() []string {
	return []string{
		"",
		"INCOMING_INVOICE",
		"OUTGOING_INVOICE",
	}
}

// String implements the fmt.Stringer interface for InvoiceType
func (t InvoiceType) String() string {
	return string(t)
}

// IsNull returns true if t is the null value InvoiceTypeNull
func (t InvoiceType) IsNull() bool {
	return t == InvoiceTypeNull
}

// IsNotNull returns true if t is not the null value InvoiceTypeNull
func (t InvoiceType) IsNotNull() bool {
	return t != InvoiceTypeNull
}

// SetNull sets the null value InvoiceTypeNull at t
func (t *InvoiceType) SetNull() {
	*t = InvoiceTypeNull
}

// MarshalJSON implements encoding/json.Marshaler for InvoiceType
// by returning the JSON null value for InvoiceTypeNull.
func (t InvoiceType) MarshalJSON() ([]byte, error) {
	if t == InvoiceTypeNull {
		return []byte("null"), nil
	}
	return json.Marshal(string(t))
}

// UnmarshalJSON implements encoding/json.Unmarshaler
func (t *InvoiceType) UnmarshalJSON(j []byte) error {
	if bytes.Equal(j, []byte("null")) {
		*t = InvoiceTypeNull
		return nil
	}
	return json.Unmarshal(j, (*string)(t))
}

// Scan implements the database/sql.Scanner interface for InvoiceType
func (t *InvoiceType) Scan(value any) error {
	switch value := value.(type) {
	case string:
		*t = InvoiceType(value)
	case []byte:
		*t = InvoiceType(value)
	case nil:
		*t = InvoiceTypeNull
	default:
		return fmt.Errorf("can't scan SQL value of type %T as invoicing.InvoiceType", value)
	}
	return nil
}

// Value implements the driver database/sql/driver.Valuer interface for InvoiceType
func (t InvoiceType) Value() (driver.Value, error) {
	if t == InvoiceTypeNull {
		return nil, nil
	}
	return string(t), nil
}

// JSONSchema returns a github.com/invopop/jsonschema.Schema for InvoiceType
func (InvoiceType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type: "string",
				Enum: []any{
					"INCOMING_INVOICE",
					"OUTGOING_INVOICE",
				},
			},
			{Type: "null"},
		},
		Default: InvoiceTypeNull,
	}
}
