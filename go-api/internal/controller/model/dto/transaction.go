package dto

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"time"
)

type CreateTransaction struct {
	Amount      int                    `json:"amount"      validate:"required"`
	Date        time.Time              `json:"date"        validate:"required"`
	PortfolioId int                    `json:"portfolioId" validate:"required"`
	Type        entity.TransactionType `json:"type"        validate:"required,transactionType"`
}
