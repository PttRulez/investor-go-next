package pgcashout

import "github.com/pttrulez/investor-go/internal/model"

func FromModelToDBCashout(c *model.Cashout) *Cashout {
	return &Cashout{
		Id:          c.Id,
		Amount:      c.Amount,
		Date:        c.Date,
		PortfolioId: c.PortfolioId,
	}
}

func FromDBToModelCashout(c *Cashout) *model.Cashout {
	return &model.Cashout{
		Id:          c.Id,
		Amount:      c.Amount,
		Date:        c.Date,
		PortfolioId: c.PortfolioId,
	}
}
