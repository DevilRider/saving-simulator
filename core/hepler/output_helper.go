package hepler

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	beans2 "saving-simulator/core/beans"
	"saving-simulator/core/constant"
	"saving-simulator/utils"
	"sort"
)

const outputPath = "app-output/"
const splitString = "-----------------------------------------------"

type OutputHelper struct {
	Summary      beans2.Summary
	Acc          *beans2.Account
	Transactions []beans2.Transaction
	Amount       decimal.Decimal
	Principal    decimal.Decimal
	Interest     decimal.Decimal
}

func NewOutputor(summary beans2.Summary, acc *beans2.Account, transactions []beans2.Transaction) *OutputHelper {
	return &OutputHelper{
		Summary:      summary,
		Acc:          acc,
		Transactions: sortTransactions(transactions),
		Amount:       decimal.Zero,
		Principal:    decimal.Zero,
		Interest:     decimal.Zero,
	}
}

func (o *OutputHelper) Output() {
	filename := fmt.Sprintf("%s%d Years %s[%d] %s.txt", outputPath, o.Summary.Years, o.Summary.Frequency, o.Summary.Step.IntPart(), o.Acc.GetTotal().Round(2).String())
	utils.WriteCsv(filename, o.rows())
	logrus.Infof("filename: %s", filename)
}

func (o *OutputHelper) rows() [][]string {
	var rows [][]string
	rows = append(rows, []string{splitString})
	for _, val := range o.headers() {
		rows = append(rows, val)
	}
	rows = append(rows, []string{splitString})
	rows = append(rows, []string{"Date", "Saving Amount", "Principal", "Interest"})
	for _, val := range o.results() {
		rows = append(rows, val)
	}
	return rows
}

func (o *OutputHelper) headers() [][]string {
	var rows [][]string
	rows = append(rows, []string{fmt.Sprintf("Save %d/%s for %d Years.", o.Summary.Step.IntPart(), o.Summary.Frequency, o.Summary.Years)})
	rows = append(rows, []string{fmt.Sprintf("Total Amount: %s", o.Acc.GetTotal().Round(2).String())})
	rows = append(rows, []string{fmt.Sprintf("Principal: %s", o.Acc.Principal.Round(2).String())})
	rows = append(rows, []string{fmt.Sprintf("Earnings: %s", o.Acc.GetEarnings().Round(2).String())})
	return rows
}

func (o *OutputHelper) results() [][]string {
	var rows [][]string
	t := o.Summary.Start
	for _, trx := range o.Transactions {
		if !t.Equal(trx.Date) {
			rows = append(rows, []string{
				trx.Date.Format("2006-01-02"), o.Amount.Round(2).String(), o.Principal.Round(2).String(), o.Interest.Round(2).String(),
			})
			t = trx.Date
		}
		if trx.Type == constant.TransactionTypeSaving {
			o.Amount = o.Amount.Add(trx.Amount)
			o.Principal = o.Principal.Add(trx.Amount)
		}

		if trx.Type == constant.TransactionTypeInterest {
			o.Amount = o.Amount.Add(trx.Amount)
			o.Interest = o.Interest.Add(trx.Amount)
		}
	}
	return rows
}

func sortTransactions(transactions []beans2.Transaction) []beans2.Transaction {
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Date.Before(transactions[j].Date)
	})
	return transactions
}
