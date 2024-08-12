package dto

import (
	"time"

	"github.com/pttrulez/investor-go/internal/entity"
)

type CreateTransaction struct {
	Amount      int                    `json:"amount"      validate:"required"`
	Date        time.Time              `json:"date"        validate:"required"`
	PortfolioID int                    `json:"portfolioId" validate:"required"`
	Type        entity.TransactionType `json:"type"        validate:"required,transactionType"`
}
