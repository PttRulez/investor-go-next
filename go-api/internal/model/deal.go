package model

import "time"

type Deal struct {
	Amount       int
	Date         time.Time
	Exchange     Exchange
	Id           int
	PortfolioId  int
	Price        float64
	Ticker       string
	SecurityId   int
	SecurityType SecurityType
	Type         DealType
}
