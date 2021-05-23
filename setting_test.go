package autop2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetting_Match_Amount(t *testing.T) {
	s := &Setting{
		Username:   "username",
		Password:   "password",
		Company:    Honestfund,
		Amount:     10000,
		PeriodMin:  0,
		PeriodMax:  6,
		RateMin:    0,
		RateMax:    12,
		Categories: []Category{PF, CorporateCredit},
	}

	assert.False(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by small period",
		Rate:         3,
		Period:       3,
		RemainAmount: 0,
		Category:     PF,
	}))
}

func TestSetting_Match_Period(t *testing.T) {
	s := &Setting{
		Username:   "username",
		Password:   "password",
		Company:    Honestfund,
		Amount:     10000,
		PeriodMin:  4,
		PeriodMax:  6,
		RateMin:    0,
		RateMax:    12,
		Categories: []Category{PF, CorporateCredit},
	}

	assert.False(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by small period",
		Rate:         3,
		Period:       3,
		RemainAmount: 100000,
		Category:     PF,
	}))

	assert.False(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by long period",
		Rate:         3,
		Period:       7,
		RemainAmount: 100000,
		Category:     PF,
	}))
}

func TestSetting_Match_Rate(t *testing.T) {
	s := &Setting{
		Username:   "username",
		Password:   "password",
		Company:    Honestfund,
		Amount:     10000,
		PeriodMin:  0,
		PeriodMax:  6,
		RateMin:    4,
		RateMax:    12,
		Categories: []Category{PF, CorporateCredit},
	}

	assert.False(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by small rate",
		Rate:         3,
		Period:       3,
		RemainAmount: 100000,
		Category:     PF,
	}))

	assert.False(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by big rate",
		Rate:         13,
		Period:       7,
		RemainAmount: 100000,
		Category:     PF,
	}))
}

func TestSetting_Match_Category(t *testing.T) {
	s := &Setting{
		Username:   "username",
		Password:   "password",
		Company:    Honestfund,
		Amount:     10000,
		PeriodMin:  0,
		PeriodMax:  6,
		RateMin:    0,
		RateMax:    12,
		Categories: []Category{CorporateCredit},
	}

	assert.False(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by amount",
		Rate:         3,
		Period:       6,
		RemainAmount: 0,
		Category:     PF,
	}))
}

func TestSetting_Match(t *testing.T) {
	s := &Setting{
		Username:   "username",
		Password:   "password",
		Company:    Honestfund,
		Amount:     10000,
		PeriodMin:  0,
		PeriodMax:  6,
		RateMin:    0,
		RateMax:    12,
		Categories: []Category{PF, CorporateCredit},
	}

	assert.True(t, s.Match(&Product{
		Id:           "1",
		Company:      Honestfund,
		Title:        "Filtered by amount",
		Rate:         3,
		Period:       6,
		RemainAmount: 1000000,
		Category:     PF,
	}))
}
