package beans

import (
	"github.com/shopspring/decimal"
	"saving-simulator/core/constant"
	"saving-simulator/utils"
	"time"
)

// Account 账户
type Account struct {
	Principal   decimal.Decimal // 本金 (账户本金)
	Balance     decimal.Decimal // 余额, 利息太少忽略不计
	CurrentPlus decimal.Decimal // 活期+资产（类似余额宝，年化2%，每日计息，随时可赎）
	Investment  decimal.Decimal // 投资资产 (类似定期，年化5%，到期赎回)
	Locked      []LockedAmount  // 冻结金额
}

func EmptyAccount() *Account {
	return &Account{
		Principal:   decimal.Zero,
		Balance:     decimal.Zero,
		CurrentPlus: decimal.Zero,
		Investment:  decimal.Zero,
	}
}

func (acc *Account) ToJson() string {
	temp, _ := utils.JSON.Marshal(acc)
	return string(temp)
}

func (acc *Account) GetTotal() decimal.Decimal {
	return acc.Balance.Add(acc.CurrentPlus).Add(acc.Investment)
}

func (acc *Account) GetTotalCurrent() decimal.Decimal {
	return acc.Balance.Add(acc.CurrentPlus)
}

func (acc *Account) GetEarnings() decimal.Decimal {
	return acc.GetTotal().Sub(acc.Principal)
}

// Save add amount to balance
func (acc *Account) Save(amount decimal.Decimal) *Account {
	acc.Balance = acc.Balance.Add(amount)
	acc.Principal = acc.Principal.Add(amount)
	return acc
}

func (acc *Account) Lock(transaction Transaction) {
	if transaction.Type == constant.TransactionTypePurchase {
		acc.Balance = acc.Balance.Sub(transaction.Amount)
	} else if transaction.Type == constant.TransactionTypeRedeem {
		if transaction.Product.IsFixedTerm() {
			acc.Investment = acc.Investment.Sub(transaction.Amount)
		} else {
			acc.CurrentPlus = acc.CurrentPlus.Sub(transaction.Amount)
		}
	}

	acc.Locked = append(acc.Locked, LockedAmount{
		Type:       transaction.Type,
		Amount:     transaction.Amount,
		SettleDate: transaction.SettleDate,
		Ref:        transaction.Id,
	})
}

func (acc *Account) Release(t time.Time, transaction Transaction) {
	for i, locked := range acc.Locked {
		if locked.Ref == transaction.Id && !t.Before(locked.SettleDate) {
			if locked.Type == constant.AccLockTypePurchase {
				if transaction.Product.IsFixedTerm() {
					acc.Investment = acc.Investment.Add(locked.Amount)
				} else {
					acc.CurrentPlus = acc.CurrentPlus.Add(locked.Amount)
				}

				acc.Locked = append(acc.Locked[:i], acc.Locked[i+1:]...)
			}
			if locked.Type == constant.AccLockTypeRedeem {
				acc.Balance = acc.Balance.Add(locked.Amount)
				acc.Locked = append(acc.Locked[:i], acc.Locked[i+1:]...)
			}
		}
	}
}
