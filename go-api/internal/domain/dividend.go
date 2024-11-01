package domain

import "time"

type Dividend struct {
	Date            time.Time
	Exchange        Exchange
	ID              int
	PaymentPeriod   string
	PaymentPerShare float64
	PortfolioID     int
	SharesCount     int
	ShortName       string
	TaxPaid         float64
	Ticker          string
	TotalPayment    float64
	UserID          int
}
