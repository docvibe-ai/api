package realestate

import (
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

// type Address struct {
// 	Street           nullable.TrimmedString  `json:"street,omitempty"`
// 	StreetVariations []notnull.TrimmedString `json:"street_variations,omitempty"`
// 	City             nullable.TrimmedString  `json:"city,omitempty"`
// 	State            nullable.TrimmedString  `json:"state,omitempty"`
// 	PostalCode       nullable.TrimmedString  `json:"postal_code,omitempty"`
// 	Country          country.NullableCode    `json:"country,omitempty"`
// }

type IdentifyObjectResult struct {
	ObjectID             notnull.TrimmedString   `json:"object_id"`
	AlternativeObjectIDs []notnull.TrimmedString `json:"alternative_object_ids,omitempty"`
}
