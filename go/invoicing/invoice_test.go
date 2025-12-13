package invoicing

import (
	"testing"

	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
)

func TestInvoice_normalizeHeaderDates(t *testing.T) {
	tests := []struct {
		name       string
		invoice    *Invoice
		wantErrors int
	}{
		{
			name:       "nil invoice returns no errors",
			invoice:    nil,
			wantErrors: 0,
		},
		{
			name: "valid dates",
			invoice: &Invoice{
				IssueDate: "2024-01-15",
				DueDate:   "2024-02-15",
				OrderDate: "2024-01-10",
			},
			wantErrors: 0,
		},
		{
			name:       "null dates are valid",
			invoice:    &Invoice{},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invoice == nil {
				return
			}
			errs := tt.invoice.normalizeHeaderDates()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeHeaderDates() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
		})
	}
}

func TestInvoice_normalizePeriodDates(t *testing.T) {
	tests := []struct {
		name          string
		invoice       *Invoice
		wantErrors    int
		wantStartNull bool
	}{
		{
			name: "valid period",
			invoice: &Invoice{
				PeriodStart: "2024-01-01",
				PeriodEnd:   "2024-01-31",
			},
			wantErrors: 0,
		},
		{
			name: "start after end",
			invoice: &Invoice{
				PeriodStart: "2024-02-01",
				PeriodEnd:   "2024-01-31",
			},
			wantErrors:    1,
			wantStartNull: true,
		},
		{
			name: "same start and end",
			invoice: &Invoice{
				PeriodStart: "2024-01-15",
				PeriodEnd:   "2024-01-15",
			},
			wantErrors: 0,
		},
		{
			name:       "null dates are valid",
			invoice:    &Invoice{},
			wantErrors: 0,
		},
		{
			name: "only start date",
			invoice: &Invoice{
				PeriodStart: "2024-01-01",
			},
			wantErrors: 0,
		},
		{
			name: "only end date",
			invoice: &Invoice{
				PeriodEnd: "2024-01-31",
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizePeriodDates()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizePeriodDates() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.wantStartNull && tt.invoice.PeriodStart.IsNotNull() {
				t.Error("PeriodStart should be null after error")
			}
		})
	}
}

func TestInvoice_normalizeIssuer(t *testing.T) {
	tests := []struct {
		name          string
		invoice       *Invoice
		wantErrors    int
		wantVATIDNull bool
	}{
		{
			name: "valid VAT ID",
			invoice: &Invoice{
				IssuerVATID: "ATU10223006", // Valid Austrian VAT ID
			},
			wantErrors: 0,
		},
		{
			name: "invalid VAT ID",
			invoice: &Invoice{
				IssuerVATID: "INVALID",
			},
			wantErrors:    1,
			wantVATIDNull: true,
		},
		{
			name:       "null VAT ID is valid",
			invoice:    &Invoice{},
			wantErrors: 0,
		},
		{
			name: "invalid address country",
			invoice: &Invoice{
				IssuerAddress: &Address{
					Country: "XX",
				},
			},
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizeIssuer()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeIssuer() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.wantVATIDNull && tt.invoice.IssuerVATID.IsNotNull() {
				t.Error("IssuerVATID should be null after error")
			}
		})
	}
}

func TestInvoice_normalizeCustomer(t *testing.T) {
	tests := []struct {
		name          string
		invoice       *Invoice
		wantErrors    int
		wantVATIDNull bool
		wantEmailNull bool
	}{
		{
			name: "valid customer data",
			invoice: &Invoice{
				CustomerVATID: "ATU10223006", // Valid Austrian VAT ID
				CustomerEmail: "test@example.com",
			},
			wantErrors: 0,
		},
		{
			name: "invalid VAT ID",
			invoice: &Invoice{
				CustomerVATID: "INVALID",
			},
			wantErrors:    1,
			wantVATIDNull: true,
		},
		{
			name: "invalid email",
			invoice: &Invoice{
				CustomerEmail: "not-an-email",
			},
			wantErrors:    1,
			wantEmailNull: true,
		},
		{
			name:       "null values are valid",
			invoice:    &Invoice{},
			wantErrors: 0,
		},
		{
			name: "invalid billing address country",
			invoice: &Invoice{
				CustomerBillingAddress: &Address{
					Country: "XX",
				},
			},
			wantErrors: 1,
		},
		{
			name: "invalid shipping address country",
			invoice: &Invoice{
				CustomerShippingAddress: &Address{
					Country: "XX",
				},
			},
			wantErrors: 1,
		},
		{
			name: "multiple invalid fields",
			invoice: &Invoice{
				CustomerVATID: "INVALID",
				CustomerEmail: "not-an-email",
				CustomerBillingAddress: &Address{
					Country: "XX",
				},
			},
			wantErrors:    3,
			wantVATIDNull: true,
			wantEmailNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizeCustomer()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeCustomer() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.wantVATIDNull && tt.invoice.CustomerVATID.IsNotNull() {
				t.Error("CustomerVATID should be null after error")
			}
			if tt.wantEmailNull && tt.invoice.CustomerEmail.IsNotNull() {
				t.Error("CustomerEmail should be null after error")
			}
		})
	}
}

func TestInvoice_normalizeAmountsAndCurrency(t *testing.T) {
	tests := []struct {
		name             string
		invoice          *Invoice
		wantErrors       int
		wantSubtotalNull bool
		wantTaxNull      bool
		wantCurrencyNull bool
	}{
		{
			name: "valid amounts",
			invoice: &Invoice{
				Subtotal: money.NullableAmountFrom(100),
				Tax:      money.NullableAmountFrom(19),
				Total:    money.NullableAmountFrom(119),
				Currency: "EUR",
			},
			wantErrors: 0,
		},
		{
			name: "negative subtotal corrected",
			invoice: &Invoice{
				Subtotal: money.NullableAmountFrom(-100),
				Total:    money.NullableAmountFrom(119),
			},
			wantErrors: 1,
		},
		{
			name: "negative tax corrected",
			invoice: &Invoice{
				Tax: money.NullableAmountFrom(-19),
			},
			wantErrors: 1,
		},
		{
			name: "negative total corrected",
			invoice: &Invoice{
				Total: money.NullableAmountFrom(-119),
			},
			wantErrors: 1,
		},
		{
			name: "subtotal greater than total",
			invoice: &Invoice{
				Subtotal: money.NullableAmountFrom(150),
				Total:    money.NullableAmountFrom(100),
			},
			wantErrors:       1,
			wantSubtotalNull: true,
		},
		{
			name: "subtotal + tax != total",
			invoice: &Invoice{
				Subtotal: money.NullableAmountFrom(100),
				Tax:      money.NullableAmountFrom(20),
				Total:    money.NullableAmountFrom(119),
			},
			wantErrors:  1,
			wantTaxNull: true,
		},
		{
			name: "subtotal + tax within one cent of total",
			invoice: &Invoice{
				Subtotal: money.NullableAmountFrom(100),
				Tax:      money.NullableAmountFrom(19.005),
				Total:    money.NullableAmountFrom(119),
			},
			wantErrors: 0,
		},
		{
			name: "invalid currency",
			invoice: &Invoice{
				Currency: "INVALID",
			},
			wantErrors:       1,
			wantCurrencyNull: true,
		},
		{
			name:       "null values are valid",
			invoice:    &Invoice{},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizeAmountsAndCurrency()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeAmountsAndCurrency() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.wantSubtotalNull && tt.invoice.Subtotal.IsNotNull() {
				t.Error("Subtotal should be null after error")
			}
			if tt.wantTaxNull && tt.invoice.Tax.IsNotNull() {
				t.Error("Tax should be null after error")
			}
			if tt.wantCurrencyNull && tt.invoice.Currency.IsNotNull() {
				t.Error("Currency should be null after error")
			}
		})
	}
}

func TestInvoice_normalizePayment(t *testing.T) {
	tests := []struct {
		name              string
		invoice           *Invoice
		wantErrors        int
		wantPaymentStatus PaymentStatus
		wantIBANNull      bool
		wantBICNull       bool
	}{
		{
			name: "valid payment data",
			invoice: &Invoice{
				PaymentStatus: PaymentStatusPaidWithBankTransfer,
				PaidDate:      "2024-01-15",
				PaymentIBAN:   "DE89370400440532013000",
				PaymentBIC:    "COBADEFFXXX",
			},
			wantErrors:        0,
			wantPaymentStatus: PaymentStatusPaidWithBankTransfer,
		},
		{
			name: "invalid payment status defaults to unpaid",
			invoice: &Invoice{
				PaymentStatus: "INVALID",
			},
			wantErrors:        1,
			wantPaymentStatus: PaymentStatusUnpaid,
		},
		{
			name: "invalid IBAN",
			invoice: &Invoice{
				PaymentStatus: PaymentStatusUnpaid,
				PaymentIBAN:   "INVALID",
			},
			wantErrors:   1,
			wantIBANNull: true,
		},
		{
			name: "invalid BIC",
			invoice: &Invoice{
				PaymentStatus: PaymentStatusUnpaid,
				PaymentBIC:    "INVALID",
			},
			wantErrors:  1,
			wantBICNull: true,
		},
		{
			name: "empty payment status is invalid",
			invoice: &Invoice{
				PaymentStatus: "",
			},
			wantErrors:        1,
			wantPaymentStatus: PaymentStatusUnpaid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizePayment()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizePayment() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.wantPaymentStatus != "" && tt.invoice.PaymentStatus != tt.wantPaymentStatus {
				t.Errorf("PaymentStatus = %v, want %v", tt.invoice.PaymentStatus, tt.wantPaymentStatus)
			}
			if tt.wantIBANNull && tt.invoice.PaymentIBAN.IsNotNull() {
				t.Error("PaymentIBAN should be null after error")
			}
			if tt.wantBICNull && tt.invoice.PaymentBIC.IsNotNull() {
				t.Error("PaymentBIC should be null after error")
			}
		})
	}
}

func TestInvoice_normalizeDiscount(t *testing.T) {
	tests := []struct {
		name            string
		invoice         *Invoice
		wantErrors      int
		wantPercentNull bool
	}{
		{
			name: "valid discount",
			invoice: &Invoice{
				DiscountPercent:   money.NullableRateFrom(10),
				DiscountAmount:    money.NullableAmountFrom(50),
				DiscountUntilDate: "2024-02-01",
			},
			wantErrors: 0,
		},
		{
			name: "negative discount percent corrected",
			invoice: &Invoice{
				DiscountPercent: money.NullableRateFrom(-10),
			},
			wantErrors: 1,
		},
		{
			name: "discount percent over 100",
			invoice: &Invoice{
				DiscountPercent: money.NullableRateFrom(150),
			},
			wantErrors:      1,
			wantPercentNull: true,
		},
		{
			name: "negative discount amount corrected",
			invoice: &Invoice{
				DiscountAmount: money.NullableAmountFrom(-50),
			},
			wantErrors: 1,
		},
		{
			name:       "null values are valid",
			invoice:    &Invoice{},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizeDiscount()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeDiscount() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.wantPercentNull && tt.invoice.DiscountPercent.IsNotNull() {
				t.Error("DiscountPercent should be null after error")
			}
		})
	}
}

func TestInvoice_normalizeItems(t *testing.T) {
	tests := []struct {
		name           string
		invoice        *Invoice
		wantErrors     int
		wantItemsCount int
	}{
		{
			name: "valid items",
			invoice: &Invoice{
				Items: []*InvoiceItem{
					{Description: "Item 1"},
					{Description: "Item 2"},
				},
			},
			wantErrors:     0,
			wantItemsCount: 2,
		},
		{
			name: "nil item removed",
			invoice: &Invoice{
				Items: []*InvoiceItem{
					{Description: "Item 1"},
					nil,
					{Description: "Item 2"},
				},
			},
			wantErrors:     0,
			wantItemsCount: 2,
		},
		{
			name: "empty item removed",
			invoice: &Invoice{
				Items: []*InvoiceItem{
					{Description: "Item 1"},
					{},
					{Description: "Item 2"},
				},
			},
			wantErrors:     0,
			wantItemsCount: 2,
		},
		{
			name: "item with invalid tax percent",
			invoice: &Invoice{
				Items: []*InvoiceItem{
					{
						Description: "Item 1",
						TaxPercent:  money.NullableRateFrom(150),
					},
				},
			},
			wantErrors:     1,
			wantItemsCount: 1,
		},
		{
			name:           "no items",
			invoice:        &Invoice{},
			wantErrors:     0,
			wantItemsCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizeItems()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeItems() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if len(tt.invoice.Items) != tt.wantItemsCount {
				t.Errorf("Items count = %d, want %d", len(tt.invoice.Items), tt.wantItemsCount)
			}
		})
	}
}

func TestInvoice_normalizeAccountingEntries(t *testing.T) {
	tests := []struct {
		name             string
		invoice          *Invoice
		wantErrors       int
		wantEntriesCount int
	}{
		{
			name: "valid entries",
			invoice: &Invoice{
				AccountingEntries: []*AccountingEntry{
					{
						Type:                       AccountingEntryTypeDebit,
						GeneralLedgerAccountNumber: notnull.TrimmedString("1000"),
						Amount:                     100,
						BookingText:                notnull.TrimmedString("Test booking"),
					},
				},
			},
			wantErrors:       0,
			wantEntriesCount: 1,
		},
		{
			name: "nil entry removed",
			invoice: &Invoice{
				AccountingEntries: []*AccountingEntry{
					{
						Type:                       AccountingEntryTypeDebit,
						GeneralLedgerAccountNumber: notnull.TrimmedString("1000"),
						Amount:                     100,
						BookingText:                notnull.TrimmedString("Test booking"),
					},
					nil,
				},
			},
			wantErrors:       0,
			wantEntriesCount: 1,
		},
		{
			name: "empty entry removed",
			invoice: &Invoice{
				AccountingEntries: []*AccountingEntry{
					{
						Type:                       AccountingEntryTypeDebit,
						GeneralLedgerAccountNumber: notnull.TrimmedString("1000"),
						Amount:                     100,
						BookingText:                notnull.TrimmedString("Test booking"),
					},
					{},
				},
			},
			wantErrors:       0,
			wantEntriesCount: 1,
		},
		{
			name: "entry with missing required fields",
			invoice: &Invoice{
				AccountingEntries: []*AccountingEntry{
					{
						Type:   AccountingEntryTypeDebit,
						Amount: 100,
					},
				},
			},
			wantErrors:       2, // missing GeneralLedgerAccountNumber and BookingText
			wantEntriesCount: 1,
		},
		{
			name:             "no entries",
			invoice:          &Invoice{},
			wantErrors:       0,
			wantEntriesCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.invoice.normalizeAccountingEntries()
			if len(errs) != tt.wantErrors {
				t.Errorf("normalizeAccountingEntries() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if len(tt.invoice.AccountingEntries) != tt.wantEntriesCount {
				t.Errorf("AccountingEntries count = %d, want %d", len(tt.invoice.AccountingEntries), tt.wantEntriesCount)
			}
		})
	}
}

func TestInvoice_Normalize(t *testing.T) {
	tests := []struct {
		name       string
		invoice    *Invoice
		wantErrors int
	}{
		{
			name:       "nil invoice",
			invoice:    nil,
			wantErrors: 0,
		},
		{
			name: "empty invoice has invalid payment status",
			invoice: &Invoice{
				PaymentStatus: PaymentStatusUnpaid, // Set valid status
			},
			wantErrors: 0,
		},
		{
			name: "valid complete invoice",
			invoice: &Invoice{
				Type:          InvoiceTypeIncoming,
				IssueDate:     "2024-01-15",
				DueDate:       "2024-02-15",
				PeriodStart:   "2024-01-01",
				PeriodEnd:     "2024-01-31",
				IssuerVATID:   "ATU10223006", // Valid Austrian VAT ID
				CustomerVATID: "ATU10223006", // Valid Austrian VAT ID
				CustomerEmail: "test@example.com",
				Subtotal:      money.NullableAmountFrom(100),
				Tax:           money.NullableAmountFrom(19),
				Total:         money.NullableAmountFrom(119),
				Currency:      "EUR",
				PaymentStatus: PaymentStatusUnpaid,
				PaymentIBAN:   "DE89370400440532013000",
			},
			wantErrors: 0,
		},
		{
			name: "invoice with multiple errors",
			invoice: &Invoice{
				Type:            "INVALID",
				PeriodStart:     "2024-02-01",
				PeriodEnd:       "2024-01-31",
				IssuerVATID:     "INVALID",
				CustomerVATID:   "INVALID",
				Subtotal:        money.NullableAmountFrom(-100),
				Currency:        "INVALID",
				PaymentStatus:   "INVALID",
				DiscountPercent: money.NullableRateFrom(150),
			},
			wantErrors: 8,
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
