package invoice

import (
	"errors"
	"fmt"
	"math"
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
	// Type of the invoice, one of: INCOMING_INVOICE, OUTGOING_INVOICE
	Type Type `json:"type,omitempty" jsonschema:"enum=INCOMING_INVOICE,enum=OUTGOING_INVOICE"`

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

	// Partner account number (vendor or client account number depending on the invoice type)
	PartnerAccountNumber nullable.TrimmedString `json:"partner_account_number,omitempty"`
	// Partner account name (vendor or client name depending on the invoice type)
	PartnerAccountName nullable.TrimmedString `json:"partner_account_name,omitempty"`

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
	CreditNote bool `json:"credit_note,omitzero"`

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
	Items []Item `json:"items,omitempty"`

	// Accounting entries of the invoice
	AccountingEntries []AccountingEntry `json:"accounting_entries,omitempty"`
}

func (inv *Invoice) Normalize() (err error) {
	if inv == nil {
		return nil
	}
	if inv.Type != "" && inv.Type != "INCOMING_INVOICE" && inv.Type != "OUTGOING_INVOICE" {
		err = errors.Join(err, fmt.Errorf("invalid invoice type: %s", inv.Type))
		inv.Type = ""
	}
	if inv.IssueDate, err = inv.IssueDate.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid issue date: %w", err))
		inv.IssueDate.SetNull()
	}
	if inv.PeriodStart, err = inv.PeriodStart.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid period start date: %w", err))
		inv.PeriodStart.SetNull()
	}
	if inv.PeriodEnd, err = inv.PeriodEnd.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid period end date: %w", err))
		inv.PeriodEnd.SetNull()
	}
	if inv.DueDate, err = inv.DueDate.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid due date: %w", err))
		inv.DueDate.SetNull()
	}
	if inv.OrderDate, err = inv.OrderDate.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid order date: %w", err))
		inv.OrderDate.SetNull()
	}
	inv.DeliveryNoteIDs = slices.DeleteFunc(inv.DeliveryNoteIDs, func(id notnull.TrimmedString) bool {
		return id.IsEmpty()
	})
	if inv.IssuerVATID, err = inv.IssuerVATID.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid issuer VAT ID: %w", err))
		inv.IssuerVATID.SetNull()
	}
	if err = inv.IssuerAddress.Normalize(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid issuer address: %w", err))
	}
	if inv.CustomerVATID, err = inv.CustomerVATID.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer VAT ID: %w", err))
		inv.CustomerVATID.SetNull()
	}
	if inv.CustomerEmail, err = inv.CustomerEmail.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer email: %w", err))
		inv.CustomerEmail.SetNull()
	}
	if err = inv.CustomerBillingAddress.Normalize(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer billing address: %w", err))
	}
	if err = inv.CustomerShippingAddress.Normalize(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid customer shipping address: %w", err))
	}
	if inv.Subtotal.IsNotNull() {
		inv.Subtotal.Set(inv.Subtotal.Get().Abs().RoundToCents())
	}
	if inv.Tax.IsNotNull() {
		inv.Tax.Set(inv.Tax.Get().Abs().RoundToCents())
	}
	if inv.Total.IsNotNull() {
		inv.Total.Set(inv.Total.Get().Abs().RoundToCents())
	}
	switch {
	case inv.Subtotal.IsNotNull() && inv.Tax.IsNotNull() && inv.Total.IsNotNull():
		if inv.Subtotal.Get() > inv.Total.Get() {
			err = errors.Join(err, fmt.Errorf("subtotal %f is greater than total %f", inv.Subtotal.Get(), inv.Total.Get()))
			inv.Subtotal, inv.Total = inv.Total, inv.Subtotal
		}
		if !(inv.Subtotal.Get() + inv.Tax.Get()).WithinOneCent(inv.Total.Get()) {
			err = errors.Join(err, fmt.Errorf("subtotal %f and tax %f does not sum up to total %f", inv.Subtotal.Get(), inv.Tax.Get(), inv.Total.Get()))
			inv.Tax.Set(inv.Total.Get() - inv.Subtotal.Get())
		}
	case inv.Subtotal.IsNotNull() && inv.Total.IsNotNull():
		// Tax is null
		if inv.Subtotal.Get() > inv.Total.Get() {
			err = errors.Join(err, fmt.Errorf("subtotal %f is greater than total %f", inv.Subtotal.Get(), inv.Total.Get()))
			inv.Subtotal, inv.Total = inv.Total, inv.Subtotal
		}
		inv.Tax.Set(inv.Total.Get() - inv.Subtotal.Get())
	case inv.Tax.IsNotNull() && inv.Total.IsNotNull():
		// Subtotal is null
		if inv.Tax.Get() > inv.Total.Get() {
			err = errors.Join(err, fmt.Errorf("tax %f is greater than total %f", inv.Tax.Get(), inv.Total.Get()))
			inv.Tax.Set(0)
		}
		inv.Subtotal.Set(inv.Total.Get() - inv.Tax.Get())
	case inv.Subtotal.IsNotNull() && inv.Tax.IsNotNull():
		// Total is null
		inv.Total.Set(inv.Subtotal.Get() + inv.Tax.Get())
	}
	if inv.Currency, err = inv.Currency.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid currency: %w", err))
		inv.Currency.SetNull()
	}
	if inv.PaymentIBAN, err = inv.PaymentIBAN.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid payment IBAN: %w", err))
		inv.PaymentIBAN.SetNull()
	}
	if inv.PaymentBIC, err = inv.PaymentBIC.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid payment BIC: %w", err))
		inv.PaymentBIC.SetNull()
	}
	if inv.DiscountPercent.IsNotNull() {
		inv.DiscountPercent.Set(inv.DiscountPercent.Get().Abs())
		if inv.DiscountPercent.Get() > 100 {
			err = errors.Join(err, fmt.Errorf("discount percent %f is greater than 100%%", inv.DiscountPercent.Get()))
			inv.DiscountPercent.SetNull()
		}
	}
	if inv.DiscountAmount.IsNotNull() {
		inv.DiscountAmount.Set(inv.DiscountAmount.Get().Abs().RoundToCents())
	}
	if inv.DiscountUntilDate, err = inv.DiscountUntilDate.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid discount until date: %w", err))
		inv.DiscountUntilDate.SetNull()
	}
	inv.Notes = slices.DeleteFunc(inv.Notes, func(note nullable.TrimmedString) bool {
		return note.IsNull()
	})
	inv.Items = slices.DeleteFunc(inv.Items, func(item Item) bool {
		return item == Item{}
	})
	for i := range inv.Items {
		if err = inv.Items[i].Normalize(); err != nil {
			err = errors.Join(err, fmt.Errorf("invalid item %d: %w", i, err))
		}
	}
	return err
}

type Item struct {
	// Position number of the item in the invoice
	PositionNumber nullable.TrimmedString `json:"position_number,omitempty"`
	// Description or name of the item
	Description nullable.TrimmedString `json:"description,omitempty"`
	// Item is a reverse charge or credit note
	CreditNote bool `json:"credit_note,omitzero"`
	// Order ID of the item
	OrderID nullable.TrimmedString `json:"order_id,omitempty"`
	// Delivery ID of the item
	DeliveryID nullable.TrimmedString `json:"delivery_id,omitempty"`
	// Product ID of the item
	ProductID nullable.TrimmedString `json:"product_id,omitempty"`
	// Quantity of the item
	Quantity nullable.Type[float64] `json:"quantity,omitempty,omitzero"`
	// Unit of the item
	Unit nullable.TrimmedString `json:"unit,omitempty"`
	// Unit price of the item
	UnitPrice money.NullableAmount `json:"unit_price,omitempty,omitzero"`
	// Total price of the item
	Subtotal money.NullableAmount `json:"subtotal,omitempty,omitzero"`
	// Tax percentage of the item
	TaxPercent money.NullableRate `json:"tax_percent,omitempty,omitzero"`
	// Tax amount of the item
	TaxAmount money.NullableAmount `json:"tax_amount,omitempty,omitzero"`
	// 3-digit currency code
	Currency money.NullableCurrency `json:"currency,omitempty"`
	// Discount percentage of the item
	DiscountPercent money.NullableRate `json:"discount_percent,omitempty,omitzero"`
	// Discount amount of the item
	DiscountAmount money.NullableAmount `json:"discount_amount,omitempty,omitzero"`
}

func (item *Item) Normalize() (err error) {
	if item == nil {
		return nil
	}
	if item.TaxPercent.IsNotNull() {
		item.TaxPercent.Set(item.TaxPercent.Get().Abs())
		if item.TaxPercent.Get() > 100 {
			err = errors.Join(err, fmt.Errorf("tax percent %f is greater than 100%%", item.TaxPercent.Get()))
			item.TaxPercent.SetNull()
		}
	}
	if item.Quantity.IsNotNull() {
		item.Quantity.Set(math.Abs(item.Quantity.Get()))
	}
	if item.Currency, err = item.Currency.Normalized(); err != nil {
		err = errors.Join(err, fmt.Errorf("invalid currency: %w", err))
		item.Currency.SetNull()
	}
	if item.DiscountPercent.IsNotNull() {
		item.DiscountPercent.Set(item.DiscountPercent.Get().Abs())
		if item.DiscountPercent.Get() > 100 {
			err = errors.Join(err, fmt.Errorf("discount percent %f is greater than 100%%", item.DiscountPercent.Get()))
			item.DiscountPercent.SetNull()
		}
	}
	return err
}

type AccountingEntryType string

const (
	AccountingEntryTypeCredit AccountingEntryType = "CREDIT"
	AccountingEntryTypeDebit  AccountingEntryType = "DEBIT"
)

type AccountingEntry struct {
	// Type of the accounting entry
	Type AccountingEntryType `json:"type,omitempty" jsonschema:"enum=CREDIT,enum=DEBIT"`

	// General Ledger Account Number of the item
	GeneralLedgerAccountNumber nullable.TrimmedString `json:"general_ledger_account_number,omitempty"`
	// Description of the general ledger account
	GeneralLedgerAccountDescription nullable.TrimmedString `json:"general_ledger_account_description,omitempty"`
	// Booking text of the item
	BookingText nullable.TrimmedString `json:"booking_text,omitempty"`

	// Date of the accounting entry
	Date date.NullableDate `json:"date,omitempty"`
	// Description of the accounting entry
	Description nullable.TrimmedString `json:"description,omitempty"`
	// Amount of the accounting entry
	Amount money.NullableAmount `json:"amount,omitempty,omitzero"`
	// Tax amount of the accounting entry
	TaxAmount money.NullableAmount `json:"tax_amount,omitempty,omitzero"`
}
