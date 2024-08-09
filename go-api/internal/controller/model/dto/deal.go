package dto

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"time"
)

type CreateDeal struct {
	Amount       int                 `json:"amount" validate:"required"`
	Date         time.Time           `json:"date" validate:"required"`
	Exchange     entity.Exchange     `json:"exchange" validate:"required,exchange"`
	PortfolioId  int                 `json:"portfolioId" validate:"required"`
	Price        float64             `json:"price" validate:"required,price"`
	SecurityType entity.SecurityType `json:"securityType" validate:"required,securityType"`
	Ticker       string              `json:"ticker" validate:"required"`
	Type         entity.DealType     `json:"type" validate:"required,dealType"`
}

type DeleteDeal struct {
	Id           int                 `json:"id" validate:"required"`
	Exchange     entity.Exchange     `json:"exchange" validate:"required,exchange"`
	SecurityType entity.SecurityType `json:"securityType" validate:"required,securityType"`
}
