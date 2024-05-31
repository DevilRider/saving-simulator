package constant

const (
	HoursOfDay                          = 24
	HoursOfWeek                         = 168
	DaysOfYear                          = 365
	ContinuousSavingInterestAccrualType = "Daily"
	SavingFrequencyOnce                 = "Once"
	SavingFrequencyWeekly               = "Weekly"
	SavingFrequencyMonthly              = "Monthly"
	TransactionStatusPending            = "Pending"
	TransactionStatusSettled            = "Settled"
	TransactionTypePurchase             = "Purchase"
	TransactionTypeRedeem               = "Redeem"
	TransactionTypeSaving               = "Saving"
	TransactionTypeInterest             = "Interest"
	AccLockTypeRedeem                   = TransactionTypeRedeem
	AccLockTypePurchase                 = TransactionTypePurchase
)
