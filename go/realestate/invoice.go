package realestate

import "github.com/docvibe-ai/api/go/invoicing"

type Invoice struct {
	invoicing.AccountingInvoice

	Section35aAmounts []*Section35aInvoiceAmount `json:"section35a_amounts,omitempty"`
	IdentifiedObjects []*Object                  `json:"identified_objects,omitempty"`
}
