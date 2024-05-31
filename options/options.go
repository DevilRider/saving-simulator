package options

import (
	"github.com/shopspring/decimal"
	"time"
)

type Options struct {
	Start         time.Time
	End           time.Time
	Amount        decimal.Decimal
	Years         int
	Frequency     string // Weekly, Monthly
	Step          decimal.Decimal
	ComplexSaving *ComplexSavingCfg
}

type ComplexSavingCfg struct {
	SavingThreshold decimal.Decimal // 存款阈值
	MaxSavingDays   int64           // 最大存款天数
}

type Option func(*Options)

func Start(start time.Time) Option {
	return func(c *Options) {
		c.Start = start
	}
}

func Frequency(Type string) Option {
	return func(c *Options) {
		c.Frequency = Type
	}
}

func Step(step int64) Option {
	return func(c *Options) {
		c.Step = decimal.NewFromInt(step)
	}
}

func End(end time.Time) Option {
	return func(c *Options) {
		c.End = end
	}
}

func Amount(amount decimal.Decimal) Option {
	return func(c *Options) {
		c.Amount = amount
	}
}

func ComplexSaving(threshold int64, max int64) Option {
	return func(c *Options) {
		c.ComplexSaving = &ComplexSavingCfg{
			SavingThreshold: decimal.NewFromInt(threshold),
			MaxSavingDays:   max,
		}
	}
}

func Years(years int) Option {
	return func(c *Options) {
		c.Years = years
	}
}
