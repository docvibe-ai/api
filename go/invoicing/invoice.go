package invoicing

import (
	"fmt"
	"slices"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/email"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/vat"
)

type Invoice struct {
	// Type of the invoice
	Type InvoiceType `json:"type,omitempty"`

	// Unique invoice identifier
	InvoiceID nullable.TrimmedString `json:"invoice_id,omitempty"`
	// Issue date of the invoice
	IssueDate date.NullableDate `json:"issue_date,omitempty"`
	// Invoice period start date
	PeriodStart date.NullableDate `json:"period_start,omitempty"`
	// Invoice period end date
	PeriodEnd date.NullableDate `json:"period_end,omitempty"`
	// Due date of the invoice
	DueDate date.NullableDate `json:"due_date,omitempty"`

	// Identifier of the order that the invoice is related to
	OrderID nullable.TrimmedString `json:"order_id,omitempty"`
	// Order date of the invoice
	OrderDate date.NullableDate `json:"order_date,omitempty"`

	// Identifier of the contract that the invoice is related to
	ContractID nullable.TrimmedString `json:"contract_id,omitempty"`
	// Unique customer identifier
	CustomerID nullable.TrimmedString `json:"customer_id,omitempty"`

	// IDs of the delivery notes that are related to the invoice
	DeliveryNoteIDs []notnull.TrimmedString `json:"delivery_note_ids,omitempty"`

	// Issuer of the invoice
	Issuer nullable.TrimmedString `json:"issuer,omitempty"`
	// Issuer's VAT ID
	IssuerVATID vat.NullableID `json:"issuer_vat_id,omitempty"`
	// Issuer's tax number other than VAT ID
	IssuerTaxNumber nullable.TrimmedString `json:"issuer_tax_number,omitempty"`
	// Issuer's address
	IssuerAddress *Address `json:"issuer_address,omitempty"`

	// Recipient of the invoice
	Customer nullable.TrimmedString `json:"customer"`
	// Recipient's VAT ID
	CustomerVATID vat.NullableID `json:"customer_vat_id,omitempty"`
	// Recipient's email
	CustomerEmail email.NullableAddress `json:"customer_email,omitempty"`
	// Recipient's phone
	CustomerPhone nullable.TrimmedString `json:"customer_phone,omitempty"`
	// Recipient's billing address
	CustomerBillingAddress *Address `json:"customer_billing_address,omitempty"`
	// Recipient's shipping address
	CustomerShippingAddress *Address `json:"customer_shipping_address,omitempty"`

	// Subtotal of the invoice
	Subtotal money.NullableAmount `json:"subtotal,omitempty,omitzero"`
	// Tax of the invoice
	Tax money.NullableAmount `json:"tax,omitempty,omitzero"`
	// Total of the invoice
	Total money.NullableAmount `json:"total,omitempty,omitzero"`
	// Currency of the invoice
	Currency money.NullableCurrency `json:"currency,omitempty"`

	// European Union reverse charge for intra-community supply or acquisition
	ReverseCharge bool `json:"reverse_charge"`
	// Reason for the reverse charge value
	ReverseChargeReason nullable.TrimmedString `json:"reverse_charge_reason"`
	// Exact text of the reverse charge clause
	ReverseChargeClauseText nullable.TrimmedString `json:"reverse_charge_clause_text"`
	// Problems indicating that the invoice is not valid for reverse charge, but marked as such
	ReverseChargeProblems nullable.TrimmedString `json:"reverse_charge_problems"`

	// The invoice is a credit note
	CreditNote bool `json:"credit_note"`
	// Exact text of the credit note clause
	CreditNoteClauseText nullable.TrimmedString `json:"credit_note_clause_text"`

	// Payment status of the invoice
	PaymentStatus PaymentStatus `json:"payment_status"`
	// Date the invoice was paid
	PaidDate date.NullableDate `json:"paid_date,omitempty"`
	// Direct debit mandate ID
	DirectDebitMandateID nullable.TrimmedString `json:"direct_debit_mandate_id,omitempty"`

	// Payment reference of the invoice
	PaymentReference nullable.TrimmedString `json:"payment_reference,omitempty"`
	// Payment terms of the invoice
	PaymentTerms nullable.TrimmedString `json:"payment_terms,omitempty"`

	// IBAN of the bank account to pay the invoice
	PaymentIBAN bank.NullableIBAN `json:"payment_iban,omitempty"`
	// SWIFTBIC of the bank account to pay the invoice
	PaymentBIC bank.NullableBIC `json:"payment_bic,omitempty"`

	// Discount percentage of the invoice (valid range: 0-100)
	DiscountPercent money.NullableRate `json:"discount_percent,omitempty,omitzero"`
	// Discount amount of the invoice
	DiscountAmount money.NullableAmount `json:"discount_amount,omitempty,omitzero"`
	// Date until the discount is valid
	DiscountUntilDate date.NullableDate `json:"discount_until_date,omitempty"`

	// Notes of the invoice
	Notes []nullable.TrimmedString `json:"notes,omitempty"`

	// Items in the invoice
	Items []*InvoiceItem `json:"items,omitempty"`

	// Accounting entries of the invoice
	AccountingEntries []*AccountingEntry `json:"accounting_entries,omitempty"`
}

// Normalize validates and normalizes all fields of the Invoice.
// It returns a slice of all validation errors found.
// Invalid fields are either corrected (e.g., negative amounts become absolute)
// or set to null/zero values. The invoice is usable after normalization
// with the returned errors describing what was corrected.
func (inv *Invoice) Normalize() (errs []error) {
	if inv == nil {
		return nil
	}

	if err := inv.Type.Validate(); err != nil {
		errs = append(errs, err)
		inv.Type = ""
	}
	inv.Notes = slices.DeleteFunc(inv.Notes, func(note nullable.TrimmedString) bool {
		return note.IsNull()
	})
	inv.DeliveryNoteIDs = slices.DeleteFunc(inv.DeliveryNoteIDs, func(id notnull.TrimmedString) bool {
		return id.IsEmpty()
	})

	errs = append(errs, inv.normalizeHeaderDates()...)
	errs = append(errs, inv.normalizePeriodDates()...)
	errs = append(errs, inv.normalizeIssuer()...)
	errs = append(errs, inv.normalizeCustomer()...)
	errs = append(errs, inv.normalizeAmountsAndCurrency()...)
	errs = append(errs, inv.normalizePayment()...)
	errs = append(errs, inv.normalizeDiscount()...)
	errs = append(errs, inv.normalizeItems()...)
	errs = append(errs, inv.normalizeAccountingEntries()...)
	return errs
}

func (inv *Invoice) normalizeCustomer() []error {
	var errs []error
	if normalized, err := inv.CustomerVATID.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid customer VAT ID: %w", err))
		inv.CustomerVATID.SetNull()
	} else {
		inv.CustomerVATID = normalized
	}
	if normalized, err := inv.CustomerEmail.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid customer email: %w", err))
		inv.CustomerEmail.SetNull()
	} else {
		inv.CustomerEmail = normalized
	}
	if addrErrs := inv.CustomerBillingAddress.Normalize(); len(addrErrs) > 0 {
		for _, e := range addrErrs {
			errs = append(errs, fmt.Errorf("invalid customer billing address: %w", e))
		}
	}
	if addrErrs := inv.CustomerShippingAddress.Normalize(); len(addrErrs) > 0 {
		for _, e := range addrErrs {
			errs = append(errs, fmt.Errorf("invalid customer shipping address: %w", e))
		}
	}
	return errs
}

func (inv *Invoice) normalizeIssuer() []error {
	var errs []error
	if normalized, err := inv.IssuerVATID.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid issuer VAT ID: %w", err))
		inv.IssuerVATID.SetNull()
	} else {
		inv.IssuerVATID = normalized
	}
	if addrErrs := inv.IssuerAddress.Normalize(); len(addrErrs) > 0 {
		for _, e := range addrErrs {
			errs = append(errs, fmt.Errorf("invalid issuer address: %w", e))
		}
	}
	return errs
}

func (inv *Invoice) normalizeHeaderDates() []error {
	var errs []error
	if normalized, err := inv.IssueDate.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid issue date: %w", err))
		inv.IssueDate.SetNull()
	} else {
		inv.IssueDate = normalized
	}
	if normalized, err := inv.DueDate.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid due date: %w", err))
		inv.DueDate.SetNull()
	} else {
		inv.DueDate = normalized
	}
	if normalized, err := inv.OrderDate.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid order date: %w", err))
		inv.OrderDate.SetNull()
	} else {
		inv.OrderDate = normalized
	}
	return errs
}

func (inv *Invoice) normalizePeriodDates() []error {
	var errs []error
	if normalized, err := inv.PeriodStart.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid period start date: %w", err))
		inv.PeriodStart.SetNull()
	} else {
		inv.PeriodStart = normalized
	}
	if normalized, err := inv.PeriodEnd.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid period end date: %w", err))
		inv.PeriodEnd.SetNull()
	} else {
		inv.PeriodEnd = normalized
	}
	if inv.PeriodStart.IsNotNull() && inv.PeriodEnd.IsNotNull() {
		if inv.PeriodStart.Get().After(inv.PeriodEnd.Get()) {
			errs = append(errs, fmt.Errorf("period start date %s is after period end date %s", inv.PeriodStart.Get(), inv.PeriodEnd.Get()))
			inv.PeriodStart.SetNull()
		}
	}
	return errs
}

func (inv *Invoice) normalizePayment() []error {
	var errs []error
	if err := inv.PaymentStatus.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("invalid payment status: %w", err))
		inv.PaymentStatus = PaymentStatusUnpaid
	}
	if normalized, err := inv.PaidDate.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid paid date: %w", err))
		inv.PaidDate.SetNull()
	} else {
		inv.PaidDate = normalized
	}
	if normalized, err := inv.PaymentIBAN.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid payment IBAN: %w", err))
		inv.PaymentIBAN.SetNull()
	} else {
		inv.PaymentIBAN = normalized
	}
	if normalized, err := inv.PaymentBIC.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid payment BIC: %w", err))
		inv.PaymentBIC.SetNull()
	} else {
		inv.PaymentBIC = normalized
	}
	return errs
}

func (inv *Invoice) normalizeDiscount() []error {
	var errs []error
	if inv.DiscountPercent.IsNotNull() {
		if inv.DiscountPercent.Get() < 0 {
			errs = append(errs, fmt.Errorf("discount percent %f is negative", inv.DiscountPercent.Get()))
			inv.DiscountPercent.Set(inv.DiscountPercent.Get().Abs())
		}
		if inv.DiscountPercent.Get() > 100 {
			errs = append(errs, fmt.Errorf("discount percent %f is greater than 100%%", inv.DiscountPercent.Get()))
			inv.DiscountPercent.SetNull()
		}
	}
	if inv.DiscountAmount.IsNotNull() && inv.DiscountAmount.Get() < 0 {
		errs = append(errs, fmt.Errorf("discount amount %f is negative", inv.DiscountAmount.Get()))
		inv.DiscountAmount.Set(inv.DiscountAmount.Get().Abs())
	}
	if normalized, err := inv.DiscountUntilDate.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid discount until date: %w", err))
		inv.DiscountUntilDate.SetNull()
	} else {
		inv.DiscountUntilDate = normalized
	}
	return errs
}

func (inv *Invoice) normalizeItems() []error {
	var errs []error
	inv.Items = slices.DeleteFunc(inv.Items, func(item *InvoiceItem) bool {
		return item == nil || *item == InvoiceItem{}
	})
	for i, item := range inv.Items {
		if itemErrs := item.Normalize(); len(itemErrs) > 0 {
			for _, e := range itemErrs {
				errs = append(errs, fmt.Errorf("invalid item %d: %w", i, e))
			}
		}
	}
	return errs
}

func (inv *Invoice) normalizeAccountingEntries() []error {
	var errs []error
	inv.AccountingEntries = slices.DeleteFunc(inv.AccountingEntries, func(entry *AccountingEntry) bool {
		return entry == nil || *entry == AccountingEntry{}
	})
	for i, entry := range inv.AccountingEntries {
		if entryErrs := entry.Normalize(); len(entryErrs) > 0 {
			for _, e := range entryErrs {
				errs = append(errs, fmt.Errorf("invalid accounting entry %d: %w", i, e))
			}
		}
	}
	return errs
}

func (inv *Invoice) normalizeAmountsAndCurrency() []error {
	var errs []error
	if inv.Subtotal.IsNotNull() && inv.Subtotal.Get() < 0 {
		errs = append(errs, fmt.Errorf("subtotal %f is negative", inv.Subtotal.Get()))
		inv.Subtotal.Set(inv.Subtotal.Get().Abs())
	}
	if inv.Tax.IsNotNull() && inv.Tax.Get() < 0 {
		errs = append(errs, fmt.Errorf("tax %f is negative", inv.Tax.Get()))
		inv.Tax.Set(inv.Tax.Get().Abs())
	}
	if inv.Total.IsNotNull() && inv.Total.Get() < 0 {
		errs = append(errs, fmt.Errorf("total %f is negative", inv.Total.Get()))
		inv.Total.Set(inv.Total.Get().Abs())
	}
	if inv.Subtotal.IsNotNull() && inv.Total.IsNotNull() {
		if inv.Subtotal.Get() > inv.Total.Get() {
			errs = append(errs, fmt.Errorf("subtotal %f is greater than total %f", inv.Subtotal.Get(), inv.Total.Get()))
			inv.Subtotal.SetNull()
		}
	}
	if inv.Subtotal.IsNotNull() && inv.Tax.IsNotNull() && inv.Total.IsNotNull() {
		if !(inv.Subtotal.Get() + inv.Tax.Get()).WithinOneCent(inv.Total.Get()) {
			errs = append(errs, fmt.Errorf("subtotal %f and tax %f does not sum up to total %f", inv.Subtotal.Get(), inv.Tax.Get(), inv.Total.Get()))
			inv.Tax.SetNull()
		}
	}
	var err error
	if inv.Currency, err = inv.Currency.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid currency: %w", err))
		inv.Currency.SetNull()
	}
	return errs
}

type EUReverseCharge struct {
	ReverseChargeDetected bool                   `json:"reverse_charge_detected"`
	ReverseChargeReason   nullable.TrimmedString `json:"reverse_charge_reason"`
}
