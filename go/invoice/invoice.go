package invoice

import (
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/money"
)

type Invoice struct {
	// Unique invoice identifier
	InvoiceID string `json:"invoice_id,omitempty"`
	// Invoice period start date
	PeriodStart date.NullableDate `json:"period_start,omitempty"`
	// Invoice period end date
	PeriodEnd date.NullableDate `json:"period_end,omitempty"`
	// Issue date of the invoice
	InvoiceIssueDate date.NullableDate `json:"invoice_issue_date,omitempty"`
	// Due date of the invoice
	InvoiceDueDate date.NullableDate `json:"invoice_due_date,omitempty"`

	// Unique order identifier
	OrderID string `json:"order_id,omitempty"`
	// Unique customer identifier
	CustomerID string `json:"customer_id,omitempty"`
	// Issuer of the invoice
	Issuer string `json:"issuer,omitempty"`
	// Issuer's address
	IssuerAddress *Address `json:"issuer_address,omitempty"`

	// Recipient of the invoice
	Customer string `json:"customer,omitempty"`
	// Recipient's email
	CustomerEmail string `json:"customer_email,omitempty"`
	// Recipient's phone
	CustomerPhone string `json:"customer_phone,omitempty"`
	// Recipient's billing address
	CustomerBilling string `json:"customer_billing_address,omitempty"`
	// Recipient's shipping address
	CustomerShipping string `json:"customer_shipping_address,omitempty"`

	// Items in the invoice
	Items []Item `json:"items,omitempty"`
	// Subtotal of the invoice
	Subtotal money.NullableAmount `json:"subtotal,omitempty"`
	// Tax of the invoice
	Tax money.NullableAmount `json:"tax,omitempty"`
	// Total of the invoice
	Total money.NullableAmount `json:"total,omitempty"`
	// Currency of the invoice
	Currency money.NullableCurrency `json:"currency,omitempty"`
	// Notes of the invoice
	Notes string `json:"notes,omitempty"`
}

type Item struct {
	// Description or name of the item
	Description string `json:"description,omitempty"`
	// Quantity of the item
	Quantity *float64 `json:"quantity,omitempty"`
	// Unit of the item
	Unit string `json:"unit,omitempty"`
	// Unit price of the item
	UnitPrice money.NullableAmount `json:"unit_price,omitempty"`
	// Total price of the item
	TotalPrice money.NullableAmount `json:"total_price,omitempty"`
	// 3-digit currency code
	Currency money.NullableCurrency `json:"currency,omitempty"`
}
