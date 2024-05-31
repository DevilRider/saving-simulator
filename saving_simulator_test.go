package saving_simulator

import (
	"github.com/sirupsen/logrus"
	"saving-simulator/core/constant"
	"saving-simulator/core/hepler"
	"saving-simulator/core/products"
	"saving-simulator/options"
	"saving-simulator/utils"
	"testing"
)

func Test_Simulate_Simple_Weekly_1Y(t *testing.T) {
	simulator := NewSimulator(
		options.Start(utils.ConvertStringToTime("2024-04-16", "2006-01-02")),
		options.Years(1),
		options.Frequency(constant.SavingFrequencyWeekly),
		options.Step(500),
	)
	_, acc, transactions := simulator.Simulate()
	logrus.Infof("Total : %s,", acc.GetTotal().Round(2).String())
	logrus.Infof("Principal : %s,", acc.Principal.Round(2).String())
	logrus.Infof("Earnings : %s,", acc.GetEarnings().Round(2).String())
	logrus.Infof("Balance : %s,", acc.Balance.Round(2).String())
	logrus.Infof("CurrentPlus : %s,", acc.CurrentPlus.Round(2).String())
	logrus.Infof("transactions: %d", len(transactions))
	//now := time.Now()
	for _, trx := range transactions {
		if trx.Type == constant.TransactionTypeSaving {
			logrus.Infof("trx: %s", trx.ToJson())
		}
	}
}

func Test_Simulate_Simple_Monthly_1Y(t *testing.T) {
	simulator := NewSimulator(
		options.Start(utils.ConvertStringToTime("2023-04-16", "2006-01-02")),
		options.Years(1),
		options.Frequency(constant.SavingFrequencyMonthly),
		options.Step(2210),
	)
	_, acc, transactions := simulator.Simulate()
	logrus.Infof("Total : %s,", acc.GetTotal().Round(2).String())
	logrus.Infof("Principal : %s,", acc.Principal.Round(2).String())
	logrus.Infof("Earnings : %s,", acc.GetEarnings().Round(2).String())
	logrus.Infof("Balance : %s,", acc.Balance.Round(2).String())
	logrus.Infof("CurrentPlus : %s,", acc.CurrentPlus.Round(2).String())
	logrus.Infof("transactions: %d", len(transactions))
	for _, trx := range transactions {
		if trx.Type == constant.TransactionTypeSaving {
			logrus.Infof("trx: %s", trx.ToJson())
		}
	}
}

func Test_Simulate_Complex_Weekly_1Y(t *testing.T) {
	simulator := NewSimulator(
		options.Start(utils.ConvertStringToTime("2023-04-16", "2006-01-02")),
		options.Years(1),
		options.Frequency(constant.SavingFrequencyWeekly),
		options.Step(500),
		options.ComplexSaving(10000, 366),
	)
	_, acc, transactions := simulator.Simulate()

	logrus.Infof("Total : %s,", acc.GetTotal().Round(2).String())
	logrus.Infof("Principal : %s,", acc.Principal.Round(2).String())
	logrus.Infof("Earnings : %s,", acc.GetEarnings().Round(2).String())
	logrus.Infof("Balance : %s,", acc.Balance.Round(2).String())
	logrus.Infof("CurrentPlus : %s,", acc.CurrentPlus.Round(2).String())
	logrus.Infof("transactions: %d", len(transactions))
}

func Test_Simulate_Complex_Monthly_1Y(t *testing.T) {
	simulator := NewSimulator(
		options.Start(utils.ConvertStringToTime("2023-04-16", "2006-01-02")),
		options.Years(1),
		options.Frequency(constant.SavingFrequencyMonthly),
		options.Step(2210),
		options.ComplexSaving(10000, 366),
	)
	summary, acc, transactions := simulator.Simulate()

	logrus.Infof("Total : %s,", acc.GetTotal().Round(2).String())
	logrus.Infof("Principal : %s,", acc.Principal.Round(2).String())
	logrus.Infof("Earnings : %s,", acc.GetEarnings().Round(2).String())
	logrus.Infof("Balance : %s,", acc.Balance.Round(2).String())
	logrus.Infof("CurrentPlus : %s,", acc.CurrentPlus.Round(2).String())
	logrus.Infof("transactions: %d", len(transactions))
	out := hepler.NewOutputor(summary, acc, transactions)
	out.Output()

}

func Test_Simulate_Complex_Monthly_2Y(t *testing.T) {
	simulator := NewSimulator(
		options.Start(utils.ConvertStringToTime("2023-04-16", "2006-01-02")),
		options.Years(2),
		options.Frequency(constant.SavingFrequencyMonthly),
		options.Step(2210),
		options.ComplexSaving(10000, 366),
	)
	summary, acc, transactions := simulator.Simulate()

	logrus.Infof("Total : %s,", acc.GetTotal().Round(2).String())
	logrus.Infof("Principal : %s,", acc.Principal.Round(2).String())
	logrus.Infof("Earnings : %s,", acc.GetEarnings().Round(2).String())
	logrus.Infof("Balance : %s,", acc.Balance.Round(2).String())
	logrus.Infof("CurrentPlus : %s,", acc.CurrentPlus.Round(2).String())
	logrus.Infof("transactions: %d", len(transactions))
	out := hepler.NewOutputor(summary, acc, transactions)
	out.Output()

}

func Test_Simulate_Complex_Monthly_24Y(t *testing.T) {
	simulator := NewSimulator(
		options.Start(utils.ConvertStringToTime("2023-04-16", "2006-01-02")),
		options.Years(24),
		options.Frequency(constant.SavingFrequencyMonthly),
		options.Step(2210),
		options.ComplexSaving(10000, 366),
	)
	summary, acc, transactions := simulator.Simulate()

	logrus.Infof("Total : %s,", acc.GetTotal().Round(2).String())
	logrus.Infof("Principal : %s,", acc.Principal.Round(2).String())
	logrus.Infof("Earnings : %s,", acc.GetEarnings().Round(2).String())
	logrus.Infof("Balance : %s,", acc.Balance.Round(2).String())
	logrus.Infof("CurrentPlus : %s,", acc.CurrentPlus.Round(2).String())
	logrus.Infof("transactions: %d", len(transactions))
	out := hepler.NewOutputor(summary, acc, transactions)
	out.Output()

}

func Test_LoadProducts(t *testing.T) {

	products.LoadProducts()
	ps := products.All()
	logrus.Infof("Products: %d", len(ps))
	for _, p := range ps {
		logrus.Infof("Product: %s", p.ToJson())
	}
}
