package realestate

import (
	"fmt"
	"slices"

	"github.com/docvibe-ai/api/go/invoicing"
)

type Invoice struct {
	invoicing.AccountingInvoice

	Section35aAmounts []*Section35aInvoiceAmount `json:"section35a_amounts,omitempty"`
	IdentifiedObjects []*Object                  `json:"identified_objects,omitempty"`
}

// Normalize validates and normalizes all fields of the Invoice.
// It returns a slice of all validation errors found.
// Invalid fields are set to null values or corrected. The invoice remains usable
// after normalization, with the returned errors describing what was corrected.
func (inv *Invoice) Normalize() []error {
	if inv == nil {
		return nil
	}

	// Normalize the embedded AccountingInvoice
	errs := inv.AccountingInvoice.Normalize()

	// Remove nil Section35a amounts and normalize the rest
	inv.Section35aAmounts = slices.DeleteFunc(inv.Section35aAmounts, func(a *Section35aInvoiceAmount) bool {
		return a == nil
	})
	for i, amount := range inv.Section35aAmounts {
		if amountErrs := amount.Normalize(); len(amountErrs) > 0 {
			for _, e := range amountErrs {
				errs = append(errs, fmt.Errorf("section35a amount %d: %w", i+1, e))
			}
		}
	}

	// Remove nil identified objects and normalize the rest
	inv.IdentifiedObjects = slices.DeleteFunc(inv.IdentifiedObjects, func(o *Object) bool {
		return o == nil
	})
	for i, obj := range inv.IdentifiedObjects {
		if objErrs := obj.Normalize(); len(objErrs) > 0 {
			for _, e := range objErrs {
				errs = append(errs, fmt.Errorf("identified object %d: %w", i+1, e))
			}
		}
	}

	return errs
}
