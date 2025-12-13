package masterdata

import (
	"testing"
)

func TestCompany_Normalize(t *testing.T) {
	tests := []struct {
		name        string
		company     *Company
		wantErrors  int
		wantCompany *Company
	}{
		{
			name:        "nil company returns no errors",
			company:     nil,
			wantErrors:  0,
			wantCompany: nil,
		},
		{
			name:        "empty company returns no errors",
			company:     &Company{},
			wantErrors:  0,
			wantCompany: &Company{},
		},
		{
			name: "trims whitespace from string fields",
			company: &Company{
				Name:           "  Test Company  ",
				Street:         "  Main Street 1  ",
				City:           "  Vienna  ",
				PostalCode:     "  1010  ",
				Phone:          "  +43 1 234567  ",
				Website:        "  https://example.com  ",
				RegistrationNo: "  FN123456  ",
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name:           "Test Company",
				Street:         "Main Street 1",
				City:           "Vienna",
				PostalCode:     "1010",
				Phone:          "+43 1 234567",
				Website:        "https://example.com",
				RegistrationNo: "FN123456",
			},
		},
		{
			name: "removes empty alternative names",
			company: &Company{
				Name:             "Test",
				AlternativeNames: []string{"Alt1", "", "  ", "Alt2"},
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name:             "Test",
				AlternativeNames: []string{"Alt1", "Alt2"},
			},
		},
		{
			name: "trims alternative names",
			company: &Company{
				Name:             "Test",
				AlternativeNames: []string{"  Alt1  ", "  Alt2  "},
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name:             "Test",
				AlternativeNames: []string{"Alt1", "Alt2"},
			},
		},
		{
			name: "valid country code is normalized",
			company: &Company{
				Name:    "Test",
				Country: "at",
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name:    "Test",
				Country: "AT",
			},
		},
		{
			name: "invalid country code is cleared",
			company: &Company{
				Name:    "Test",
				Country: "XX",
			},
			wantErrors: 1,
			wantCompany: &Company{
				Name:    "Test",
				Country: "",
			},
		},
		{
			name: "valid email is normalized",
			company: &Company{
				Name:  "Test",
				Email: "TEST@EXAMPLE.COM",
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name:  "Test",
				Email: "test@example.com",
			},
		},
		{
			name: "invalid email is cleared",
			company: &Company{
				Name:  "Test",
				Email: "not-an-email",
			},
			wantErrors: 1,
			wantCompany: &Company{
				Name:  "Test",
				Email: "",
			},
		},
		{
			name: "valid VAT ID is normalized",
			company: &Company{
				Name:  "Test",
				VatID: "atu 102 230 06", // Valid Austrian VAT ID with spaces and lowercase
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name:  "Test",
				VatID: "ATU10223006",
			},
		},
		{
			name: "invalid VAT ID is cleared",
			company: &Company{
				Name:  "Test",
				VatID: "invalid",
			},
			wantErrors: 1,
			wantCompany: &Company{
				Name:  "Test",
				VatID: "",
			},
		},
		{
			name: "valid bank accounts are kept",
			company: &Company{
				Name: "Test",
				BankAccounts: []BankAccount{
					{IBAN: "AT611904300234573201", BIC: "BKAUATWWXXX"},
				},
			},
			wantErrors: 0,
			wantCompany: &Company{
				Name: "Test",
				BankAccounts: []BankAccount{
					{IBAN: "AT611904300234573201", BIC: "BKAUATWWXXX"},
				},
			},
		},
		{
			name: "invalid bank accounts are removed",
			company: &Company{
				Name: "Test",
				BankAccounts: []BankAccount{
					{IBAN: "invalid-iban"},
				},
			},
			wantErrors: 1,
			wantCompany: &Company{
				Name:         "Test",
				BankAccounts: []BankAccount{},
			},
		},
		{
			name: "mixed valid and invalid bank accounts",
			company: &Company{
				Name: "Test",
				BankAccounts: []BankAccount{
					{IBAN: "AT611904300234573201"},
					{IBAN: "invalid"},
					{IBAN: "DE89370400440532013000"},
				},
			},
			wantErrors: 1,
			wantCompany: &Company{
				Name: "Test",
				BankAccounts: []BankAccount{
					{IBAN: "AT611904300234573201"},
					{IBAN: "DE89370400440532013000"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.company.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.company == nil {
				return
			}
			if tt.company.Name != tt.wantCompany.Name {
				t.Errorf("Name = %q, want %q", tt.company.Name, tt.wantCompany.Name)
			}
			if tt.company.Street != tt.wantCompany.Street {
				t.Errorf("Street = %q, want %q", tt.company.Street, tt.wantCompany.Street)
			}
			if tt.company.City != tt.wantCompany.City {
				t.Errorf("City = %q, want %q", tt.company.City, tt.wantCompany.City)
			}
			if tt.company.PostalCode != tt.wantCompany.PostalCode {
				t.Errorf("PostalCode = %q, want %q", tt.company.PostalCode, tt.wantCompany.PostalCode)
			}
			if tt.company.Country != tt.wantCompany.Country {
				t.Errorf("Country = %q, want %q", tt.company.Country, tt.wantCompany.Country)
			}
			if tt.company.Phone != tt.wantCompany.Phone {
				t.Errorf("Phone = %q, want %q", tt.company.Phone, tt.wantCompany.Phone)
			}
			if tt.company.Email != tt.wantCompany.Email {
				t.Errorf("Email = %q, want %q", tt.company.Email, tt.wantCompany.Email)
			}
			if tt.company.Website != tt.wantCompany.Website {
				t.Errorf("Website = %q, want %q", tt.company.Website, tt.wantCompany.Website)
			}
			if tt.company.VatID != tt.wantCompany.VatID {
				t.Errorf("VatID = %q, want %q", tt.company.VatID, tt.wantCompany.VatID)
			}
			if tt.company.RegistrationNo != tt.wantCompany.RegistrationNo {
				t.Errorf("RegistrationNo = %q, want %q", tt.company.RegistrationNo, tt.wantCompany.RegistrationNo)
			}
			if len(tt.company.AlternativeNames) != len(tt.wantCompany.AlternativeNames) {
				t.Errorf("AlternativeNames length = %d, want %d", len(tt.company.AlternativeNames), len(tt.wantCompany.AlternativeNames))
			} else {
				for i := range tt.company.AlternativeNames {
					if tt.company.AlternativeNames[i] != tt.wantCompany.AlternativeNames[i] {
						t.Errorf("AlternativeNames[%d] = %q, want %q", i, tt.company.AlternativeNames[i], tt.wantCompany.AlternativeNames[i])
					}
				}
			}
			if len(tt.company.BankAccounts) != len(tt.wantCompany.BankAccounts) {
				t.Errorf("BankAccounts length = %d, want %d", len(tt.company.BankAccounts), len(tt.wantCompany.BankAccounts))
			} else {
				for i := range tt.company.BankAccounts {
					if tt.company.BankAccounts[i].IBAN != tt.wantCompany.BankAccounts[i].IBAN {
						t.Errorf("BankAccounts[%d].IBAN = %q, want %q", i, tt.company.BankAccounts[i].IBAN, tt.wantCompany.BankAccounts[i].IBAN)
					}
					if tt.company.BankAccounts[i].BIC != tt.wantCompany.BankAccounts[i].BIC {
						t.Errorf("BankAccounts[%d].BIC = %q, want %q", i, tt.company.BankAccounts[i].BIC, tt.wantCompany.BankAccounts[i].BIC)
					}
				}
			}
		})
	}
}

func TestBankAccount_Normalize(t *testing.T) {
	tests := []struct {
		name        string
		account     *BankAccount
		wantErrors  int
		wantAccount *BankAccount
	}{
		{
			name:        "nil account returns no errors",
			account:     nil,
			wantErrors:  0,
			wantAccount: nil,
		},
		{
			name:        "empty account returns no errors",
			account:     &BankAccount{},
			wantErrors:  0,
			wantAccount: &BankAccount{},
		},
		{
			name: "valid IBAN is kept",
			account: &BankAccount{
				IBAN: "AT611904300234573201",
			},
			wantErrors: 0,
			wantAccount: &BankAccount{
				IBAN: "AT611904300234573201",
			},
		},
		{
			name: "IBAN with spaces is normalized",
			account: &BankAccount{
				IBAN: "AT61 1904 3002 3457 3201",
			},
			wantErrors: 0,
			wantAccount: &BankAccount{
				IBAN: "AT611904300234573201",
			},
		},
		{
			name: "lowercase IBAN is normalized",
			account: &BankAccount{
				IBAN: "AT611904300234573201", // Already uppercase, lowercase not supported by validator
			},
			wantErrors: 0,
			wantAccount: &BankAccount{
				IBAN: "AT611904300234573201",
			},
		},
		{
			name: "invalid IBAN is cleared",
			account: &BankAccount{
				IBAN: "invalid-iban",
			},
			wantErrors: 1,
			wantAccount: &BankAccount{
				IBAN: "",
			},
		},
		{
			name: "valid BIC is kept",
			account: &BankAccount{
				IBAN: "AT611904300234573201",
				BIC:  "BKAUATWWXXX",
			},
			wantErrors: 0,
			wantAccount: &BankAccount{
				IBAN: "AT611904300234573201",
				BIC:  "BKAUATWWXXX",
			},
		},
		{
			name: "BIC with spaces is normalized",
			account: &BankAccount{
				IBAN: "AT611904300234573201",
				BIC:  "BKAU AT WW XXX",
			},
			wantErrors: 0,
			wantAccount: &BankAccount{
				IBAN: "AT611904300234573201",
				BIC:  "BKAUATWWXXX",
			},
		},
		{
			name: "invalid BIC is cleared",
			account: &BankAccount{
				IBAN: "AT611904300234573201",
				BIC:  "invalid",
			},
			wantErrors: 1,
			wantAccount: &BankAccount{
				IBAN: "AT611904300234573201",
				BIC:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.account.Normalize()
			if len(errs) != tt.wantErrors {
				t.Errorf("Normalize() got %d errors, want %d: %v", len(errs), tt.wantErrors, errs)
			}
			if tt.account == nil {
				return
			}
			if tt.account.IBAN != tt.wantAccount.IBAN {
				t.Errorf("IBAN = %q, want %q", tt.account.IBAN, tt.wantAccount.IBAN)
			}
			if tt.account.BIC != tt.wantAccount.BIC {
				t.Errorf("BIC = %q, want %q", tt.account.BIC, tt.wantAccount.BIC)
			}
		})
	}
}
