package realestate

//go:generate go tool go-enum $GOFILE

import (
	"fmt"

	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/nullable"
	"github.com/invopop/jsonschema"
)

type Section35aInvoiceAmount struct {
	Type        Section35aType         `json:"type"`
	NetAmount   money.NullableAmount   `json:"net_amount"`
	GrossAmount money.NullableAmount   `json:"gross_amount"`
	Purpose     nullable.TrimmedString `json:"purpose"`
}

type Section35aType string //#enum,jsonschema

const (
	Section35aTypeFormallyEmployedWorker Section35aType = "FORMALLY_EMPLOYED_WORKER" // Personalkosten für sozialversicherungspflichtige Beschäftigungsverhältnisse im Privathaushalt
	Section35aTypeHouseholdServices      Section35aType = "HOUSEHOLD_SERVICES"       // Haushaltsnahe Dienstleistungen, Hilfe im Haushalt
	Section35aTypeCraftsmanServices      Section35aType = "CRAFTSMAN_SERVICES"       // Handwerkerleistungen
)

// Valid indicates if s is any of the valid values for Section35aType
func (s Section35aType) Valid() bool {
	switch s {
	case
		Section35aTypeFormallyEmployedWorker,
		Section35aTypeHouseholdServices,
		Section35aTypeCraftsmanServices:
		return true
	}
	return false
}

// Validate returns an error if s is none of the valid values for Section35aType
func (s Section35aType) Validate() error {
	if !s.Valid() {
		return fmt.Errorf("invalid value %#v for type realestate.Section35aType", s)
	}
	return nil
}

// Enums returns all valid values for Section35aType
func (Section35aType) Enums() []Section35aType {
	return []Section35aType{
		Section35aTypeFormallyEmployedWorker,
		Section35aTypeHouseholdServices,
		Section35aTypeCraftsmanServices,
	}
}

// EnumStrings returns all valid values for Section35aType as strings
func (Section35aType) EnumStrings() []string {
	return []string{
		"FORMALLY_EMPLOYED_WORKER",
		"HOUSEHOLD_SERVICES",
		"CRAFTSMAN_SERVICES",
	}
}

// String implements the fmt.Stringer interface for Section35aType
func (s Section35aType) String() string {
	return string(s)
}

// JSONSchema returns a github.com/invopop/jsonschema.Schema for Section35aType
func (Section35aType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			"FORMALLY_EMPLOYED_WORKER",
			"HOUSEHOLD_SERVICES",
			"CRAFTSMAN_SERVICES",
		},
	}
}
