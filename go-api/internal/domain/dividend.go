package domain

import "time"

type Dividend struct {
	Date            time.Time
	Exchange        Exchange
	ID              int
	PaymentPeriod   string
	PaymentPerShare float64
	PortfolioID     int
	Ticker          string
	SharesCount     int
}
