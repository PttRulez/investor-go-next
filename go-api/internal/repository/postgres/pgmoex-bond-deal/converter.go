package pgmoexbonddeal

import "github.com/pttrulez/investor-go/internal/model"

func FromDBToDeal(moexBondDeal *MoexBondDeal) *model.Deal {
	return &model.Deal{
		Amount:       moexBondDeal.Amount,
		Date:         moexBondDeal.Date,
		Exchange:     model.EXCH_Moex,
		Id:           moexBondDeal.Id,
		PortfolioId:  moexBondDeal.PortfolioId,
		Price:        moexBondDeal.Price,
		Ticker:       moexBondDeal.Isin,
		SecurityId:   moexBondDeal.SecurityId,
		SecurityType: model.ST_Bond,
		Type:         moexBondDeal.Type,
	}
}

func FromDealToDB(deal *model.Deal) *MoexBondDeal {
	return &MoexBondDeal{
		Amount:      deal.Amount,
		Date:        deal.Date,
		Id:          deal.Id,
		PortfolioId: deal.PortfolioId,
		Price:       deal.Price,
		SecurityId:  deal.SecurityId,
		Type:        deal.Type,
		Isin:        deal.Ticker,
	}
}
