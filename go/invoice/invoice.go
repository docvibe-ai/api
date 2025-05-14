package invoice

import (
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/email"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/nullable"
)

type Invoice struct {
	// Type of the invoice, one of: INCOMING_INVOICE, OUTGOING_INVOICE
	Type nullable.NonEmptyString `json:"type,omitempty" jsonschema:"enum=INCOMING_INVOICE,enum=OUTGOING_INVOICE"`
	// Client account number
	ClientAccountNumber nullable.TrimmedString `json:"client_account_number,omitempty"`
	// Vendor account number
	VendorAccountNumber nullable.TrimmedString `json:"vendor_account_number,omitempty"`

	// Unique invoice identifier
	InvoiceID nullable.TrimmedString `json:"invoice_id"`
	// Issue date of the invoice
	IssueDate date.Date `json:"issue_date"`
	// Invoice period start date
	PeriodStart date.NullableDate `json:"period_start,omitempty"`
	// Invoice period end date
	PeriodEnd date.NullableDate `json:"period_end,omitempty"`
	// Due date of the invoice
	DueDate date.NullableDate `json:"due_date,omitempty"`

	// Unique order identifier
	OrderID nullable.TrimmedString `json:"order_id,omitempty"`
	// Unique customer identifier
	CustomerID nullable.TrimmedString `json:"customer_id,omitempty"`
	// Issuer of the invoice
	Issuer nullable.TrimmedString `json:"issuer"`
	// Issuer's address
	IssuerAddress *Address `json:"issuer_address"`

	// Recipient of the invoice
	Customer nullable.TrimmedString `json:"customer"`
	// Recipient's email
	CustomerEmail email.NullableAddress `json:"customer_email,omitempty"`
	// Recipient's phone
	CustomerPhone nullable.TrimmedString `json:"customer_phone,omitempty"`
	// Recipient's billing address
	CustomerBillingAddress *Address `json:"customer_billing_address"`
	// Recipient's shipping address
	CustomerShippingAddress *Address `json:"customer_shipping_address,omitempty"`

	// Subtotal of the invoice
	Subtotal money.NullableAmount `json:"subtotal,omitempty,omitzero"`
	// Tax of the invoice
	Tax money.NullableAmount `json:"tax,omitempty,omitzero"`
	// Total of the invoice
	Total money.Amount `json:"total"`
	// Currency of the invoice
	Currency money.Currency `json:"currency"`

	// The invoice is a credit note
	CreditNote bool `json:"credit_note"`

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
}

type Item struct {
	// Description or name of the item
	Description string `json:"description"`
	// Item is a reverse charge or credit note
	CreditNote bool `json:"credit_note"`
	// Quantity of the item
	Quantity nullable.Type[float64] `json:"quantity,omitempty"`
	// Unit of the item
	Unit nullable.TrimmedString `json:"unit,omitempty"`
	// Unit price of the item
	UnitPrice money.NullableAmount `json:"unit_price,omitempty"`
	// Total price of the item
	TotalPrice money.NullableAmount `json:"total_price,omitempty"`
	// 3-digit currency code
	Currency money.NullableCurrency `json:"currency,omitempty"`
	// Discount percentage of the item
	DiscountPercent money.NullableRate `json:"discount_percent,omitempty,omitzero"`
	// Discount amount of the item
	DiscountAmount money.NullableAmount `json:"discount_amount,omitempty,omitzero"`
	// General Ledger Account Number of the item
	GeneralLedgerAccountNumber nullable.TrimmedString `json:"general_ledger_account_number,omitempty"`
	// Description of the general ledger account
	GeneralLedgerAccountDescription nullable.TrimmedString `json:"general_ledger_account_description,omitempty"`
	// Booking text of the item
	BookingText nullable.TrimmedString `json:"booking_text,omitempty"`
}
