package peoplefund

import (
	"encoding/json"
	"github.com/Joddev/autop2p"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ApiMock struct {
	mock.Mock
}

func (m *ApiMock) ListProducts(status string) *ListProductResponse {
	args := m.Called(status)
	return args.Get(0).(*ListProductResponse)
}

func (m *ApiMock) Login(email string, password string) string {
	args := m.Called(email, password)
	return args.Get(0).(string)
}

func (m *ApiMock) Invest(sessionId string, uri string, loanId int, investAmount int, pointAmount int) {
	m.Called(sessionId, uri, loanId, investAmount, pointAmount)
}

func (m *ApiMock) CheckInvestment(sessionId string, loanId int) *CheckInvestmentResponse {
	args := m.Called(sessionId, loanId)
	return args.Get(0).(*CheckInvestmentResponse)
}

func (m *ApiMock) ListInvestedProducts(sessionId string) *ListInvestedProductsResponse {
	args := m.Called(sessionId)
	return args.Get(0).(*ListInvestedProductsResponse)
}

func TestServiceImpl_ListProducts(t *testing.T) {
	jsonString := `{
	  "status": "success",
	  "message": "success",
      "data": {
		"list": [
		  {
			"uri": "ml4980",
			"loan_application_id": 1,
			"loan_type": "아파트담보",
			"detailed_loan_type": "아파트담보",
			"interest_rate": 9,
			"loan_application_term": 12,
			"remain_amount": 100000,
			"loan_title": "아파트 담보(투자시 부자동) 2144"
		  },
		  {
			"uri": "ml5053",
			"loan_application_id": 2,
			"loan_type": "아파트담보",
			"detailed_loan_type": "아파트담보",
			"interest_rate": 9.4,
			"loan_application_term": 9,
			"remain_amount": 2000000,
			"loan_title": "아파트 담보(투자시 벼락동) 2170"
		  }
		]
	  }
	}`
	resp := &ListProductResponse{}
	if err := json.Unmarshal([]byte(jsonString), resp); err != nil {
		panic(err)
	}

	mockApi := &ApiMock{}
	mockApi.On("ListProducts", "투자모집중").Return(resp)

	s := NewService(mockApi)
	p := s.ListProducts()
	assert.Len(t, p, 2)
	assert.Contains(t, p, autop2p.Product{
		Id:           "ml4980-1",
		Company:      autop2p.Peoplefund,
		Title:        "아파트 담보(투자시 부자동) 2144",
		Rate:         9,
		Period:       12,
		RemainAmount: 100000,
		Category:     autop2p.MortgageRealEstate,
	})
	assert.Contains(t, p, autop2p.Product{
		Id:           "ml5053-2",
		Company:      autop2p.Peoplefund,
		Title:        "아파트 담보(투자시 벼락동) 2170",
		Rate:         9.4,
		Period:       9,
		RemainAmount: 2000000,
		Category:     autop2p.MortgageRealEstate,
	})
}

func TestServiceImpl_Login(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("Login", "email", "password").Return("SESSID")

	s := NewService(mockApi)
	sessionId := s.Login("email", "password")

	assert.Equal(t, sessionId, "SESSID")
}

func TestServiceImpl_CheckAndInvest_InsufficientBalance(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("CheckInvestment", "sessionId", 1).Return(&CheckInvestmentResponse{
		Status:  "success",
		Message: "success",
		Data: struct {
			MaxInvestableAmount int `json:"max_investable_amount"`
			Cash                int
		}{
			MaxInvestableAmount: 100000,
			Cash:                0,
		},
	})

	s := NewService(mockApi)
	err := s.CheckAndInvest("sessionId", "ml1-1", 10000)

	assert.Equal(t, err.Code, autop2p.InsufficientBalance)
}

func TestServiceImpl_CheckAndInvest_InsufficientCapacity(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("CheckInvestment", "sessionId", 1).Return(&CheckInvestmentResponse{
		Status:  "success",
		Message: "success",
		Data: struct {
			MaxInvestableAmount int `json:"max_investable_amount"`
			Cash                int
		}{
			MaxInvestableAmount: 0,
			Cash:                100000,
		},
	})

	s := NewService(mockApi)
	err := s.CheckAndInvest("sessionId", "ml1-1", 10000)

	assert.Equal(t, err.Code, autop2p.InsufficientCapacity)
}

func TestServiceImpl_CheckAndInvest(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("CheckInvestment", "sessionId", 1).Return(&CheckInvestmentResponse{
		Status:  "success",
		Message: "success",
		Data: struct {
			MaxInvestableAmount int `json:"max_investable_amount"`
			Cash                int
		}{
			MaxInvestableAmount: 100000,
			Cash:                100000,
		},
	})
	mockApi.On("Invest", "sessionId", "ml1", 1, 10000, 0)

	s := NewService(mockApi)
	err := s.CheckAndInvest("sessionId", "ml1-1", 10000)

	assert.Nil(t, err)
}

func TestServiceImpl_ListInvestedProductTitles(t *testing.T) {
	jsonString := `{
	  "status": "success",
	  "message": "success",
	  "data": {
		"list": [
		  {
			"title": "아파트 담보(투자시 부자동) 2144-1",
            "loan_application_status": "투자모집중" 
		  },
		  {
			"title": "아파트 담보(투자시 부자동) 2144-2",
            "loan_application_status": "투자모집마감" 
		  },
		  {
			"title": "아파트 담보(투자시 손실동) 144-1",
            "loan_application_status": "상환중" 
		  },
		  {
			"title": "아파트 담보(투자시 두배동) 2057",
            "loan_application_status": "매각완료" 
		  },
		  {
			"title": "아파트 담보(투자시 세배동) 1057",
            "loan_application_status": "채권종결" 
		  },
		  {
			"title": "아파트 담보(투자시 세배동) 1057-1",
            "loan_application_status": "상환완료" 
		  },
		  {
			"title": "아파트 담보(투자시 세배동) 1057-2",
            "loan_application_status": "단기지연" 
		  }
		]
	  }
	}`
	resp := &ListInvestedProductsResponse{}
	if err := json.Unmarshal([]byte(jsonString), resp); err != nil {
		panic(err)
	}

	mockApi := &ApiMock{}
	mockApi.On("ListInvestedProducts", mock.Anything).Return(resp)

	s := NewService(mockApi)
	ret := s.ListInvestedProductTitles("sessionId")

	assert.Len(t, ret, 3)
	assert.Contains(t, ret, "아파트 담보(투자시 부자동) 2144")
	assert.Contains(t, ret, "아파트 담보(투자시 손실동) 144")
	assert.Contains(t, ret, "아파트 담보(투자시 세배동) 1057")
}
