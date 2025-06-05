package masterdata

type Company struct {
	Name             string   `json:"name"`
	AlternativeNames []string `json:"alternative_names,omitempty"`
	Street           string   `json:"street,omitempty"`
	City             string   `json:"city,omitempty"`
	PostalCode       string   `json:"postal_code,omitempty"`
	Country          string   `json:"country,omitempty"`
	Phone            string   `json:"phone,omitempty"`
	Email            string   `json:"email,omitempty"`
	Website          string   `json:"website,omitempty"`
	VatID            string   `json:"vat_id,omitempty"`
	RegistrationNo   string   `json:"registration_no,omitempty"`
}

type PartnerCompany struct {
	ClientAccountNumber string `json:"client_account_number"`
	VendorAccountNumber string `json:"vendor_account_number"`
	Company
}
