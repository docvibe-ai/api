package invoicing

import (
	"testing"

	"github.com/domonda/go-types/money"
)

func TestAccountingEntry_Normalize(t *testing.T) {
	tests := []struct {
		name                  string
		entry                 *AccountingEntry
		wantErrors            int
		wantTypeEmpty         bool
		wantGLAccountEmpty    bool
		wantBookingTextEmpty  bool
		wantTaxPercentNull    bool
		wantAmountPositive    bool
		wantTaxAmountPositive bool
	}{
		{
			name:       "nil entry returns no errors",
			entry:      nil,
			wantErrors: 0,
		},
		{
			name:                 "empty entry returns errors for missing required fields",
			entry:                &AccountingEntry{},
			wantErrors:           3, // invalid type, empty GL account, empty booking text
			wantTypeEmpty:        true,
			wantGLAccountEmpty:   true,
			wantBookingTextEmpty: true,
		},
		{
			name: "valid entry with all fields",
			entry: &AccountingEntry{
				Type:                       AccountingEntryTypeCredit,
				GeneralLedgerAccountNumber: "4000",
				Amount:                     100.50,
				TaxAmount:                  money.NullableAmountFrom(19.10),
				TaxPercent:                 money.NullableRateFrom(19),
				BookingText:                "Test booking",
			},
			wantErrors: 0,
		},
		{
			name: "invalid type is set to empty",
			entry: &AccountingEntry{
				Type:                       "INVALID",
				GeneralLedgerAccountNumber: "4000",
				BookingText:                "Test booking",
			},
			wantErrors:    1,
			wantTypeEmpty: true,
		},
		{
			name: "empty general ledger account number returns error",
			entry: &AccountingEntry{
				Type:        AccountingEntryTypeDebit,
				BookingText: "Test booking",
			},
			wantErrors:         1,
			wantGLAccountEmpty: true,
		},
		{
			name: "negative amount becomes absolute and rounded",
			entry: &AccountingEntry{
				Type:                       AccountingEntryTypeCredit,
				GeneralLedgerAccountNumber: "4000",
				Amount:                     -123.456,
				BookingText:                "Test booking",
			},
			wantErrors:         0,
			wantAmountPositive: true,
		},
		{
			name: "negative tax amount becomes absolute and rounded",
			entry: &AccountingEntry{
				Type:                       AccountingEntryTypeDebit,
				GeneralLedgerAccountNumber: "4000",
				Amount:                     100,
				TaxAmount:                  money.NullableAmountFrom(-19.999),
				BookingText:                "Test booking",
			},
			wantErrors:            0,
			wantTaxAmountPositive: true,
		},
		{
			name: "negative tax percent becomes absolute",
			entry: &AccountingEntry{
				Type:                       AccountingEntryTypeCredit,
				GeneralLedgerAccountNumber: "4000",
				Amount:                     100,
				TaxPercent:                 money.NullableRateFrom(-19),
				BookingText:                "Test booking",
			},
			wantErrors: 0,
		},
		{
			name: "tax percent over 100 is set to null",
			entry: &AccountingEntry{
				Type:                       AccountingEntryTypeDebit,
				GeneralLedgerAccountNumber: "4000",
				Amount:                     100,
				TaxPercent:                 money.NullableRateFrom(150),
				BookingText:                "Test booking",
			},
			wantErrors:         1,
			wantTaxPercentNull: true,
		},
		{
			name: "empty booking text returns error",
			entry: &AccountingEntry{
				Type:                       AccountingEntryTypeCredit,
				GeneralLedgerAccountNumber: "4000",
				Amount:                     100,
			},
			wantErrors:           1,
			wantBookingTextEmpty: true,
		},
		{
			name: "multiple errors",
			entry: &AccountingEntry{
				Type:       "INVALID",
				TaxPercent: money.NullableRateFrom(200),
			},
			wantErrors:           4, // invalid type, empty GL account, tax > 100, empty booking text
			wantTypeEmpty:        true,
			wantGLAccountEmpty:   true,
			wantTaxPercentNull:   true,
			wantBookingTextEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.entry.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.entry == nil {
				return
			}
			if tt.wantTypeEmpty && tt.entry.Type != "" {
				t.Error("Type should be empty after error")
			}
			if tt.wantGLAccountEmpty && !tt.entry.GeneralLedgerAccountNumber.IsEmpty() {
				t.Error("GeneralLedgerAccountNumber should be empty after error")
			}
			if tt.wantBookingTextEmpty && !tt.entry.BookingText.IsEmpty() {
				t.Error("BookingText should be empty after error")
			}
			if tt.wantTaxPercentNull && tt.entry.TaxPercent.IsNotNull() {
				t.Error("TaxPercent should be null after error")
			}
			if tt.wantAmountPositive && tt.entry.Amount < 0 {
				t.Error("Amount should be positive after normalization")
			}
			if tt.wantTaxAmountPositive && tt.entry.TaxAmount.IsNotNull() && tt.entry.TaxAmount.Get() < 0 {
				t.Error("TaxAmount should be positive after normalization")
			}
		})
	}
}

func TestAccountingEntry_Normalize_AmountRounding(t *testing.T) {
	entry := &AccountingEntry{
		Type:                       AccountingEntryTypeCredit,
		GeneralLedgerAccountNumber: "4000",
		Amount:                     123.456,
		TaxAmount:                  money.NullableAmountFrom(19.999),
		BookingText:                "Test booking",
	}

	errs := entry.Normalize()
	if len(errs) != 0 {
		t.Errorf("Normalize() got unexpected errors: %v", errs)
	}

	// Amount should be rounded to cents (123.46)
	if entry.Amount != 123.46 {
		t.Errorf("Amount = %v, want 123.46", entry.Amount)
	}

	// TaxAmount should be rounded to cents (20.00)
	if entry.TaxAmount.Get() != 20.00 {
		t.Errorf("TaxAmount = %v, want 20.00", entry.TaxAmount.Get())
	}
}
