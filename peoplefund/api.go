package peoplefund

import (
	"errors"
	"fmt"
	"github.com/Joddev/autop2p/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Api interface {
	ListProducts(status string) *ListProductResponse
	Login(email string, password string) string
	Invest(sessionId string, uri string, loanId int, investAmount int, pointAmount int)
	CheckInvestment(sessionId string, loanId int) *CheckInvestmentResponse
	ListInvestedProducts(sessionId string) *ListInvestedProductsResponse
}

type ApiImpl struct {
	client *http.Client
}

func NewApi(client *http.Client) Api {
	return &ApiImpl{client}
}

func (a *ApiImpl) ListProducts(status string) *ListProductResponse {
	req, _ := http.NewRequest(
		"GET",
		"https://static.peoplefund.co.kr/showcase/newlistGetAjax/1/",
		nil,
	)

	q := req.URL.Query()
	q.Add("status", status)
	req.URL.RawQuery = q.Encode()

	res := util.HandleResponse(a.client.Do(req))

	ret := &ListProductResponse{}
	util.DecodeJsonResponse(res, ret)

	return ret
}

type ListProductResponse struct {
	Status  string
	Message string
	Data    struct {
		List []struct {
			Uri                 string
			LoanApplicationId   int     `json:"loan_application_id"`
			LoanType            string  `json:"loan_type"`
			DetailedLoanType    string  `json:"detailed_loan_type"`
			InterestRate        float64 `json:"interest_rate"`
			LoanApplicationTerm int     `json:"loan_application_term"`
			RemainAmount        int     `json:"remain_amount"`
			LoanTitle           string  `json:"loan_title"`
		}
	}
}

func (a *ApiImpl) Login(email string, password string) string {
	res := util.HandleResponse(a.client.PostForm(
		"https://www.peoplefund.co.kr/auth/loginAjax/",
		url.Values{
			"type":     {"email"},
			"email":    {email},
			"password": {password},
		},
	))
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		if cookie.Name == "SESSID" && cookie.Value != "" {
			return cookie.Value
		}
	}
	panic(errors.New("can't find SESSID from cookies"))
}

func (a *ApiImpl) Invest(sessionId string, uri string, loanId int, investAmount int, pointAmount int) {
	data := url.Values{
		"showcase_uri":        {uri},
		"loan_application_id": {strconv.Itoa(loanId)},
		"invest_amount":       {strconv.Itoa(investAmount)},
		"point_amount":        {strconv.Itoa(pointAmount)},
	}
	httpReq, _ := http.NewRequest(
		"POST",
		"https://www.peoplefund.co.kr/showcase/investSubmitAjax",
		strings.NewReader(data.Encode()),
	)

	addSessionCookie(httpReq, sessionId)

	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	httpReq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res := util.HandleResponse(a.client.Do(httpReq))
	defer res.Body.Close()
}

func (a *ApiImpl) CheckInvestment(sessionId string, loanId int) *CheckInvestmentResponse {
	httpReq, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://www.peoplefund.co.kr/showcase/maxInvestableAmountGetAjax/%d/", loanId),
		nil,
	)

	addSessionCookie(httpReq, sessionId)

	res := util.HandleResponse(a.client.Do(httpReq))
	ret := &CheckInvestmentResponse{}
	util.DecodeJsonResponse(res, ret)

	return ret
}

type CheckInvestmentResponse struct {
	Status  string
	Message string
	Data    struct {
		MaxInvestableAmount int `json:"max_investable_amount"`
		Cash                int
	}
}

func (a *ApiImpl) ListInvestedProducts(sessionId string) *ListInvestedProductsResponse {
	httpReq, _ := http.NewRequest(
		"GET",
		"https://www.peoplefund.co.kr/mypage/investlistAjax?type=showcase",
		nil,
	)

	addSessionCookie(httpReq, sessionId)

	res := util.HandleResponse(a.client.Do(httpReq))
	ret := &ListInvestedProductsResponse{}
	util.DecodeJsonResponse(res, ret)

	return ret
}

type ListInvestedProductsResponse struct {
	Status  string
	Message string
	Data    struct {
		List []struct {
			Uri                   string
			Title                 string
			LoanApplicationId     int    `json:"loan_application_id"`
			LoanType              string `json:"loan_type"`
			LoanApplicationStatus string `json:"loan_application_status"`
		}
	}
}

func addSessionCookie(req *http.Request, sessionId string) {
	req.AddCookie(&http.Cookie{Name: "SESSID", Value: sessionId})
}
