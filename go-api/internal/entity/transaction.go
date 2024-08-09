package entity

import "time"

type Cashout struct {
	Amount      int
	Date        time.Time
	Id          int
	PortfolioId int
	UserId      int
}
