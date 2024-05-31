package beans

import (
	"github.com/shopspring/decimal"
	"saving-simulator/utils"
	"time"
)

type Summary struct {
	Start     time.Time
	Years     int
	Frequency string // Weekly, Monthly
	Step      decimal.Decimal
}

func (d *Summary) ToJson() string {
	temp, _ := utils.JSON.Marshal(d)
	return string(temp)
}
