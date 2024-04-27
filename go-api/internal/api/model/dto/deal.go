package dto

import (
	"time"

	"github.com/pttrulez/investor-go/internal/model"
)

type CreateDeal struct {
	Amount       int                `json:"amount" validate:"required"`
	Date         time.Time          `json:"date" validate:"required"`
	Exchange     model.Exchange     `json:"exchange" validate:"required,exchange"`
	PortfolioId  int                `json:"portfolioId" validate:"required"`
	Price        float64            `json:"price" validate:"required,price"`
	Ticker       string             `json:"ticker" validate:"required"`
	SecurityType model.SecurityType `json:"securityType" validate:"required,securityType"`
	Type         model.DealType     `json:"type" validate:"required,dealType"`
}

type DeleteDeal struct {
	Id           int                `json:"id" validate:"required"`
	Exchange     model.Exchange     `json:"exchange" validate:"required,exchange"`
	SecurityType model.SecurityType `json:"securityType" validate:"required,securityType"`
}
