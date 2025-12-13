package realestate

import (
	"testing"

	"github.com/domonda/go-types/money"
)

func TestSection35aInvoiceAmount_Normalize(t *testing.T) {
	tests := []struct {
		name            string
		amount          *Section35aInvoiceAmount
		wantErrors      int
		wantNetNull     bool
		wantGrossNull   bool
		wantNetAmount   money.Amount
		wantGrossAmount money.Amount
		wantType        Section35aType
	}{
		{
			name:       "nil returns no errors",
			amount:     nil,
			wantErrors: 0,
		},
		{
			name:          "empty struct returns error for invalid type",
			amount:        &Section35aInvoiceAmount{},
			wantErrors:    1,
			wantNetNull:   true,
			wantGrossNull: true,
			wantType:      Section35aTypeCraftsmanServices,
		},
		{
			name: "valid amounts are kept",
			amount: &Section35aInvoiceAmount{
				Type:        Section35aTypeCraftsmanServices,
				NetAmount:   money.NullableAmountFrom(100),
				GrossAmount: money.NullableAmountFrom(119),
			},
			wantErrors:      0,
			wantNetNull:     false,
			wantGrossNull:   false,
			wantNetAmount:   100,
			wantGrossAmount: 119,
			wantType:        Section35aTypeCraftsmanServices,
		},
		{
			name: "negative amounts are converted to absolute",
			amount: &Section35aInvoiceAmount{
				Type:        Section35aTypeHouseholdServices,
				NetAmount:   money.NullableAmountFrom(-50),
				GrossAmount: money.NullableAmountFrom(-59.50),
			},
			wantErrors:      0,
			wantNetNull:     false,
			wantGrossNull:   false,
			wantNetAmount:   50,
			wantGrossAmount: 59.50,
			wantType:        Section35aTypeHouseholdServices,
		},
		{
			name: "zero net amount becomes null",
			amount: &Section35aInvoiceAmount{
				Type:        Section35aTypeCraftsmanServices,
				NetAmount:   money.NullableAmountFrom(0),
				GrossAmount: money.NullableAmountFrom(100),
			},
			wantErrors:      0,
			wantNetNull:     true,
			wantGrossNull:   false,
			wantGrossAmount: 100,
			wantType:        Section35aTypeCraftsmanServices,
		},
		{
			name: "zero gross amount becomes null",
			amount: &Section35aInvoiceAmount{
				Type:        Section35aTypeCraftsmanServices,
				NetAmount:   money.NullableAmountFrom(100),
				GrossAmount: money.NullableAmountFrom(0),
			},
			wantErrors:    0,
			wantNetNull:   false,
			wantGrossNull: true,
			wantNetAmount: 100,
			wantType:      Section35aTypeCraftsmanServices,
		},
		{
			name: "net amount greater than gross sets net to null",
			amount: &Section35aInvoiceAmount{
				Type:        Section35aTypeCraftsmanServices,
				NetAmount:   money.NullableAmountFrom(150),
				GrossAmount: money.NullableAmountFrom(100),
			},
			wantErrors:      1,
			wantNetNull:     true,
			wantGrossNull:   false,
			wantGrossAmount: 100,
			wantType:        Section35aTypeCraftsmanServices,
		},
		{
			name: "net amount equal to gross is valid",
			amount: &Section35aInvoiceAmount{
				Type:        Section35aTypeFormallyEmployedWorker,
				NetAmount:   money.NullableAmountFrom(100),
				GrossAmount: money.NullableAmountFrom(100),
			},
			wantErrors:      0,
			wantNetNull:     false,
			wantGrossNull:   false,
			wantNetAmount:   100,
			wantGrossAmount: 100,
			wantType:        Section35aTypeFormallyEmployedWorker,
		},
		{
			name: "invalid type is set to default",
			amount: &Section35aInvoiceAmount{
				Type:        "INVALID_TYPE",
				NetAmount:   money.NullableAmountFrom(100),
				GrossAmount: money.NullableAmountFrom(119),
			},
			wantErrors:      1,
			wantNetNull:     false,
			wantGrossNull:   false,
			wantNetAmount:   100,
			wantGrossAmount: 119,
			wantType:        Section35aTypeCraftsmanServices,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.amount.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.amount == nil {
				return
			}
			if tt.amount.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", tt.amount.Type, tt.wantType)
			}
			if tt.wantNetNull && tt.amount.NetAmount.IsNotNull() {
				t.Error("NetAmount should be null")
			}
			if !tt.wantNetNull && tt.amount.NetAmount.IsNull() {
				t.Error("NetAmount should not be null")
			}
			if !tt.wantNetNull && tt.amount.NetAmount.Get() != tt.wantNetAmount {
				t.Errorf("NetAmount = %v, want %v", tt.amount.NetAmount.Get(), tt.wantNetAmount)
			}
			if tt.wantGrossNull && tt.amount.GrossAmount.IsNotNull() {
				t.Error("GrossAmount should be null")
			}
			if !tt.wantGrossNull && tt.amount.GrossAmount.IsNull() {
				t.Error("GrossAmount should not be null")
			}
			if !tt.wantGrossNull && tt.amount.GrossAmount.Get() != tt.wantGrossAmount {
				t.Errorf("GrossAmount = %v, want %v", tt.amount.GrossAmount.Get(), tt.wantGrossAmount)
			}
		})
	}
}
