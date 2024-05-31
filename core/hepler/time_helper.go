package hepler

import (
	"saving-simulator/core/constant"
	"time"
)

type TimeHelper struct {
	Type string
	T    time.Time
	End  time.Time
}

func NewDailyHelper(start time.Time, end time.Time) *TimeHelper {
	return &TimeHelper{
		Type: constant.ContinuousSavingInterestAccrualType,
		T:    start,
		End:  end,
	}
}

func (impl *TimeHelper) IsEnd() bool {
	return impl.T.After(impl.End)
}

func (impl *TimeHelper) Increase() *TimeHelper {
	switch impl.Type {
	case "Weekly":
		impl.T = impl.T.AddDate(0, 0, 7)
		break
	case "Monthly":
		impl.T = impl.T.AddDate(0, 1, 0)
		break
	default: // Daily
		impl.T = impl.T.AddDate(0, 0, 1)
	}
	return impl
}

func (impl *TimeHelper) IncreaseInterestAccrualPeriod() *TimeHelper {
	switch impl.Type {
	case "Weekly":
		impl.T = impl.T.AddDate(0, 0, 7)
		break
	case "Monthly":
		impl.T = impl.T.AddDate(0, 1, 0)
		break
	default: // Daily
		impl.T = impl.T.AddDate(0, 0, 1)
	}

	return impl
}
