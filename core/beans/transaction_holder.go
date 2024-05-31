package beans

import (
	"saving-simulator/core/constant"
	"saving-simulator/utils"
)

type TransactionHolder struct {
	Pending []Transaction
	Settled []Transaction
}

// OptimizePendingTransactions remove settled transactions from pending
func (h *TransactionHolder) OptimizePendingTransactions() {
	if len(h.Pending) == 0 {
		return
	}
	tmp := h.Pending
	h.Pending = nil
	for _, trx := range tmp {
		if trx.Status == constant.TransactionStatusPending {
			h.AppendPending(trx)
		}
	}
}

func (h *TransactionHolder) HasPendingTransactions() bool {
	return len(h.Pending) > 0
}

func (h *TransactionHolder) AppendPending(t Transaction) {
	h.Pending = append(h.Pending, t)
}

func (h *TransactionHolder) AppendSettled(t Transaction) {
	h.Settled = append(h.Settled, t)
}

func (h *TransactionHolder) Settle(idx int, t Transaction) {
	t.Status = constant.TransactionStatusSettled
	h.Pending[idx] = t
	h.AppendSettled(t)
}

func (h *TransactionHolder) ToJson() string {
	temp, _ := utils.JSON.Marshal(h)
	return string(temp)
}
