package types

import "time"

type Deal struct {
	Amount      int       `json:"amount" db:"amount" validate:"required"`
	Date        time.Time `json:"date" db:"date" validate:"required"`
	Id          int       `json:"id" db:"id"`
	PortfolioId int       `json:"portfolioId" db:"portfolio_id" validate:"required"`
	Price       float64   `json:"price" db:"price"  validate:"required,price"`
	SecurityId  int       `json:"securityId" db:"security_id"`
	Ticker      string    `json:"ticker" db:"ticker"`
	Type        DealType  `json:"type" db:"type" validate:"required,dealType"`
}

type DealType string

const (
	Buy  DealType = "Buy"
	Sell DealType = "Sell"
)

func (e DealType) Validate() bool {
	switch e {
	case Buy:
	case Sell:
	default:
		return false
	}
	return true
}
