package domain

import "time"

type Coupon struct {
	BondsCount    int
	CouponAmount  float64
	Date          time.Time
	Exchange      Exchange
	ID            int
	PaymentPeriod string
	PortfolioID   int
	ShortName     string
	TaxPaid       float64
	TotalPayment  float64
	Ticker        string
	UserID        int
}
