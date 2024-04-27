package pgdeposit

import "time"

type Deposit struct {
	Id          int       `db:"id"`
	Amount      int       `db:"amount"`
	Date        time.Time `db:"date"`
	PortfolioId int       `db:"portfolio_id"`
}
