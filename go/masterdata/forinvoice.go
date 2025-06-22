package masterdata

type ForInvoice struct {
	ExtractingCompany Company `json:"extracting_company"`
}

type ForAccountingInvoice struct {
	ExtractingCompany     Company                `json:"extracting_company"`
	PartnerCompanies      []PartnerCompany       `json:"partner_companies,omitempty"`
	GeneralLedgerAccounts []GeneralLedgerAccount `json:"general_ledger_accounts,omitempty"`
}
