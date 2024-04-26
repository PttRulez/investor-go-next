package model

import "time"

type Cashout struct {
	Amount      int
	Date        time.Time
	Id          int
	PortfolioId int
}
