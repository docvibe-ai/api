package invoice

//go:generate go tool go-enum $GOFILE

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Type string //#enum

const (
	TypeNull            Type = "" //#null
	TypeIncomingInvoice Type = "INCOMING_INVOICE"
	TypeOutgoingInvoice Type = "OUTGOING_INVOICE"
)

// Valid indicates if t is any of the valid values for Type
func (t Type) Valid() bool {
	switch t {
	case
		TypeNull,
		TypeIncomingInvoice,
		TypeOutgoingInvoice:
		return true
	}
	return false
}

// Validate returns an error if t is none of the valid values for Type
func (t Type) Validate() error {
	if !t.Valid() {
		return fmt.Errorf("invalid value %#v for type invoice.Type", t)
	}
	return nil
}

// Enums returns all valid values for Type
func (Type) Enums() []Type {
	return []Type{
		TypeNull,
		TypeIncomingInvoice,
		TypeOutgoingInvoice,
	}
}

// EnumStrings returns all valid values for Type as strings
func (Type) EnumStrings() []string {
	return []string{
		"",
		"INCOMING_INVOICE",
		"OUTGOING_INVOICE",
	}
}

// String implements the fmt.Stringer interface for Type
func (t Type) String() string {
	return string(t)
}

// IsNull returns true if t is the null value TypeNull
func (t Type) IsNull() bool {
	return t == TypeNull
}

// IsNotNull returns true if t is not the null value TypeNull
func (t Type) IsNotNull() bool {
	return t != TypeNull
}

// SetNull sets the null value TypeNull at t
func (t *Type) SetNull() {
	*t = TypeNull
}

// MarshalJSON implements encoding/json.Marshaler for Type
// by returning the JSON null value for TypeNull.
func (t Type) MarshalJSON() ([]byte, error) {
	if t == TypeNull {
		return []byte("null"), nil
	}
	return json.Marshal(string(t))
}

// UnmarshalJSON implements encoding/json.Unmarshaler
func (t *Type) UnmarshalJSON(j []byte) error {
	if bytes.Equal(j, []byte("null")) {
		*t = TypeNull
		return nil
	}
	return json.Unmarshal(j, (*string)(t))
}

// Scan implements the database/sql.Scanner interface for Type
func (t *Type) Scan(value any) error {
	switch value := value.(type) {
	case string:
		*t = Type(value)
	case []byte:
		*t = Type(value)
	case nil:
		*t = TypeNull
	default:
		return fmt.Errorf("can't scan SQL value of type %T as invoice.Type", value)
	}
	return nil
}

// Value implements the driver database/sql/driver.Valuer interface for Type
func (t Type) Value() (driver.Value, error) {
	if t == TypeNull {
		return nil, nil
	}
	return string(t), nil
}
