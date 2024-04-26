package dto

import "time"

type CreateDeposit struct {
	Amount      int       `json:"amount"      validate:"required"`
	Date        time.Time `json:"date"        validate:"required"`
	PortfolioId int       `json:"portfolioId" validate:"required"`
}
