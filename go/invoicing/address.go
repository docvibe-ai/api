package invoicing

import (
	"fmt"

	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/nullable"
)

type Address struct {
	Street     nullable.TrimmedString `json:"street,omitempty"`
	City       nullable.TrimmedString `json:"city,omitempty"`
	State      nullable.TrimmedString `json:"state,omitempty"`
	PostalCode nullable.TrimmedString `json:"postal_code,omitempty"`
	Country    country.NullableCode   `json:"country,omitempty"`
}

// Normalize validates and normalizes all fields of the Address.
// It returns a slice of all validation errors found.
// Invalid fields are set to null values. The address remains usable
// after normalization, with the returned errors describing what was corrected.
func (a *Address) Normalize() []error {
	if a == nil {
		return nil
	}
	var err error
	var errs []error
	if a.Country, err = a.Country.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid address country code: %w", err))
		a.Country.SetNull()
	}
	return errs
}
