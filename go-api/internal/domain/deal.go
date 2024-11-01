package domain

import (
	"time"
)

type Deal struct {
	Amount       int
	Commission   float64
	Date         time.Time
	Exchange     Exchange
	Nkd          *float64
	ID           int
	PortfolioID  int
	Price        float64
	Ticker       string
	SecurityID   int
	SecurityType SecurityType
	ShortName    string
	Type         DealType
	UserID       int
}

type DealType string

const (
	DTBuy  DealType = "BUY"
	DTSell DealType = "SELL"
)

func (e DealType) Validate() bool {
	switch e {
	case DTBuy:
	case DTSell:
	default:
		return false
	}
	return true
}
