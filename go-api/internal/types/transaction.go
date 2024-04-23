package types

import "time"

type Transaction struct {
	Amount      int       `json:"amount" db:"amount" validate:"required"`
	Date        time.Time `json:"date" db:"date" validate:"required"`
	Id          int       `json:"id" db:"id"`
	PortfolioId int       `json:"portfolioId" db:"portfolio_id" validate:"required"`
}
type Cashout struct {
	Transaction
}
type Deposit struct {
	Transaction
}
