package auto

import (
	"fmt"
	"github.com/Joddev/autop2p"
	"github.com/Joddev/autop2p/honestfund"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Run() {
	for _, setting := range loadSettings() {
		runner := buildRunner(&setting)
		products := runner.ListProducts()

		candidates := filter(products, setting)

		count := 0
		for _, p := range candidates {
			err := runner.InvestProduct(&p, setting.Amount)
			if err != nil {
				switch err.Code {
				case autop2p.Duplicated:
					continue
				default:
					panic(err)
				}
			}
			count += 1
		}
		fmt.Printf("%d건 총 투자 금액 %d원", count, setting.Amount*count)
	}
}

func filter(products []autop2p.Product, setting autop2p.Setting) []autop2p.Product {
	var ret []autop2p.Product
	for _, p := range products {
		if setting.Match(&p) {
			ret = append(ret, p)
		}
	}
	return ret
}

func buildRunner(setting *autop2p.Setting) autop2p.Runner {
	switch setting.Company {
	case autop2p.Honestfund:
		return honestfund.Build(setting)
	default:
		panic("unsupported type")
	}
}

func loadSettings() []autop2p.Setting {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		panic(err)
	}

	conf := &autop2p.Conf{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		panic(err)
	}

	return conf.Settings
}
