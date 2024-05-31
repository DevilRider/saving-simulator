package core

import (
	"github.com/shopspring/decimal"
	"saving-simulator/core/beans"
	"saving-simulator/core/constant"
	"saving-simulator/core/products"
	"saving-simulator/options"
	"time"
)

type Planner struct {
	Start  time.Time
	End    time.Time
	Amount decimal.Decimal
	Period int64
}

func NewFinancePlanner(opts ...options.Option) *Planner {
	cfg := &options.Options{}
	for _, opt := range opts {
		opt(cfg)
	}

	end := cfg.End
	if end.IsZero() {
		end = cfg.Start.AddDate(cfg.Years, 0, 0)
	}

	var period int64

	if cfg.ComplexSaving != nil && cfg.ComplexSaving.MaxSavingDays != 0 {
		period = cfg.ComplexSaving.MaxSavingDays
	} else {
		period = -1 // default -1, means not set, no fixed period
	}

	return &Planner{
		Start:  cfg.Start,
		End:    end,
		Period: period,
	}
}

func (p *Planner) Generate() *beans.FinancePlan {
	product := products.Match(p.Period)
	return recursive(p.Start, p.End, product)
}

func recursive(start time.Time, end time.Time, product *beans.Product) *beans.FinancePlan {
	plan := planing(start, end, product)
	if plan != nil {
		return plan
	}
	previous := product.Previous
	if previous != nil {
		return recursive(start, end, previous)
	}
	return nil
}

func planing(start time.Time, end time.Time, product *beans.Product) *beans.FinancePlan {
	days := end.Sub(start).Hours() / constant.HoursOfDay
	counts := decimal.NewFromFloat(days).Div(intToDecimal(product.Period))
	if counts.IntPart() == 0 {
		return nil
	}
	return &beans.FinancePlan{
		Product:            product,
		Start:              start,
		ValueDate:          start.AddDate(0, 0, int(product.SubscriptionDays)),
		End:                start.AddDate(0, 0, int(product.Period)-int(product.SettlementDays)),
		SettleDate:         start.AddDate(0, 0, int(product.Period)),
		Days:               product.Days(),
		AnnualInterestRate: product.AnnualInterestRate,
	}
}

func intToDecimal(i int64) decimal.Decimal {
	return decimal.NewFromInt(i)
}
