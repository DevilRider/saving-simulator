package beans

import (
	"github.com/shopspring/decimal"
	"saving-simulator/utils"
)

type Product struct {
	Type               string          `json:"type" validate:"required"`
	Name               string          `json:"name" validate:"required"`
	Period             int64           `json:"period" validate:"required"`
	SubscriptionDays   int64           `json:"subscriptionDays" validate:"required"`
	SettlementDays     int64           `json:"settlementDays" validate:"required"`
	AnnualInterestRate decimal.Decimal `json:"annualInterestRate" validate:"required"`
	Limit              int64           `json:"limit"`
	Balance            int64           // 用于统计 limit 剩余情况， 如果是-1 则表示无限制
	Previous           *Product
}

// SortKey ( 0.7 * a - 0.3 * (b / 14)) * 10000
func (p *Product) SortKey() int64 {
	return decimal.NewFromFloat(0.7).Mul(p.AnnualInterestRate).
		Sub(decimal.NewFromFloat(0.3).Mul(decimal.NewFromInt(p.blockDays()).Div(decimal.NewFromInt(14)))).
		Mul(decimal.NewFromInt(10000)).IntPart()
}

func (p *Product) blockDays() int64 {
	return p.SubscriptionDays + p.SettlementDays
}

func (p *Product) Days() int64 {
	return p.SubscriptionDays + p.Period + p.SettlementDays
}

func (p *Product) IsFixedTerm() bool {
	return p.Type != "current"
}

func (p *Product) ToJson() string {
	temp, _ := utils.JSON.Marshal(p)
	return string(temp)
}
