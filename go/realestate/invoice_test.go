package realestate

import (
	"testing"

	"github.com/domonda/go-types/money"
)

func TestInvoice_Normalize(t *testing.T) {
	tests := []struct {
		name       string
		invoice    *Invoice
		wantErrors int
	}{
		{
			name:       "nil returns no errors",
			invoice:    nil,
			wantErrors: 0,
		},
		{
			name:       "empty struct returns error from embedded invoice",
			invoice:    &Invoice{},
			wantErrors: 1, // invalid payment status from embedded Invoice.Normalize()
		},
		{
			name: "valid Section35a amounts are normalized",
			invoice: &Invoice{
				Section35aAmounts: []*Section35aInvoiceAmount{
					{
						Type:        Section35aTypeCraftsmanServices,
						NetAmount:   money.NullableAmountFrom(100),
						GrossAmount: money.NullableAmountFrom(119),
					},
				},
			},
			wantErrors: 1, // +1 from embedded invoice
		},
		{
			name: "invalid Section35a type returns error",
			invoice: &Invoice{
				Section35aAmounts: []*Section35aInvoiceAmount{
					{
						Type:        "INVALID_TYPE",
						NetAmount:   money.NullableAmountFrom(100),
						GrossAmount: money.NullableAmountFrom(119),
					},
				},
			},
			wantErrors: 2, // 1 from section35a + 1 from embedded invoice
		},
		{
			name: "multiple Section35a amounts with errors",
			invoice: &Invoice{
				Section35aAmounts: []*Section35aInvoiceAmount{
					{
						Type:        "INVALID_TYPE",
						NetAmount:   money.NullableAmountFrom(100),
						GrossAmount: money.NullableAmountFrom(119),
					},
					{
						Type:        Section35aTypeCraftsmanServices,
						NetAmount:   money.NullableAmountFrom(150),
						GrossAmount: money.NullableAmountFrom(100),
					},
				},
			},
			wantErrors: 3, // 2 from section35a + 1 from embedded invoice
		},
		{
			name: "valid identified objects are normalized",
			invoice: &Invoice{
				IdentifiedObjects: []*Object{
					{
						ID:      "obj1",
						Country: "at",
					},
				},
			},
			wantErrors: 1, // +1 from embedded invoice
		},
		{
			name: "invalid identified object country returns error",
			invoice: &Invoice{
				IdentifiedObjects: []*Object{
					{
						ID:      "obj1",
						Country: "XX",
					},
				},
			},
			wantErrors: 2, // 1 from object + 1 from embedded invoice
		},
		{
			name: "multiple identified objects with errors",
			invoice: &Invoice{
				IdentifiedObjects: []*Object{
					{
						ID:      "obj1",
						Country: "XX",
					},
					{
						ID:      "obj2",
						Country: "YY",
					},
				},
			},
			wantErrors: 3, // 2 from objects + 1 from embedded invoice
		},
		{
			name: "combined errors from section35a and objects",
			invoice: &Invoice{
				Section35aAmounts: []*Section35aInvoiceAmount{
					{
						Type:        "INVALID_TYPE",
						NetAmount:   money.NullableAmountFrom(100),
						GrossAmount: money.NullableAmountFrom(119),
					},
				},
				IdentifiedObjects: []*Object{
					{
						ID:      "obj1",
						Country: "XX",
					},
				},
			},
			wantErrors: 3, // 1 from section35a + 1 from object + 1 from embedded invoice
		},
		{
			name: "nil elements in slices are handled",
			invoice: &Invoice{
				Section35aAmounts: []*Section35aInvoiceAmount{nil},
				IdentifiedObjects: []*Object{nil},
			},
			wantErrors: 1, // +1 from embedded invoice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
		})
	}
}

func TestInvoice_Normalize_Section35aAmounts(t *testing.T) {
	invoice := &Invoice{
		Section35aAmounts: []*Section35aInvoiceAmount{
			{
				Type:        Section35aTypeCraftsmanServices,
				NetAmount:   money.NullableAmountFrom(-100),
				GrossAmount: money.NullableAmountFrom(-119),
			},
		},
	}

	errs := invoice.Normalize()
	if len(errs) != 1 { // 1 error from embedded invoice (invalid payment status)
		t.Errorf("Normalize() got %d errors, want 1: %v", len(errs), errs)
	}

	// Verify amounts were converted to absolute values
	if invoice.Section35aAmounts[0].NetAmount.Get() != 100 {
		t.Errorf("NetAmount = %v, want 100", invoice.Section35aAmounts[0].NetAmount.Get())
	}
	if invoice.Section35aAmounts[0].GrossAmount.Get() != 119 {
		t.Errorf("GrossAmount = %v, want 119", invoice.Section35aAmounts[0].GrossAmount.Get())
	}
}

func TestInvoice_Normalize_IdentifiedObjects(t *testing.T) {
	invoice := &Invoice{
		IdentifiedObjects: []*Object{
			{
				ID:      "obj1",
				Country: "at",
			},
		},
	}

	errs := invoice.Normalize()
	if len(errs) != 1 { // 1 error from embedded invoice (invalid payment status)
		t.Errorf("Normalize() got %d errors, want 1: %v", len(errs), errs)
	}

	// Verify country was normalized to uppercase
	if invoice.IdentifiedObjects[0].Country != "AT" {
		t.Errorf("Country = %q, want %q", invoice.IdentifiedObjects[0].Country, "AT")
	}
}

func TestInvoice_Normalize_RemovesNilElements(t *testing.T) {
	invoice := &Invoice{
		Section35aAmounts: []*Section35aInvoiceAmount{
			nil,
			{
				Type:        Section35aTypeCraftsmanServices,
				NetAmount:   money.NullableAmountFrom(100),
				GrossAmount: money.NullableAmountFrom(119),
			},
			nil,
		},
		IdentifiedObjects: []*Object{
			nil,
			{
				ID:      "obj1",
				Country: "AT",
			},
			nil,
		},
	}

	errs := invoice.Normalize()
	if len(errs) != 1 { // 1 error from embedded invoice (invalid payment status)
		t.Errorf("Normalize() got %d errors, want 1: %v", len(errs), errs)
	}

	// Verify nil elements were removed
	if len(invoice.Section35aAmounts) != 1 {
		t.Errorf("Section35aAmounts length = %d, want 1", len(invoice.Section35aAmounts))
	}
	if len(invoice.IdentifiedObjects) != 1 {
		t.Errorf("IdentifiedObjects length = %d, want 1", len(invoice.IdentifiedObjects))
	}
}
