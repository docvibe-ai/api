package masterdata

import (
	"fmt"
	"slices"
	"strings"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/email"
	"github.com/domonda/go-types/vat"
)

type Company struct {
	Name             string                `json:"name"`
	AlternativeNames []string              `json:"alternative_names,omitempty"`
	Street           string                `json:"street,omitempty"`
	City             string                `json:"city,omitempty"`
	PostalCode       string                `json:"postal_code,omitempty"`
	Country          country.NullableCode  `json:"country,omitempty"`
	Phone            string                `json:"phone,omitempty"`
	Email            email.NullableAddress `json:"email,omitempty"`
	Website          string                `json:"website,omitempty"`
	VatID            vat.NullableID        `json:"vat_id,omitempty"`
	RegistrationNo   string                `json:"registration_no,omitempty"`
	BankAccounts     []BankAccount         `json:"bank_accounts,omitempty"`
}

// Normalize validates and normalizes all fields of the Company.
// It returns a slice of all validation errors found.
// Invalid fields are set to null values. The company remains usable
// after normalization, with the returned errors describing what was corrected.
func (c *Company) Normalize() (errs []error) {
	if c == nil {
		return nil
	}

	// Trim string fields
	c.Name = strings.TrimSpace(c.Name)
	c.Street = strings.TrimSpace(c.Street)
	c.City = strings.TrimSpace(c.City)
	c.PostalCode = strings.TrimSpace(c.PostalCode)
	c.Phone = strings.TrimSpace(c.Phone)
	c.Website = strings.TrimSpace(c.Website)
	c.RegistrationNo = strings.TrimSpace(c.RegistrationNo)

	// Remove empty alternative names
	c.AlternativeNames = slices.DeleteFunc(c.AlternativeNames, func(name string) bool {
		return strings.TrimSpace(name) == ""
	})
	for i, name := range c.AlternativeNames {
		c.AlternativeNames[i] = strings.TrimSpace(name)
	}

	// Normalize country code
	var err error
	if c.Country, err = c.Country.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid country code: %w", err))
		c.Country.SetNull()
	}

	// Normalize email
	if c.Email, err = c.Email.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid email: %w", err))
		c.Email.SetNull()
	}

	// Normalize VAT ID
	if c.VatID, err = c.VatID.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid VAT ID: %w", err))
		c.VatID.SetNull()
	}

	// Normalize bank accounts
	validAccounts := make([]BankAccount, 0, len(c.BankAccounts))
	for i := range c.BankAccounts {
		accErrs := c.BankAccounts[i].Normalize()
		for _, e := range accErrs {
			errs = append(errs, fmt.Errorf("bank account %d: %w", i+1, e))
		}
		// Only keep accounts with valid IBAN
		if c.BankAccounts[i].IBAN != "" {
			validAccounts = append(validAccounts, c.BankAccounts[i])
		}
	}
	c.BankAccounts = validAccounts

	return errs
}

type BankAccount struct {
	IBAN bank.IBAN        `json:"iban"`
	BIC  bank.NullableBIC `json:"bic,omitempty"`
}

// Normalize validates and normalizes the BankAccount fields.
// It returns a slice of all validation errors found.
func (b *BankAccount) Normalize() []error {
	if b == nil {
		return nil
	}
	var errs []error
	var err error

	// Normalize IBAN (only if not empty, since bank.IBAN treats empty as invalid)
	if b.IBAN != "" {
		if b.IBAN, err = b.IBAN.Normalized(); err != nil {
			errs = append(errs, fmt.Errorf("invalid IBAN: %w", err))
			b.IBAN = ""
		}
	}

	// Normalize BIC (NullableBIC handles empty as null/valid)
	if b.BIC, err = b.BIC.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("invalid BIC: %w", err))
		b.BIC.SetNull()
	}

	return errs
}

type PartnerCompany struct {
	ClientAccountNumber string `json:"client_account_number"`
	VendorAccountNumber string `json:"vendor_account_number"`
	Company
}
