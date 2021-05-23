package peoplefund

import (
	"github.com/Joddev/autop2p"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) ListProducts() []autop2p.Product {
	args := m.Called()
	return args.Get(0).([]autop2p.Product)
}

func (m *ServiceMock) Login(email string, password string) string {
	args := m.Called(email, password)
	return args.Get(0).(string)
}

func (m *ServiceMock) CheckAndInvest(sessionId string, productId string, amount int) *autop2p.InvestError {
	args := m.Called(sessionId, productId, amount)
	return args.Error(0).(*autop2p.InvestError)
}

func (m *ServiceMock) ListInvestedProductTitles(sessionId string) map[string]struct{} {
	args := m.Called(sessionId)
	return args.Get(0).(map[string]struct{})
}

func TestNewRunner(t *testing.T) {
	m := &ServiceMock{}
	m.On("Login", "hf@peoplefund.kr", "1234password!@#$").Return(
		"SESSION_ID#1414",
	)

	r := NewRunner(&autop2p.Setting{
		Username: "hf@peoplefund.kr",
		Password: "1234password!@#$",
	}, m)

	assert.Equal(t, r.sessionId, "SESSION_ID#1414")
}

func TestRunner_ListProducts(t *testing.T) {
	m := &ServiceMock{}
	m.On("ListInvestedProductTitles", "SESSION_ID#143").Return(
		map[string]struct{}{
			"TITLE#1":      {},
			"Second Title": {},
			"P2P":          {},
		},
	)
	m.On("ListProducts").Return([]autop2p.Product{
		{Title: "TITLE#1"},
		{Title: "P2P-2"},
		{Title: "Third Title"},
	})

	r := Runner{
		sessionId: "SESSION_ID#143",
		service:   m,
	}
	p := r.ListProducts()

	assert.Len(t, p, 2)
	assert.Contains(t, p, autop2p.Product{Title: "Third Title"})
}
