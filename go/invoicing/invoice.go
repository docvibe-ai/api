package invoicing

import (
	"errors"
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

	// European Union reverse charge
	ReverseCharge bool `json:"reverse_charge"`
	// Intra-Community supply
	IntraCommunitySupply bool `json:"intra_community_supply"`

	// The invoice is a credit note
	CreditNote bool `json:"credit_note"`

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

	// Discount percentage of the invoice
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

func (inv *Invoice) Normalize() error {
	if inv == nil {
		return nil
	}
	var e, err error
	if e = inv.Type.Validate(); e != nil {
		err = errors.Join(err, e)
		inv.Type = ""
	}
	if inv.IssueDate, e = inv.IssueDate.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid issue date: %w", e))
		inv.IssueDate.SetNull()
	}
	if inv.PeriodStart, e = inv.PeriodStart.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid period start date: %w", e))
		inv.PeriodStart.SetNull()
	}
	if inv.PeriodEnd, e = inv.PeriodEnd.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid period end date: %w", e))
		inv.PeriodEnd.SetNull()
	}
	if inv.PeriodStart.IsNotNull() && inv.PeriodEnd.IsNotNull() {
		if inv.PeriodStart.Get().After(inv.PeriodEnd.Get()) {
			err = errors.Join(err, fmt.Errorf("period start date %s is after period end date %s", inv.PeriodStart.Get(), inv.PeriodEnd.Get()))
			inv.PeriodStart.SetNull()
		}
	}
	if inv.DueDate, e = inv.DueDate.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid due date: %w", e))
		inv.DueDate.SetNull()
	}
	if inv.OrderDate, e = inv.OrderDate.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid order date: %w", e))
		inv.OrderDate.SetNull()
	}
	inv.DeliveryNoteIDs = slices.DeleteFunc(inv.DeliveryNoteIDs, func(id notnull.TrimmedString) bool {
		return id.IsEmpty()
	})
	if inv.IssuerVATID, e = inv.IssuerVATID.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid issuer VAT ID: %w", e))
		inv.IssuerVATID.SetNull()
	}
	if err = inv.IssuerAddress.Normalize(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid issuer address: %w", e))
	}
	if inv.CustomerVATID, e = inv.CustomerVATID.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer VAT ID: %w", e))
		inv.CustomerVATID.SetNull()
	}
	if inv.CustomerEmail, e = inv.CustomerEmail.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer email: %w", e))
		inv.CustomerEmail.SetNull()
	}
	if err = inv.CustomerBillingAddress.Normalize(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer billing address: %w", e))
	}
	if err = inv.CustomerShippingAddress.Normalize(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer shipping address: %w", e))
	}

	if inv.Subtotal.IsNotNull() && inv.Subtotal.Get() < 0 {
		err = errors.Join(err, fmt.Errorf("subtotal %f is negative", inv.Subtotal.Get()))
		inv.Subtotal.Set(inv.Subtotal.Get().Abs())
	}
	if inv.Tax.IsNotNull() && inv.Tax.Get() < 0 {
		err = errors.Join(err, fmt.Errorf("tax %f is negative", inv.Tax.Get()))
		inv.Tax.Set(inv.Tax.Get().Abs())
	}
	if inv.Total.IsNotNull() && inv.Total.Get() < 0 {
		err = errors.Join(err, fmt.Errorf("total %f is negative", inv.Total.Get()))
		inv.Total.Set(inv.Total.Get().Abs())
	}
	if inv.Subtotal.IsNotNull() && inv.Total.IsNotNull() {
		if inv.Subtotal.Get() > inv.Total.Get() {
			err = errors.Join(err, fmt.Errorf("subtotal %f is greater than total %f", inv.Subtotal.Get(), inv.Total.Get()))
			inv.Subtotal.SetNull()
		}
	}
	if inv.Subtotal.IsNotNull() && inv.Tax.IsNotNull() && inv.Total.IsNotNull() {
		if !(inv.Subtotal.Get() + inv.Tax.Get()).WithinOneCent(inv.Total.Get()) {
			err = errors.Join(err, fmt.Errorf("subtotal %f and tax %f does not sum up to total %f", inv.Subtotal.Get(), inv.Tax.Get(), inv.Total.Get()))
			inv.Tax.SetNull()
		}
	}
	if inv.Currency, e = inv.Currency.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid currency: %w", e))
		inv.Currency.SetNull()
	}

	if e = inv.PaymentStatus.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid payment status: %w", e))
		inv.PaymentStatus = PaymentStatusUnpaid
	}
	if inv.PaidDate, e = inv.PaidDate.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid paid date: %w", e))
		inv.PaidDate.SetNull()
	}
	if inv.PaymentIBAN, e = inv.PaymentIBAN.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid payment IBAN: %w", e))
		inv.PaymentIBAN.SetNull()
	}
	if inv.PaymentBIC, e = inv.PaymentBIC.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid payment BIC: %w", e))
		inv.PaymentBIC.SetNull()
	}
	if inv.DiscountPercent.IsNotNull() {
		if inv.DiscountPercent.Get() < 0 {
			err = errors.Join(err, fmt.Errorf("discount percent %f is negative", inv.DiscountPercent.Get()))
			inv.DiscountPercent.Set(inv.DiscountPercent.Get().Abs())
		}
		if inv.DiscountPercent.Get() > 100 {
			err = errors.Join(err, fmt.Errorf("discount percent %f is greater than 100%%", inv.DiscountPercent.Get()))
			inv.DiscountPercent.SetNull()
		}
	}
	if inv.DiscountAmount.IsNotNull() && inv.DiscountAmount.Get() < 0 {
		err = errors.Join(err, fmt.Errorf("discount amount %f is negative", inv.DiscountAmount.Get()))
		inv.DiscountAmount.Set(inv.DiscountAmount.Get().Abs())
	}
	if inv.DiscountUntilDate, e = inv.DiscountUntilDate.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid discount until date: %w", e))
		inv.DiscountUntilDate.SetNull()
	}
	inv.Notes = slices.DeleteFunc(inv.Notes, func(note nullable.TrimmedString) bool {
		return note.IsNull()
	})
	inv.Items = slices.DeleteFunc(inv.Items, func(item *InvoiceItem) bool {
		return item == nil || *item == InvoiceItem{}
	})
	for i, item := range inv.Items {
		if e = item.Normalize(); e != nil {
			err = errors.Join(err, fmt.Errorf("invalid item %d: %w", i, e))
		}
	}
	inv.AccountingEntries = slices.DeleteFunc(inv.AccountingEntries, func(entry *AccountingEntry) bool {
		return entry == nil || *entry == AccountingEntry{}
	})
	for i, entry := range inv.AccountingEntries {
		if e = entry.Normalize(); e != nil {
			err = errors.Join(err, fmt.Errorf("invalid accounting entry %d: %w", i, e))
		}
	}
	return err
}
