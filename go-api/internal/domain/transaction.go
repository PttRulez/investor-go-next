package domain

import "time"

type Transaction struct {
	Amount      float64
	Date        time.Time
	ID          int
	PortfolioID int
	Type        TransactionType
	UserID      int
}

type TransactionType string

const (
	TTDeposit TransactionType = "DEPOSIT"
	TTCashout TransactionType = "CASHOUT"
)

func (e TransactionType) Validate() bool {
	switch e {
	case TTDeposit:
	case TTCashout:
	default:
		return false
	}
	return true
}
