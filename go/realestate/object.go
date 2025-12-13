package realestate

import (
	"fmt"
	"slices"

	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
)

type Object struct {
	ID    notnull.TrimmedString  `json:"id"`
	Type  nullable.TrimmedString `json:"type,omitempty"`
	Notes nullable.TrimmedString `json:"notes,omitempty"`
	// Status nullable.TrimmedString `json:"status,omitempty"`

	Street           nullable.TrimmedString  `json:"street,omitempty"`
	StreetVariations []notnull.TrimmedString `json:"street_variations,omitempty"`
	City             nullable.TrimmedString  `json:"city,omitempty"`
	State            nullable.TrimmedString  `json:"state,omitempty"`
	PostalCode       nullable.TrimmedString  `json:"postal_code,omitempty"`
	Country          country.NullableCode    `json:"country,omitempty"`
}

// Normalize validates and normalizes all fields of the Object.
// It returns a slice of all validation errors found.
// Invalid fields are set to null values. The object remains usable
// after normalization, with the returned errors describing what was corrected.
func (o *Object) Normalize() (errs []error) {
	if o == nil {
		return nil
	}
	var err error

	// Normalize country code
	if o.Country, err = o.Country.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid country code: %w", err))
		o.Country.SetNull()
	}

	// Remove empty street variations
	o.StreetVariations = slices.DeleteFunc(o.StreetVariations, func(v notnull.TrimmedString) bool {
		return v.IsEmpty()
	})

	return errs
}

// type Address struct {
// 	Street           nullable.TrimmedString  `json:"street,omitempty"`
// 	StreetVariations []notnull.TrimmedString `json:"street_variations,omitempty"`
// 	City             nullable.TrimmedString  `json:"city,omitempty"`
// 	State            nullable.TrimmedString  `json:"state,omitempty"`
// 	PostalCode       nullable.TrimmedString  `json:"postal_code,omitempty"`
// 	Country          country.NullableCode    `json:"country,omitempty"`
// }
