package products

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"saving-simulator/core/beans"
	"sort"
)

const name = "products.json"

var (
	CurrentPlus *beans.Product
	FixedTerms  *ProductPool
)

type ProductPool struct {
	Products []beans.Product
}

func LoadProducts() {
	// init current product
	CurrentPlus = &beans.Product{
		Name:               "CurrentPlus",
		Type:               "current",
		Period:             1,
		SubscriptionDays:   0,
		SettlementDays:     0,
		AnnualInterestRate: decimal.NewFromFloat(0.0216),
		Limit:              -1,
	}

	// load fixed products
	data, err := os.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("failed to read file: %v", err))
	}
	var ps []beans.Product
	err = json.Unmarshal(data, &ps)
	if err != nil {
		panic(fmt.Sprintf("failed to read file: %v", err))
	}

	FixedTerms = &ProductPool{
		Products: ps,
	}

	sort.SliceStable(FixedTerms.Products, func(i, j int) bool {
		return FixedTerms.Products[i].SortKey() >= FixedTerms.Products[j].SortKey()
	})

	for idx, _ := range FixedTerms.Products {
		p := FixedTerms.Products[idx]
		if idx == 0 {
			p.Previous = nil
		} else {
			p.Previous = &FixedTerms.Products[idx-1]
		}
		FixedTerms.Products[idx] = p
	}
}
