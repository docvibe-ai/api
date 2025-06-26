package invoicing

import (
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
)

//go:generate go tool go-enum $GOFILE

type PaymentStatus string //#enum,jsonschema

const (
	PaymentStatusUnpaid                          PaymentStatus = "UNPAID"
	PaymentStatusNotPayable                      PaymentStatus = "NOT_PAYABLE"
	PaymentStatusPaidWithCash                    PaymentStatus = "PAID_WITH_CASH"
	PaymentStatusPaidWithCreditcard              PaymentStatus = "PAID_WITH_CREDITCARD"
	PaymentStatusPaidWithBankTransfer            PaymentStatus = "PAID_WITH_BANK_TRANSFER"
	PaymentStatusPaidWithDirectDebit             PaymentStatus = "PAID_WITH_DIRECT_DEBIT"
	PaymentStatusPaidWithStripe                  PaymentStatus = "PAID_WITH_STRIPE"
	PaymentStatusPaidWithPaypal                  PaymentStatus = "PAID_WITH_PAYPAL"
	PaymentStatusPaidWithGooglePay               PaymentStatus = "PAID_WITH_GOOGLE_PAY"
	PaymentStatusPaidWithApplePay                PaymentStatus = "PAID_WITH_APPLE_PAY"
	PaymentStatusPaidWithAmazonPay               PaymentStatus = "PAID_WITH_AMAZON_PAY"
	PaymentStatusPaidWithTransferwise            PaymentStatus = "PAID_WITH_TRANSFERWISE"
	PaymentStatusPaidWithElectronicPaymentMethod PaymentStatus = "PAID_WITH_ELECTRONIC_PAYMENT_METHOD"
)

func (p PaymentStatus) IsPaid() bool {
	return strings.HasPrefix(string(p), "PAID")
}

// Valid indicates if p is any of the valid values for PaymentStatus
func (p PaymentStatus) Valid() bool {
	switch p {
	case
		PaymentStatusUnpaid,
		PaymentStatusNotPayable,
		PaymentStatusPaidWithCash,
		PaymentStatusPaidWithCreditcard,
		PaymentStatusPaidWithBankTransfer,
		PaymentStatusPaidWithDirectDebit,
		PaymentStatusPaidWithStripe,
		PaymentStatusPaidWithPaypal,
		PaymentStatusPaidWithGooglePay,
		PaymentStatusPaidWithApplePay,
		PaymentStatusPaidWithAmazonPay,
		PaymentStatusPaidWithTransferwise,
		PaymentStatusPaidWithElectronicPaymentMethod:
		return true
	}
	return false
}

// Validate returns an error if p is none of the valid values for PaymentStatus
func (p PaymentStatus) Validate() error {
	if !p.Valid() {
		return fmt.Errorf("invalid value %#v for type invoicing.PaymentStatus", p)
	}
	return nil
}

// Enums returns all valid values for PaymentStatus
func (PaymentStatus) Enums() []PaymentStatus {
	return []PaymentStatus{
		PaymentStatusUnpaid,
		PaymentStatusNotPayable,
		PaymentStatusPaidWithCash,
		PaymentStatusPaidWithCreditcard,
		PaymentStatusPaidWithBankTransfer,
		PaymentStatusPaidWithDirectDebit,
		PaymentStatusPaidWithStripe,
		PaymentStatusPaidWithPaypal,
		PaymentStatusPaidWithGooglePay,
		PaymentStatusPaidWithApplePay,
		PaymentStatusPaidWithAmazonPay,
		PaymentStatusPaidWithTransferwise,
		PaymentStatusPaidWithElectronicPaymentMethod,
	}
}

// EnumStrings returns all valid values for PaymentStatus as strings
func (PaymentStatus) EnumStrings() []string {
	return []string{
		"UNPAID",
		"NOT_PAYABLE",
		"PAID_WITH_CASH",
		"PAID_WITH_CREDITCARD",
		"PAID_WITH_BANK_TRANSFER",
		"PAID_WITH_DIRECT_DEBIT",
		"PAID_WITH_STRIPE",
		"PAID_WITH_PAYPAL",
		"PAID_WITH_GOOGLE_PAY",
		"PAID_WITH_APPLE_PAY",
		"PAID_WITH_AMAZON_PAY",
		"PAID_WITH_TRANSFERWISE",
		"PAID_WITH_ELECTRONIC_PAYMENT_METHOD",
	}
}

// String implements the fmt.Stringer interface for PaymentStatus
func (p PaymentStatus) String() string {
	return string(p)
}

// JSONSchema returns a github.com/invopop/jsonschema.Schema for PaymentStatus
func (PaymentStatus) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			"UNPAID",
			"NOT_PAYABLE",
			"PAID_WITH_CASH",
			"PAID_WITH_CREDITCARD",
			"PAID_WITH_BANK_TRANSFER",
			"PAID_WITH_DIRECT_DEBIT",
			"PAID_WITH_STRIPE",
			"PAID_WITH_PAYPAL",
			"PAID_WITH_GOOGLE_PAY",
			"PAID_WITH_APPLE_PAY",
			"PAID_WITH_AMAZON_PAY",
			"PAID_WITH_TRANSFERWISE",
			"PAID_WITH_ELECTRONIC_PAYMENT_METHOD",
		},
	}
}
