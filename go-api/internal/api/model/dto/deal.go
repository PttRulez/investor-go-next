package dto

import (
	"time"

	"github.com/pttrulez/investor-go/internal/model"
)

type CreateDeal struct {
	Amount      int            `json:"amount"       validate:"required"`
	Date        time.Time      `json:"date"         validate:"required"`
	Exchange    model.Exchange `json:"exchange"     validate:"required,is-exchange"`
	PortfolioId int            `json:"portfolioId"  validate:"required"`
	Price       float64        `json:"price"        validate:"required,price"`
	Secid       string         `json:"secid"        validate:"required"`
	Type        model.DealType `json:"type"         validate:"required,dealType"`
}
