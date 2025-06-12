package invoice

import (
	"errors"
	"fmt"
	"math"

	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/nullable"
)

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

func (item *Item) Normalize() error {
	if item == nil {
		return nil
	}
	var e, err error
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
	if item.Currency, e = item.Currency.Normalized(); e != nil {
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
