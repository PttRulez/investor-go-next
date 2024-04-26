package model

import "time"

type Deposit struct {
	Amount      int
	Date        time.Time
	Id          int
	PortfolioId int
}
