package honestfund

import (
	"github.com/Joddev/autop2p"
	"strings"
)

type Runner struct {
	accessToken             string
	investedProductTitleSet map[string]struct{}
}

func Build(setting *autop2p.Setting) *Runner {
	accessToken, err := Login(setting.Username, setting.Password)
	if err != nil {
		panic(err)
	}

	return &Runner{
		accessToken:             accessToken,
		investedProductTitleSet: ListInvestedProductTitles(accessToken),
	}
}

func (r *Runner) ListProducts() []autop2p.Product {
	var products []autop2p.Product
	for _, product := range ListProducts() {
		if _, ok := r.investedProductTitleSet[strings.Trim(product.Title, " ")]; ok {
			if strings.Contains(product.Title, "SCF") {
				products = append(products, product)
			}
		} else {
			products = append(products, product)
		}
	}
	return products
}

func (r *Runner) ListInvestedProducts() []autop2p.Product {
	return []autop2p.Product{}
}

func (r *Runner) InvestProduct(product *autop2p.Product, amount int) *autop2p.InvestError {
	return CheckAndInvest(r.accessToken, product.Id, amount)
}
