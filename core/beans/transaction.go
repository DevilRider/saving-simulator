package beans

import (
	"github.com/shopspring/decimal"
	"saving-simulator/core/constant"
	"saving-simulator/utils"
	"time"
)

type Transaction struct {
	Id         string
	Product    *Product
	Date       time.Time // 交易日
	SettleDate time.Time // 结算日
	Type       string    // 交易类型: 买入 Purchase，赎回 Redeem，分红 Interest
	Status     string    // 交易状态: 待结算 Pending，已结算 Settled
	Amount     decimal.Decimal
	Ref        string // 用于关联买入和赎回
}

func NewPurchaseTransaction(t time.Time, amount decimal.Decimal, p *Product) Transaction {
	return Transaction{
		Id:         utils.NewObjectID().Hex(),
		Date:       t,
		Product:    p,
		SettleDate: t.AddDate(0, 0, int(p.SubscriptionDays)),
		Status:     constant.TransactionStatusPending,
		Type:       constant.TransactionTypePurchase,
		Amount:     amount,
	}
}

func NewRedeemTransaction(t time.Time, amount decimal.Decimal, p *Product) Transaction {
	return Transaction{
		Id:         utils.NewObjectID().Hex(),
		Date:       t,
		Product:    p,
		SettleDate: t.AddDate(0, 0, int(p.SettlementDays)),
		Status:     constant.TransactionStatusPending,
		Type:       constant.TransactionTypeRedeem,
		Amount:     amount,
	}
}

func NewInterestTransaction(t time.Time, interest decimal.Decimal, p *Product) Transaction {
	return Transaction{
		Id:         utils.NewObjectID().Hex(),
		Product:    p,
		Date:       t,
		SettleDate: t.AddDate(0, 0, int(p.SettlementDays)),
		Status:     constant.TransactionStatusPending,
		Type:       constant.TransactionTypeInterest,
		Amount:     interest,
	}
}

func (t *Transaction) ToJson() string {
	temp, _ := utils.JSON.Marshal(t)
	return string(temp)
}
