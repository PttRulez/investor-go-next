package entity

import (
	"time"
)

type Deal struct {
	Amount       int
	Commission   float64
	Date         time.Time
	Exchange     Exchange
	Id           int
	PortfolioId  int
	Price        float64
	SecurityId   int
	SecurityType SecurityType
	Type         DealType
	Ticker       string
	UserId       int
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
