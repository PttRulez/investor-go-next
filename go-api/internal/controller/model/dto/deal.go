package dto

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"time"

	"github.com/pttrulez/investor-go/internal/types"
)

type CreateDeal struct {
	Amount       int                `json:"amount" validate:"required"`
	Date         time.Time          `json:"date" validate:"required"`
	Exchange     types.Exchange     `json:"exchange" validate:"required,exchange"`
	PortfolioId  int                `json:"portfolioId" validate:"required"`
	Price        float64            `json:"price" validate:"required,price"`
	Ticker       string             `json:"ticker" validate:"required"`
	SecurityType types.SecurityType `json:"securityType" validate:"required,securityType"`
	Type         entity.Type        `json:"type" validate:"required,dealType"`
}

type DeleteDeal struct {
	Id           int                `json:"id" validate:"required"`
	Exchange     types.Exchange     `json:"exchange" validate:"required,exchange"`
	SecurityType types.SecurityType `json:"securityType" validate:"required,securityType"`
}
