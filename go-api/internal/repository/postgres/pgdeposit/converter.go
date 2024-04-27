package pgdeposit

import "github.com/pttrulez/investor-go/internal/model"

func FromModelToDBDeposit(d *model.Deposit) *Deposit {
	return &Deposit{
		Id:          d.Id,
		Amount:      d.Amount,
		Date:        d.Date,
		PortfolioId: d.PortfolioId,
	}
}

func FromDBToModelDeposit(d *Deposit) *model.Deposit {
	return &model.Deposit{
		Id:          d.Id,
		Amount:      d.Amount,
		Date:        d.Date,
		PortfolioId: d.PortfolioId,
	}
}
