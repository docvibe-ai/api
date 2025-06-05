package invoice

import (
	"errors"
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

func (a *Address) Normalize() (err error) {
	if a == nil {
		return nil
	}
	if a.Country, err = a.Country.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid country code: %w", err))
		a.Country.SetNull()
	}
	return err
}
