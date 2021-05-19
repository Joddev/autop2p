package autop2p

type Conf struct {
	Settings []Setting
}

type Setting struct {
	Username   string
	Password   string
	Company    CompanyType
	Amount     int
	PeriodMin  int     `yaml:"periodMin"`
	PeriodMax  int     `yaml:"periodMax"`
	RateMin    float64 `yaml:"rateMin"`
	RateMax    float64 `yaml:"rateMax"`
	Categories []Category
}

func (s *Setting) Match(product *Product) bool {
	if s.Amount > product.RemainAmount {
		return false
	}
	if s.PeriodMax < product.Period || s.PeriodMin > product.Period {
		return false
	}
	if s.RateMax < product.Rate || s.RateMin > product.Rate {
		return false
	}
	for _, c := range s.Categories {
		if c == product.Category {
			return true
		}
	}
	return false
}
