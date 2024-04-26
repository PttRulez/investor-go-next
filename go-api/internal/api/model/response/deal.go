package response

import (
	"time"

	"github.com/pttrulez/investor-go/internal/model"
)

type MoexBondDeal struct {
	Amount      int            `json:"amount"`
	Date        time.Time      `json:"date"`
	Id          int            `json:"id"`
	PortfolioId int            `json:"portfolio_id"`
	Price       float64        `json:"price"`
	SecurityId  int            `json:"security_id"`
	Type        model.DealType `json:"type"`
	Isin        string         `json:"isin"`
}

type MoexShareDeal struct {
	Amount      int            `json:"amount"`
	Date        time.Time      `json:"date"`
	Id          int            `json:"id"`
	PortfolioId int            `json:"portfolio_id"`
	Price       float64        `json:"price"`
	SecurityId  int            `json:"security_id"`
	Type        model.DealType `json:"type"`
	Ticker      string         `json:"ticker"`
}
