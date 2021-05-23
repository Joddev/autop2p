package peoplefund

import (
	"fmt"
	"github.com/Joddev/autop2p"
	"regexp"
	"strconv"
	"strings"
)

type Service interface {
	ListProducts() []autop2p.Product
	Login(email string, password string) string
	CheckAndInvest(sessionId string, productId string, amount int) *autop2p.InvestError
	ListInvestedProductTitles(sessionId string) map[string]struct{}
}

type ServiceImpl struct {
	api Api
}

func NewService(api Api) Service {
	return &ServiceImpl{api}
}

func (s *ServiceImpl) ListProducts() []autop2p.Product {
	resp := s.api.ListProducts("투자모집중")

	return convertToProducts(resp)
}

func convertToProducts(res *ListProductResponse) []autop2p.Product {
	products := make([]autop2p.Product, len(res.Data.List))
	for i, p := range res.Data.List {
		products[i] = autop2p.Product{
			Id:           fmt.Sprintf("%s-%d", p.Uri, p.LoanApplicationId),
			Title:        p.LoanTitle,
			Rate:         p.InterestRate,
			Period:       p.LoanApplicationTerm,
			Company:      autop2p.Peoplefund,
			RemainAmount: p.RemainAmount,
			Category:     convertCategory(p.LoanType),
		}
	}
	return products
}

func convertCategory(loanType string) autop2p.Category {
	switch loanType {
	case "아파트담보":
		return autop2p.MortgageRealEstate
	default:
		return autop2p.UNKNOWN
	}
}

func (s *ServiceImpl) Login(email string, password string) string {
	return s.api.Login(email, password)
}

func (s *ServiceImpl) CheckAndInvest(sessionId string, productId string, amount int) *autop2p.InvestError {
	slice := strings.Split(productId, "-")
	loanId, _ := strconv.Atoi(slice[1])
	err := s.checkInvestment(sessionId, loanId, amount)
	if err != nil {
		return err
	}
	s.api.Invest(sessionId, slice[0], loanId, amount, 0)
	return nil
}

func (s *ServiceImpl) checkInvestment(sessionId string, loanId int, amount int) *autop2p.InvestError {
	info := s.api.CheckInvestment(sessionId, loanId)

	if info.Data.Cash < amount {
		return &autop2p.InvestError{Code: autop2p.InsufficientBalance}
	}
	if info.Data.MaxInvestableAmount < amount {
		return &autop2p.InvestError{Code: autop2p.InsufficientCapacity}
	}
	return nil
}

func (s *ServiceImpl) ListInvestedProductTitles(sessionId string) map[string]struct{} {
	list := s.api.ListInvestedProducts(sessionId)

	container := make(map[string]struct{})

	matcher, _ := regexp.Compile("(-\\d+)?$")

	filter := map[string]struct{}{
		"매각완료": {}, "채권종결": {}, "상환완료": {},
	}

	for _, p := range list.Data.List {
		if _, ok := filter[p.LoanApplicationStatus]; !ok {
			container[strings.Trim(matcher.ReplaceAllString(p.Title, ""), " ")] = struct{}{}
		}
	}
	return container
}
