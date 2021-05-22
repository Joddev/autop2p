package honestfund

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

func (m *ApiMock) ListProducts(req *ListProductRequest) *ListProductResponse {
	args := m.Called(req)
	return args.Get(0).(*ListProductResponse)
}

func (m *ApiMock) Login(email string, password string) string {
	args := m.Called(email, password)
	return args.Get(0).(string)
}

func (m *ApiMock) Invest(accessToken string, req *InvestRequest) {
	m.Called(accessToken, req)
}

func (m *ApiMock) GetInvestConfirmHtml(accessToken string, productId string, amount int) []byte {
	args := m.Called(accessToken, productId, amount)
	return args.Get(0).([]byte)
}

func (m *ApiMock) ListInvestedProduct(accessToken string, req *ListInvestedProductsRequest) *ListInvestedProductsResponse {
	args := m.Called(accessToken, req)
	return args.Get(0).(*ListInvestedProductsResponse)
}

func TestServiceImpl_ListProducts(t *testing.T) {
	jsonString := `{
	  "code": 200,
	  "data": {
		"products": [
		  {
			"uid": 12384,
			"category": 3,
			"titleWithoutSeq": "SCF 플러스",
			"rate": 6.5,
			"period": 2,
			"goalAmount": 500000000,
			"progressPercentage": 10
		  },
		  {
			"uid": 12383,
			"category": 1,
			"titleWithoutSeq": "여수 마리나항만 프리미엄 생활형숙박시설 신축",
			"rate": 13,
			"period": 3,
			"goalAmount": 100000000,
			"progressPercentage": 80
		  }
		]
	  }
	}`
	resp := &ListProductResponse{}
	if err := json.Unmarshal([]byte(jsonString), resp); err != nil {
		panic(err)
	}

	mockApi := &ApiMock{}
	mockApi.On("ListProducts", &ListProductRequest{
		Category:     []string{},
		PageSize:     50,
		Scroll:       false,
		State:        []int{2},
		Tendency:     []string{},
		TitleKeyword: "",
	}).Return(resp)

	s := NewService(mockApi)
	p := s.ListProducts()
	assert.Len(t, p, 2)
	assert.Contains(t, p, autop2p.Product{
		Id:           "12384",
		Company:      autop2p.Honestfund,
		Title:        "SCF 플러스",
		Rate:         6.5,
		Period:       2,
		RemainAmount: 45000000000,
		Category:     autop2p.CorporateCredit,
	})
	assert.Contains(t, p, autop2p.Product{
		Id:           "12383",
		Company:      autop2p.Honestfund,
		Title:        "여수 마리나항만 프리미엄 생활형숙박시설 신축",
		Rate:         13,
		Period:       3,
		RemainAmount: 2000000000,
		Category:     autop2p.PF,
	})
}

func TestServiceImpl_Login(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("Login", "email", "password").Return("ACCESS_TOKEN")

	s := NewService(mockApi)
	accessToken := s.Login("email", "password")

	assert.Equal(t, accessToken, "ACCESS_TOKEN")
}

func TestServiceImpl_CheckAndInvest_Duplicated(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("GetInvestConfirmHtml", "accessToken", "1", 10000).Return([]byte(`
		<!DOCTYPE html>
		<html lang="ko" ng-app="app">
		<head></head>
		<body class="page-invest _confirm">
		  <div>
			<script>
			app = angular.module("app");
			app.constant('preload', {"invest":{"investedAmount":400000,"serviceInvestTerms":true}});
			</script>
		  </div>
		</body>
		</html>
   `))

	s := NewService(mockApi)
	err := s.CheckAndInvest("accessToken", "1", 10000)

	assert.Equal(t, err.Code, autop2p.Duplicated)
}

func TestServiceImpl_CheckAndInvest(t *testing.T) {
	mockApi := &ApiMock{}
	mockApi.On("GetInvestConfirmHtml", "accessToken", "1", 10000).Return([]byte(`
		<!DOCTYPE html>
		<html lang="ko" ng-app="app">
		<head></head>
		<body class="page-invest _confirm">
		  <div>
			<script>
			app = angular.module("app");
			app.constant('preload', {"invest":{"investedAmount":null,"serviceInvestTerms":true}});
			</script>
		  </div>
		</body>
		</html>
   `))
	mockApi.On("Invest", "accessToken", mock.Anything).Return()

	s := NewService(mockApi)
	err := s.CheckAndInvest("accessToken", "1", 10000)

	assert.Nil(t, err)
}

func TestServiceImpl_ListInvestedProductTitles(t *testing.T) {
	jsonString1 := `{
	  "code": 200,
	  "data": {
		"investments": [
		  { "title": "어펀 1호 1차" },
		  { "title": "SCF 베이직 131호" },
          { "title": "중간에 1호 혹은 2차 같은게 들어 있어도 된다 131호" },
		  { "title": "이 페이지에 25개가 들어있는 셈 치자 2호 3차" }
		],
		"totalInvestmentsCount": 27
	  }
	}`
	jsonString2 := `{
	  "code": 200,
	  "data": {
		"investments": [
		  { "title": "여수 마리나 항만  1호 1차" },
		  { "title": "어펀  1호 12차" }
		],
		"totalInvestmentsCount": 27
	  }
	}`
	page1 := &ListInvestedProductsResponse{}
	if err := json.Unmarshal([]byte(jsonString1), page1); err != nil {
		panic(err)
	}
	page2 := &ListInvestedProductsResponse{}
	if err := json.Unmarshal([]byte(jsonString2), page2); err != nil {
		panic(err)
	}

	mockApi := &ApiMock{}
	mockApi.On("ListInvestedProduct", mock.Anything, &ListInvestedProductsRequest{
		Category:     -1,
		Index:        0,
		InvestState:  nil,
		IsOngoing:    true,
		PageSize:     25,
		TitleKeyword: "",
	}).Return(page1)
	mockApi.On("ListInvestedProduct", mock.Anything, &ListInvestedProductsRequest{
		Category:     -1,
		Index:        25,
		InvestState:  nil,
		IsOngoing:    true,
		PageSize:     25,
		TitleKeyword: "",
	}).Return(page2)

	s := NewService(mockApi)
	ret := s.ListInvestedProductTitles("accessToken")

	assert.Len(t, ret, 5)
	assert.Contains(t, ret, "어펀")
	assert.Contains(t, ret, "SCF 베이직")
	assert.Contains(t, ret, "중간에 1호 혹은 2차 같은게 들어 있어도 된다")
	assert.Contains(t, ret, "이 페이지에 25개가 들어있는 셈 치자")
	assert.Contains(t, ret, "여수 마리나 항만")
}
