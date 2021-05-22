package honestfund

import (
	"encoding/json"
	"github.com/Joddev/autop2p"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Service struct {
	api *Api
}

func BuildService() *Service {
	return &Service{
		api: &Api{
			&http.Client{},
		},
	}
}

func (s *Service) ListProducts() []autop2p.Product {
	resp := s.api.ListProducts(&ListProductRequest{
		Category:     []string{},
		PageSize:     50,
		Scroll:       false,
		State:        []int{2},
		Tendency:     []string{},
		TitleKeyword: "",
	})

	return convertToProducts(resp)
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
		return autop2p.PF
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

func (s *Service) Login(email string, password string) string {
	return s.api.Login(email, password)
}

func (s* Service) CheckAndInvest(accessToken string, productId string, amount int) *autop2p.InvestError {
	err := s.checkInvestment(accessToken, productId, amount)
	if err != nil {
		return err
	}
	productUid, _ := strconv.Atoi(productId)
	s.api.Invest(accessToken, &InvestRequest{
		ProductUid:   productUid,
		InvestAmount: amount,
	})
	return nil
}

func (s* Service) checkInvestment(accessToken string, productId string, amount int) *autop2p.InvestError {
	data := s.api.GetInvestConfirmHtml(accessToken, productId, amount)

	matcher, _ := regexp.Compile("app\\.constant\\('preload', (.+)\\)")

	info := &PreloadInvest{}
	if err := json.Unmarshal(matcher.FindSubmatch(data)[1], info); err != nil {
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

func (s* Service) ListInvestedProductTitles(accessToken string) map[string]struct{} {
	index, pageSize := 0, 25
	totalCount := pageSize + 1

	container := make(map[string]struct{})

	matcher, _ := regexp.Compile("(\\s+(\\d+호))?(\\s+(\\d+차))?$")

	for totalCount > (index+1)*pageSize {
		res := s.api.ListInvestedProduct(accessToken, &ListInvestedProductsRequest{
			Category:     -1,
			Index:        index * pageSize,
			InvestState:  nil,
			IsOngoing:    true,
			PageSize:     pageSize,
			TitleKeyword: "",
		})

		for _, i := range res.Data.Investments {
			container[strings.Trim(matcher.ReplaceAllString(i.Title, ""), " ")] = struct{}{}
		}

		totalCount = res.Data.TotalInvestmentsCount
		index += 1
	}
	return container
}