package autop2p

type Product struct {
	Id           string
	Company      CompanyType
	Title        string
	Rate         float64
	Period       int
	RemainAmount int
	Category     Category
}
