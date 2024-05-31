package beans

import (
	"github.com/shopspring/decimal"
	"saving-simulator/utils"
	"time"
)

type LockedAmount struct {
	Type       string // 购入 Purchase, 赎回 Redeem
	Amount     decimal.Decimal
	SettleDate time.Time
	Ref        string // 用于关联 Redeem
}

func (d *LockedAmount) ToJson() string {
	temp, _ := utils.JSON.Marshal(d)
	return string(temp)
}
