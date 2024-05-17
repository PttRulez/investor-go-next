package response

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"time"
)

type MoexBondDeal struct {
	Amount      int         `json:"amount"`
	Date        time.Time   `json:"date"`
	Id          int         `json:"id"`
	PortfolioId int         `json:"portfolio_id"`
	Price       float64     `json:"price"`
	SecurityId  int         `json:"security_id"`
	Type        entity.Type `json:"type"`
	Isin        string      `json:"isin"`
}

type MoexShareDeal struct {
	Amount      int         `json:"amount"`
	Date        time.Time   `json:"date"`
	Id          int         `json:"id"`
	PortfolioId int         `json:"portfolio_id"`
	Price       float64     `json:"price"`
	SecurityId  int         `json:"security_id"`
	Type        entity.Type `json:"type"`
	Ticker      string      `json:"ticker"`
}
