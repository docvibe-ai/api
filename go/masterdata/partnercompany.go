package masterdata

type Company struct {
	Name             string
	AlternativeNames []string `json:",omitempty"`
	Street           string   `json:",omitempty"`
	City             string   `json:",omitempty"`
	PostalCode       string   `json:",omitempty"`
	Country          string   `json:",omitempty"`
	Phone            string   `json:",omitempty"`
	Email            string   `json:",omitempty"`
	Website          string   `json:",omitempty"`
	VatID            string   `json:",omitempty"`
	RegistrationNo   string   `json:",omitempty"`
}

type PartnerCompany struct {
	ClientAccountNumber string
	VendorAccountNumber string
	Company
}
