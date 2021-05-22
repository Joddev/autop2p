package honestfund

import (
	"errors"
	"github.com/Joddev/autop2p/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Api interface {
	ListProducts(req *ListProductRequest) *ListProductResponse
	Login(email string, password string) string
	Invest(accessToken string, req *InvestRequest)
	GetInvestConfirmHtml(accessToken string, productId string, amount int) []byte
	ListInvestedProduct(accessToken string, req *ListInvestedProductsRequest) *ListInvestedProductsResponse
}

type ApiImpl struct {
	client *http.Client
}

func NewApi(client *http.Client) Api {
	return &ApiImpl{client}
}

func (a *ApiImpl) ListProducts(req *ListProductRequest) *ListProductResponse {
	resp := util.HandleResponse(a.client.Post(
		"https://www.honestfund.kr/api/search/product/cl",
		"application/json",
		util.EncodeJsonRequest(req),
	))

	ret := &ListProductResponse{}
	util.DecodeJsonResponse(resp, ret)

	return ret
}

type ListProductRequest struct {
	Category     []string `json:"category"`
	PageSize     int      `json:"pageSize"`
	Scroll       bool     `json:"scroll"`
	State        []int    `json:"state"`
	Tendency     []string `json:"tendency"`
	TitleKeyword string   `json:"titleKeyword"`
}

type ListProductResponse struct {
	Code int
	Data struct {
		Products []struct {
			Uid                int
			TitleWithoutSeq    string
			Rate               float64
			Period             int
			GoalAmount         int
			ProgressPercentage float64
			Category           int
		}
	}
}

func (a *ApiImpl) Login(email string, password string) string {
	res := util.HandleResponse(a.client.PostForm(
		"https://www.honestfund.kr/login",
		url.Values{
			"email":             {email},
			"password":          {password},
			"deviceType":        {"1"},
			"next":              {"/"},
			"checkLoginKeeping": {"false"},
		},
	))
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		if cookie.Name == "accessToken" && cookie.Value != "" {
			return cookie.Value
		}
	}
	panic(errors.New("can't find accessToken from cookies"))
}

func (a *ApiImpl) Invest(accessToken string, req *InvestRequest) {
	httpReq, _ := http.NewRequest(
		"POST",
		"https://www.honestfund.kr/invest/confirm",
		util.EncodeJsonRequest(req),
	)

	addJsonContentType(httpReq)
	addAccessTokenCookie(httpReq, accessToken)

	res := util.HandleResponse(a.client.Do(httpReq))
	defer res.Body.Close()
}

type InvestRequest struct {
	ProductUid   int `json:"productUid"`
	InvestAmount int `json:"investAmount"`
}

func (a *ApiImpl) GetInvestConfirmHtml(accessToken string, productId string, amount int) []byte {
	req, _ := http.NewRequest(
		"GET",
		"https://www.honestfund.kr/invest/confirm",
		nil,
	)

	q := req.URL.Query()
	q.Add("productUid", productId)
	q.Add("investAmount", strconv.Itoa(amount))
	req.URL.RawQuery = q.Encode()

	addAccessTokenCookie(req, accessToken)

	res := util.HandleResponse(a.client.Do(req))
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return data
}

func (a *ApiImpl) ListInvestedProduct(accessToken string, req *ListInvestedProductsRequest) *ListInvestedProductsResponse {
	httpReq, _ := http.NewRequest(
		"POST",
		"https://www.honestfund.kr/mypage/investor/investments/search",
		util.EncodeJsonRequest(req),
	)

	addJsonContentType(httpReq)
	addAccessTokenCookie(httpReq, accessToken)

	res := util.HandleResponse(a.client.Do(httpReq))
	defer res.Body.Close()

	data := &ListInvestedProductsResponse{}
	util.DecodeJsonResponse(res, data)

	return data
}

type ListInvestedProductsRequest struct {
	Category     int    `json:"category"`
	Index        int    `json:"index"`
	InvestState  *int   `json:"investState"`
	IsOngoing    bool   `json:"isOngoing"`
	PageSize     int    `json:"pageSize"`
	TitleKeyword string `json:"titleKeyword"`
}

type ListInvestedProductsResponse struct {
	Code int
	Data struct {
		Investments []struct {
			Title string
		}
		TotalInvestmentsCount int
	}
}

func addJsonContentType(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

func addAccessTokenCookie(req *http.Request, accessToken string) {
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: accessToken})
}
