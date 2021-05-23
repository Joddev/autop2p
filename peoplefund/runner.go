package peoplefund

import (
	"github.com/Joddev/autop2p"
	"strings"
)

type Runner struct {
	sessionId string
	service   Service
}

func NewRunner(setting *autop2p.Setting, service Service) *Runner {
	sessionId := service.Login(setting.Username, setting.Password)

	return &Runner{
		sessionId: sessionId,
		service:   service,
	}
}

func (r *Runner) ListProducts() []autop2p.Product {
	investedProductTitleSet := r.service.ListInvestedProductTitles(r.sessionId)

	var products []autop2p.Product
	for _, product := range r.service.ListProducts() {
		if _, ok := investedProductTitleSet[strings.Trim(product.Title, " ")]; !ok {
			products = append(products, product)
		}
	}
	return products
}

func (r *Runner) InvestProduct(product *autop2p.Product, amount int) *autop2p.InvestError {
	return r.service.CheckAndInvest(r.sessionId, product.Id, amount)
}
