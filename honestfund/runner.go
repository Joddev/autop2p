package honestfund

import (
	"github.com/Joddev/autop2p"
	"strings"
)

type Runner struct {
	accessToken string
	service     Service
}

func NewRunner(setting *autop2p.Setting, service Service) *Runner {
	accessToken := service.Login(setting.Username, setting.Password)

	return &Runner{
		accessToken: accessToken,
		service:     service,
	}
}

func (r *Runner) ListProducts() []autop2p.Product {
	investedProductTitleSet := r.service.ListInvestedProductTitles(r.accessToken)

	var products []autop2p.Product
	for _, product := range r.service.ListProducts() {
		if _, ok := investedProductTitleSet[strings.Trim(product.Title, " ")]; ok {
			if strings.Contains(product.Title, "SCF") {
				products = append(products, product)
			}
		} else {
			products = append(products, product)
		}
	}
	return products
}

func (r *Runner) InvestProduct(product *autop2p.Product, amount int) *autop2p.InvestError {
	return r.service.CheckAndInvest(r.accessToken, product.Id, amount)
}
