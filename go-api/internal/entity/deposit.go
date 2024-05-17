package entity

import "time"

type Deposit struct {
	Amount      int
	Date        time.Time
	Id          int
	PortfolioId int
	UserId      int
}
