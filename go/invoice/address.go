package invoice

import (
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
