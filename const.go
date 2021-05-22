package autop2p

type CompanyType string

const (
	Honestfund CompanyType = "Honestfund"
)

type Category string

const (
	MortgageRealEstate Category = "MortgageRealEstate"
	CorporateCredit    Category = "CorporateCredit"
	PersonalCredit     Category = "PersonalCredit"
	PF                 Category = "PF"
	UNKNOWN            Category = "UNKNOWN"
)

func (c Category) isRealState() bool {
	return c == MortgageRealEstate || c == PF
}
