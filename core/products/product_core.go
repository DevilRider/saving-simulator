package products

import (
	"saving-simulator/core/beans"
)

func All() []beans.Product {
	return FixedTerms.Products
}

func Match(period int64) *beans.Product {
	if period == -1 {
		return CurrentPlus
	}
	for _, product := range FixedTerms.Products {
		if product.Period <= period {
			return &product
		}
	}
	return nil
}
