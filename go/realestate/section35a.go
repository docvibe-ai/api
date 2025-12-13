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

// Normalize validates and normalizes all fields of the Section35aInvoiceAmount.
// It returns a slice of all validation errors found.
// Money amounts are normalized to their absolute (positive) values.
// Zero amounts are set to null. If both NetAmount and GrossAmount are not null,
// NetAmount must not be greater than GrossAmount (otherwise NetAmount is set to null).
func (s *Section35aInvoiceAmount) Normalize() (errs []error) {
	if s == nil {
		return nil
	}
	var err error

	// Validate type
	if err = s.Type.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("invalid section 35a type: %w, using %s", err, Section35aTypeCraftsmanServices))
		s.Type = Section35aTypeCraftsmanServices
	}

	// Normalize NetAmount: absolute value, zero becomes null
	if s.NetAmount.IsNotNull() {
		s.NetAmount.Set(s.NetAmount.Get().Abs())
		if s.NetAmount.Get() == 0 {
			s.NetAmount.SetNull()
		}
	}

	// Normalize GrossAmount: absolute value, zero becomes null
	if s.GrossAmount.IsNotNull() {
		s.GrossAmount.Set(s.GrossAmount.Get().Abs())
		if s.GrossAmount.Get() == 0 {
			s.GrossAmount.SetNull()
		}
	}

	// Validate NetAmount <= GrossAmount
	if s.NetAmount.IsNotNull() && s.GrossAmount.IsNotNull() {
		if s.NetAmount.Get() > s.GrossAmount.Get() {
			errs = append(errs, fmt.Errorf("net amount %s is greater than gross amount %s, setting net amount to null", s.NetAmount.Get(), s.GrossAmount.Get()))
			s.NetAmount.SetNull()
		}
	}

	return errs
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
