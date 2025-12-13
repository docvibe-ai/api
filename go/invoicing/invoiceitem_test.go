package invoicing

import (
	"testing"

	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/nullable"
)

func TestInvoiceItem_Normalize(t *testing.T) {
	tests := []struct {
		name                string
		item                *InvoiceItem
		wantErrors          int
		wantTaxPercentNull  bool
		wantDiscountNull    bool
		wantCurrencyNull    bool
		wantQuantityPositive bool
	}{
		{
			name:       "nil item returns no errors",
			item:       nil,
			wantErrors: 0,
		},
		{
			name:       "empty item is valid",
			item:       &InvoiceItem{},
			wantErrors: 0,
		},
		{
			name: "valid item with all fields",
			item: &InvoiceItem{
				Description:     "Test item",
				Quantity:        nullable.TypeFrom(2.0),
				TaxPercent:      money.NullableRateFrom(19),
				Currency:        "EUR",
				DiscountPercent: money.NullableRateFrom(10),
			},
			wantErrors: 0,
		},
		{
			name: "negative tax percent becomes absolute",
			item: &InvoiceItem{
				TaxPercent: money.NullableRateFrom(-19),
			},
			wantErrors: 0,
		},
		{
			name: "tax percent over 100 is set to null",
			item: &InvoiceItem{
				TaxPercent: money.NullableRateFrom(150),
			},
			wantErrors:         1,
			wantTaxPercentNull: true,
		},
		{
			name: "negative quantity becomes absolute",
			item: &InvoiceItem{
				Quantity: nullable.TypeFrom(-5.0),
			},
			wantErrors:           0,
			wantQuantityPositive: true,
		},
		{
			name: "invalid currency is set to null",
			item: &InvoiceItem{
				Currency: "INVALID",
			},
			wantErrors:       1,
			wantCurrencyNull: true,
		},
		{
			name: "valid currency is normalized",
			item: &InvoiceItem{
				Currency: "eur",
			},
			wantErrors: 0,
		},
		{
			name: "negative discount percent becomes absolute",
			item: &InvoiceItem{
				DiscountPercent: money.NullableRateFrom(-10),
			},
			wantErrors: 0,
		},
		{
			name: "discount percent over 100 is set to null",
			item: &InvoiceItem{
				DiscountPercent: money.NullableRateFrom(150),
			},
			wantErrors:       1,
			wantDiscountNull: true,
		},
		{
			name: "multiple errors",
			item: &InvoiceItem{
				TaxPercent:      money.NullableRateFrom(150),
				Currency:        "INVALID",
				DiscountPercent: money.NullableRateFrom(200),
			},
			wantErrors:         3,
			wantTaxPercentNull: true,
			wantCurrencyNull:   true,
			wantDiscountNull:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.item.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.item == nil {
				return
			}
			if tt.wantTaxPercentNull && tt.item.TaxPercent.IsNotNull() {
				t.Error("TaxPercent should be null after error")
			}
			if tt.wantDiscountNull && tt.item.DiscountPercent.IsNotNull() {
				t.Error("DiscountPercent should be null after error")
			}
			if tt.wantCurrencyNull && tt.item.Currency.IsNotNull() {
				t.Error("Currency should be null after error")
			}
			if tt.wantQuantityPositive && tt.item.Quantity.IsNotNull() && tt.item.Quantity.Get() < 0 {
				t.Error("Quantity should be positive after normalization")
			}
		})
	}
}
