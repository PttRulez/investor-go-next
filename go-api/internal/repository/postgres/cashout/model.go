package cashout 

import "time"

type Cashout struct {
	Id          int       `db:"id"`
	Amount      int       `db:"amount"`
	Date        time.Time `db:"date"`
	PortfolioId int       `db:"portfolio_id"`
}
