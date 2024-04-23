package types

import "time"

type Cashout struct {
	Amount      int       `json:"amount" db:"amount" validate:"required"`
	Date        time.Time `json:"date" db:"date" validate:"required"`
	Id          int       `json:"id" db:"id"`
	PortfolioId int       `json:"portfolioId" db:"portfolio_id" validate:"required"`
}
