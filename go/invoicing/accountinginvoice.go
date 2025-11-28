package invoicing

//go:generate go tool go-enum $GOFILE

import (
	"errors"
	"fmt"

	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/invopop/jsonschema"
)

type AccountingInvoice struct {
	Invoice

	// Partner account number (vendor or client account number depending on the invoice type)
	PartnerAccountNumber nullable.TrimmedString `json:"partner_account_number,omitempty"`
	// Partner account name (vendor or client name depending on the invoice type)
	PartnerAccountName nullable.TrimmedString `json:"partner_account_name,omitempty"`

	AccountingEntries []*AccountingEntry `json:"accounting_entries,omitempty"`
}

type AccountingEntry struct {
	// Type of the accounting entry
	Type AccountingEntryType `json:"type"`

	// General Ledger Account Number of the item
	GeneralLedgerAccountNumber notnull.TrimmedString `json:"general_ledger_account_number"`
	// Description of the general ledger account
	GeneralLedgerAccountDescription nullable.TrimmedString `json:"general_ledger_account_description,omitempty"`

	// Amount of the accounting entry
	Amount money.Amount `json:"amount"`
	// Tax amount of the accounting entry
	TaxAmount money.NullableAmount `json:"tax_amount,omitempty,omitzero"`
	// Tax percentage of the accounting entry
	TaxPercent money.NullableRate `json:"tax_percent,omitempty,omitzero"`

	// Booking text of the item
	BookingText notnull.TrimmedString `json:"booking_text"`
}

// Normalize validates and normalizes all fields of the AccountingEntry.
// It returns an aggregated error of all validation issues found.
// Invalid fields are either corrected (e.g., amounts become absolute and rounded)
// or set to null/zero values. The entry remains usable after normalization,
// with the returned error describing what was corrected.
func (a *AccountingEntry) Normalize() error {
	if a == nil {
		return nil
	}
	var err, result error
	if err = a.Type.Validate(); err != nil {
		result = errors.Join(result, err)
		a.Type = ""
	}
	if a.GeneralLedgerAccountNumber.IsEmpty() {
		result = errors.Join(result, fmt.Errorf("general ledger account number is empty"))
		a.GeneralLedgerAccountNumber = ""
	}
	a.Amount = a.Amount.Abs().RoundToCents()
	if a.TaxAmount.IsNotNull() {
		a.TaxAmount.Set(a.TaxAmount.Get().Abs().RoundToCents())
	}
	if a.TaxPercent.IsNotNull() {
		a.TaxPercent.Set(a.TaxPercent.Get().Abs())
		if a.TaxPercent.Get() > 100 {
			result = errors.Join(result, fmt.Errorf("tax percent %f is greater than 100%%", a.TaxPercent.Get()))
			a.TaxPercent.SetNull()
		}
	}
	if a.BookingText.IsEmpty() {
		result = errors.Join(result, fmt.Errorf("booking text is empty"))
		a.BookingText = ""
	}
	return result
}

type AccountingEntryType string //#enum,jsonschema

const (
	AccountingEntryTypeCredit AccountingEntryType = "CREDIT"
	AccountingEntryTypeDebit  AccountingEntryType = "DEBIT"
)

// Valid indicates if t is any of the valid values for AccountingEntryType
func (t AccountingEntryType) Valid() bool {
	switch t {
	case
		AccountingEntryTypeCredit,
		AccountingEntryTypeDebit:
		return true
	}
	return false
}

// Validate returns an error if t is none of the valid values for AccountingEntryType
func (t AccountingEntryType) Validate() error {
	if !t.Valid() {
		return fmt.Errorf("invalid value %#v for type invoicing.AccountingEntryType", t)
	}
	return nil
}

// Enums returns all valid values for AccountingEntryType
func (AccountingEntryType) Enums() []AccountingEntryType {
	return []AccountingEntryType{
		AccountingEntryTypeCredit,
		AccountingEntryTypeDebit,
	}
}

// EnumStrings returns all valid values for AccountingEntryType as strings
func (AccountingEntryType) EnumStrings() []string {
	return []string{
		"CREDIT",
		"DEBIT",
	}
}

// String implements the fmt.Stringer interface for AccountingEntryType
func (t AccountingEntryType) String() string {
	return string(t)
}

// JSONSchema returns a github.com/invopop/jsonschema.Schema for AccountingEntryType
func (AccountingEntryType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			"CREDIT",
			"DEBIT",
		},
	}
}
