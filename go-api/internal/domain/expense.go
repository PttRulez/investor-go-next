package domain

import "time"

type Expense struct {
	Amount      float64
	Date        time.Time
	Description string
	ID          int
	PortfolioID int
}
