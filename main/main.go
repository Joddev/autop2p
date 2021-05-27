package main

import (
	"context"
	"fmt"
	"github.com/Joddev/autop2p"
	"github.com/Joddev/autop2p/honestfund"
	"github.com/Joddev/autop2p/peoplefund"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func Run(ctx context.Context) {
	auto()
	ctx.Done()
}

func auto() {
	for _, setting := range loadSettings() {
		runner := newRunner(&setting)
		products := runner.ListProducts()

		candidates := filter(products, setting)

		count := 0
		for _, p := range candidates {
			err := runner.InvestProduct(&p, setting.Amount)
			if err != nil {
				switch err.Code {
				case autop2p.Duplicated:
				case autop2p.InsufficientCapacity:
					continue
				case autop2p.InsufficientBalance:
					break
				default:
					panic(err)
				}
			} else {
				count += 1
			}
		}
		fmt.Printf("%s %s %d건 총 투자 금액 %d원\n",
			setting.Company, setting.Username, count, setting.Amount*count)
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

func newRunner(setting *autop2p.Setting) autop2p.Runner {
	switch setting.Company {
	case autop2p.Honestfund:
		return honestfund.NewRunner(setting, HonestfundService)
	case autop2p.Peoplefund:
		return peoplefund.NewRunner(setting, PeoplefundService)
	default:
		panic("unsupported company type")
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

func main() {
	lambda.Start(Run)
}


