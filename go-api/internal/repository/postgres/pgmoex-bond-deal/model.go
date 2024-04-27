package pgmoexbonddeal

import (
	"time"

	"github.com/pttrulez/investor-go/internal/model"
)

type MoexBondDeal struct {
	Amount      int            `db:"amount"`
	Date        time.Time      `db:"date"`
	Id          int            `db:"id"`
	PortfolioId int            `db:"portfolio_id"`
	Price       float64        `db:"price"`
	SecurityId  int            `db:"security_id"`
	Type        model.DealType `db:"type"`
	Isin        string         `db:"isin"`
}
