package honestfund

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Joddev/autop2p"
	"github.com/Joddev/autop2p/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func ListProducts() []autop2p.Product {
	resp := util.HandleResponse(http.Post(
		"https://www.honestfund.kr/api/search/product/cl",
		"application/json",
		bytes.NewBuffer(buildListProductRequest()),
	))

	data := &ListProductResponse{}
	util.DecodeResponse(resp, data)

	return convertToProducts(data)
}

type ListProductRequest struct {
	Category     []string `json:"category"`
	PageSize     int      `json:"pageSize"`
	Scroll       bool     `json:"scroll"`
	State        []int    `json:"state"`
	Tendency     []string `json:"tendency"`
	TitleKeyword string   `json:"titleKeyword"`
}

func buildListProductRequest() []byte {
	data, err := json.Marshal(ListProductRequest{
		Category:     []string{},
		PageSize:     50,
		Scroll:       false,
		State:        []int{2},
		Tendency:     []string{},
		TitleKeyword: "",
	})
	if err != nil {
		panic(err)
	}
	return data
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

func convertToProducts(res *ListProductResponse) []autop2p.Product {
	products := make([]autop2p.Product, len(res.Data.Products))
	for i, p := range res.Data.Products {
		products[i] = autop2p.Product{
			Id:           strconv.Itoa(p.Uid),
			Title:        p.TitleWithoutSeq,
			Rate:         p.Rate,
			Period:       p.Period,
			Company:      autop2p.Honestfund,
			RemainAmount: int(float64(p.GoalAmount) * (100 - p.ProgressPercentage)),
			Category:     convertCategory(p.Category),
		}
	}
	return products
}

func convertCategory(category int) autop2p.Category {
	switch category {
	case 1:
		return autop2p.PfRealEstate
	case 2:
		return autop2p.MortgageRealEstate
	case 3:
		return autop2p.CorporateCredit
	case 4:
		return autop2p.PersonalCredit
	default:
		return autop2p.UNKNOWN
	}
}

func Login(email string, password string) (string, error) {
	res := util.HandleResponse(http.PostForm(
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
			return cookie.Value, nil
		}
	}
	return "", errors.New("can't find accessToken from cookies")
}

func CheckAndInvest(accessToken string, productId string, amount int) *autop2p.InvestError {
	err := checkInvestment(accessToken, productId, amount)
	if err != nil {
		return err
	}
	invest(accessToken, productId, amount)
	return nil
}

func invest(accessToken string, productId string, amount int) {
	req, _ := http.NewRequest(
		"POST",
		"https://www.honestfund.kr/invest/confirm",
		bytes.NewBuffer(buildInvestRequest(productId, amount)),
	)

	addJsonContentType(req)
	addAccessTokenCookie(req, accessToken)

	res := util.HandleResponse((&http.Client{}).Do(req))
	defer res.Body.Close()
}

type InvestRequest struct {
	ProductUid   int `json:"productUid"`
	InvestAmount int `json:"investAmount"`
}

func buildInvestRequest(productId string, amount int) []byte {
	productUid, _ := strconv.Atoi(productId)
	data, err := json.Marshal(InvestRequest{
		ProductUid:   productUid,
		InvestAmount: amount,
	})
	if err != nil {
		panic(err)
	}
	return data
}

func checkInvestment(accessToken string, productId string, amount int) *autop2p.InvestError {
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

	res := util.HandleResponse((&http.Client{}).Do(req))
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	matcher, _ := regexp.Compile("app\\.constant\\('preload', (.+)\\)")
	info := &PreloadInvest{}
	if err = json.Unmarshal(matcher.FindSubmatch(data)[1], info); err != nil {
		panic(err)
	}

	if info.Invest.InvestedAmount != 0 {
		return &autop2p.InvestError{Code: autop2p.Duplicated}
	}
	return nil
}

type PreloadInvest struct {
	Invest struct {
		InvestedAmount int
	}
}

const investedProductsUrl = "https://www.honestfund.kr/mypage/investor/investments/search"

func ListInvestedProductTitles(accessToken string) map[string]struct{} {
	index, pageSize := 0, 25
	totalCount := pageSize + 1

	container := make(map[string]struct{})

	matcher, _ := regexp.Compile("(\\s+(\\d+호))?(\\s+(\\d+차))?$")

	for totalCount > (index+1)*pageSize {
		req, _ := http.NewRequest(
			"POST",
			investedProductsUrl,
			bytes.NewBuffer(buildListInvestedProductsRequest(index, pageSize)),
		)

		addJsonContentType(req)
		addAccessTokenCookie(req, accessToken)

		res := util.HandleResponse((&http.Client{}).Do(req))

		data := &ListInvestedProductsResponse{}
		util.DecodeResponse(res, data)

		for _, i := range data.Data.Investments {
			container[strings.Trim(matcher.ReplaceAllString(i.Title, ""), " ")] = struct{}{}
		}

		totalCount = data.Data.TotalInvestmentsCount
		index += 1

		res.Body.Close()
	}
	return container
}

type ListInvestedProductsRequest struct {
	Category     int    `json:"category"`
	Index        int    `json:"index"`
	InvestState  *int   `json:"investState"`
	IsOngoing    bool   `json:"isOngoing"`
	PageSize     int    `json:"pageSize"`
	TitleKeyword string `json:"titleKeyword"`
}

func buildListInvestedProductsRequest(index int, pageSize int) []byte {
	doc, _ := json.Marshal(ListInvestedProductsRequest{
		Category:     -1,
		Index:        index * pageSize,
		InvestState:  nil,
		IsOngoing:    true,
		PageSize:     pageSize,
		TitleKeyword: "",
	})
	return doc
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
