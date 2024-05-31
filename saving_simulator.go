package saving_simulator

import (
	"saving-simulator/core"
	"saving-simulator/core/beans"
	"saving-simulator/core/constant"
	"saving-simulator/core/hepler"
	"saving-simulator/core/products"
	"saving-simulator/options"
	"saving-simulator/utils"
	"time"

	"github.com/shopspring/decimal"
)

type SavingSimulator struct {
	Start         time.Time
	Years         int
	Frequency     string                    // Once, Weekly, Monthly
	Step          decimal.Decimal           // 每次存款金额
	ComplexSaving *options.ComplexSavingCfg // 就整储蓄

	timer        *hepler.TimeHelper
	end          time.Time
	account      *beans.Account
	transactions beans.TransactionHolder
}

func NewSimulator(ops ...options.Option) *SavingSimulator {
	products.LoadProducts()
	cfg := &options.Options{}
	for _, opt := range ops {
		opt(cfg)
	}
	if cfg.Start.IsZero() {
		cfg.Start = time.Now()
	}
	simulator := &SavingSimulator{
		Start:     cfg.Start,
		Years:     cfg.Years,
		Frequency: cfg.Frequency,
		Step:      cfg.Step,
		end:       cfg.Start.AddDate(cfg.Years, 0, 0),
		account:   beans.EmptyAccount(),
	}
	simulator.timer = hepler.NewDailyHelper(simulator.Start, simulator.end)
	if cfg.ComplexSaving != nil {
		simulator.ComplexSaving = cfg.ComplexSaving
	}
	return simulator
}

// 认购期 subscription period
// 清算期 settlement period
// 买入 purchase
// 赎回 redemption

// Simulate
// 1.根据时间判断是否需要存钱
// 1.1. 每周存钱/每月存钱
// 2.根据余额(current + current plus)判断是否需买入理财产品
// 如果 balance 直接就满足阈值 这直接进行长期产品买入
// 如果 需要current plus中的金额，则先进行赎回
// 如果 小于阈值时，使用全部balance买入current plus
// 2.1.赎回
// 赎回交易流水 --> 将金额放入锁定金额
// 3.计算利息
// 计算活期理财产品利息
// 长期理财产品无需利息计算，买入时 生成出买入流水、计息流水（按日生成）、赎回流水
// 4.结算
// 结算交易流水
// 基于交易流水进行 余额 活期 理财产品 锁定金额的转换
func (s *SavingSimulator) Simulate() (beans.Summary, *beans.Account, []beans.Transaction) {
	for {
		if s.timer.IsEnd() {
			s.handlePendingTransactionsIfNecessary()
			return beans.Summary{
				Start:     s.Start,
				Years:     s.Years,
				Frequency: s.Frequency,
				Step:      s.Step,
			}, s.account, s.transactions.Settled
		}

		s.save(). // 资金流入
				purchase().       // 资金流出，可能触发赎回
				redeem().         // 如果有自动赎回，则先锁定金额，已实现赎回日停息
				interest().       // 流入，但是需要结算触发
				settle(s.timer.T) // 结算完成后， 涉及到重新买入

		// 结算完成后，利息会进入余额，所以还需要走一次买入结算流程（因为current plus，当日即可买入）。
		if !s.account.Balance.Equal(decimal.Zero) {
			s.purchase().settle(s.timer.T)
		}

		s.timer.Increase()
	}
}

func (s *SavingSimulator) handlePendingTransactionsIfNecessary() {
	hasPendingTrx := s.transactions.HasPendingTransactions()
	t := s.end.AddDate(0, 0, 1)
	for {
		if !hasPendingTrx {
			break
		}
		s.interest().settle(t)
		t = t.AddDate(0, 0, 1)
		hasPendingTrx = s.transactions.HasPendingTransactions()
	}
}

// save 1.根据时间判断是否需要存钱 每周存钱/每月存钱
func (s *SavingSimulator) save() *SavingSimulator {
	t := s.timer.T
	if !s.Start.Equal(t) {
		if s.end.Equal(t) { // 最后一天不买入
			return s
		}

		if s.Frequency == constant.SavingFrequencyOnce {
			return s
		}

		if s.Frequency == constant.SavingFrequencyWeekly && int(t.Sub(s.Start).Hours())%constant.HoursOfWeek != 0 {
			return s
		}

		if s.Frequency == constant.SavingFrequencyMonthly {
			tmp := s.Start.AddDate(0, 1, 0)
			for {
				if tmp.After(t) {
					return s
				}
				if tmp.Equal(t) {
					break
				}
				tmp = tmp.AddDate(0, 1, 0)
			}
		}
	}

	s.account.Save(s.Step)
	s.transactions.AppendSettled(beans.Transaction{
		Id:         utils.NewObjectID().Hex(),
		Date:       t,
		SettleDate: t,
		Status:     constant.TransactionStatusSettled,
		Type:       constant.TransactionTypeSaving,
		Amount:     s.Step,
	})
	return s
}

// purchase 2.根据余额(balance + current plus)判断是否需买入理财产品
// 循环执行买入操作， 退出条件： 余额为0
// 如果 Balance 直接就满足阈值 这直接进行长期产品买入
// 如果 需要Current Plus中的金额，则先进行赎回
// 如果 小于阈值时，使用全部balance买入current plus活期理财
func (s *SavingSimulator) purchase() *SavingSimulator {
	if s.isComplexSavingEnabled() {
		planner := core.NewFinancePlanner(
			options.Start(s.timer.T),
			options.End(s.end),
			options.Amount(s.ComplexSaving.SavingThreshold),
			options.ComplexSaving(s.ComplexSaving.SavingThreshold.IntPart(), s.ComplexSaving.MaxSavingDays),
		)
		plan := planner.Generate()

		// 是否可以直接购买长期理财: 金额满足，且有可用长期产品
		if s.account.Balance.GreaterThanOrEqual(s.ComplexSaving.SavingThreshold) && plan != nil && plan.Product != nil {
			// 长期理财产品无需利息计算，买入时生成 买入流水，计息流水，赎回流水
			purchaseTrx := beans.NewPurchaseTransaction(plan.Start, s.ComplexSaving.SavingThreshold, plan.Product)
			s.transactions.AppendPending(purchaseTrx)
			s.account.Lock(purchaseTrx)

			s.transactions.AppendPending(beans.NewInterestTransaction(plan.End, plan.Interest(s.ComplexSaving.SavingThreshold), plan.Product))
			s.transactions.AppendPending(beans.NewRedeemTransaction(plan.End, s.ComplexSaving.SavingThreshold, plan.Product))
		}

		// 是否需要赎回current plus: 金额满足，且有可用长期产品
		if s.account.GetTotalCurrent().GreaterThanOrEqual(s.ComplexSaving.SavingThreshold) && plan != nil && plan.Product != nil {
			// 2. 先赎回current plus，赎回后活期（二次循环后走1.进行长期产品购买）
			amountToRedeem := s.ComplexSaving.SavingThreshold
			if s.account.CurrentPlus.LessThan(s.ComplexSaving.SavingThreshold) {
				amountToRedeem = s.ComplexSaving.SavingThreshold.Sub(s.account.Balance)
			}
			transaction := beans.NewRedeemTransaction(s.timer.T, amountToRedeem, products.CurrentPlus)
			s.transactions.AppendPending(transaction)
			return s
		}

	}

	if !s.account.Balance.Equal(decimal.Zero) {
		// 3. 购买活期理财产品
		transaction := beans.NewPurchaseTransaction(s.timer.T, s.account.Balance, products.CurrentPlus)
		s.transactions.AppendPending(transaction)

		s.account.Lock(transaction)
	}

	return s
}

// interest 3.计算利息
// 计算活期理财产品利息
// 长期理财产品无需利息计算，买入时 生成出 买入流水，计息流水，赎回流水
func (s *SavingSimulator) interest() *SavingSimulator {
	if s.account.CurrentPlus.Equal(decimal.Zero) {
		return s
	}
	interest := products.CurrentPlus.AnnualInterestRate.Div(decimal.NewFromInt(constant.DaysOfYear)).Mul(s.account.CurrentPlus)
	s.transactions.AppendPending(beans.NewInterestTransaction(s.timer.T, interest, products.CurrentPlus))
	return s
}

// Redeem 执行赎回操作， 根据transaction lock 相应金额
// 赎回操作需要配合处理 金额锁定
// t==trx.time --> 将金额放入锁定金额
func (s *SavingSimulator) redeem() *SavingSimulator {
	for _, trx := range s.transactions.Pending {
		// t==trx.time --> 将金额放入锁定金额
		if trx.Status == constant.TransactionStatusPending && trx.Type == constant.TransactionTypeRedeem && trx.Date.Equal(s.timer.T) {
			s.account.Lock(trx)
		}
	}
	return s
}

// 5.结算
// 结算交易流水
// 基于交易流水进行 余额 余额+ 理财产品 锁定金额的转换
// 赎回操作需要配合处理
// t>=trx.settleTime --> 将金额从锁定金额移动到 current
func (s *SavingSimulator) settle(t time.Time) *SavingSimulator {
	// 针对已有 pending transaction 进行结算
	day := t.Weekday()
	for idx, trx := range s.transactions.Pending {
		if trx.Status == constant.TransactionStatusSettled || trx.SettleDate.After(t) || // 已结算 或 结算日期大于结算日期的不做处理
			day == time.Sunday || day == time.Saturday || day == time.Friday { // 周五 周六 周日 不结算
			continue
		}
		// 结算 purchase
		switch trx.Type {
		case constant.TransactionTypeSaving:
			continue
		case constant.TransactionTypePurchase:
			s.account.Release(t, trx)
			s.transactions.Settle(idx, trx)
		case constant.TransactionTypeRedeem:
			s.account.Release(t, trx)
			s.transactions.Settle(idx, trx)
		case constant.TransactionTypeInterest:
			s.account.Balance = s.account.Balance.Add(trx.Amount)
			s.transactions.Settle(idx, trx)
		default:
			panic("transaction type not support")
		}
	}

	// remove settled trx from pending list
	s.transactions.OptimizePendingTransactions()
	return s
}

func (s *SavingSimulator) isComplexSavingEnabled() bool {
	return s.ComplexSaving != nil && !s.account.Balance.Equal(decimal.Zero)
}
