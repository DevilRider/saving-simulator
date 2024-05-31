package beans

import (
	"saving-simulator/utils"
	"time"

	"github.com/shopspring/decimal"
)

type FinancePlan struct {
	Product            *Product
	Start              time.Time
	ValueDate          time.Time
	End                time.Time
	SettleDate         time.Time
	Days               int64
	AnnualInterestRate decimal.Decimal
}

func (p *FinancePlan) ToJson() string {
	temp, _ := utils.JSON.Marshal(p)
	return string(temp)
}

func (p *FinancePlan) Interest(amount decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(p.Days).Mul(p.AnnualInterestRate).Div(decimal.NewFromInt(365)).Mul(amount)
}
