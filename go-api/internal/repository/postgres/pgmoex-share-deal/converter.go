package pgmoexsharedeal

import "github.com/pttrulez/investor-go/internal/model"

func FromDBToMoexShareDeal(moexShareDeal *MoexShareDeal) *model.Deal {
	return &model.Deal{
		Amount:       moexShareDeal.Amount,
		Date:         moexShareDeal.Date,
		Exchange:     model.EXCH_Moex,
		Id:           moexShareDeal.Id,
		PortfolioId:  moexShareDeal.PortfolioId,
		Price:        moexShareDeal.Price,
		Ticker:       moexShareDeal.Ticker,
		SecurityId:   moexShareDeal.SecurityId,
		SecurityType: model.ST_Share,
		Type:         moexShareDeal.Type,
	}
}

func FromMoexShareDealToDB(deal *model.Deal) *MoexShareDeal {
	return &MoexShareDeal{
		Amount:      deal.Amount,
		Date:        deal.Date,
		Id:          deal.Id,
		PortfolioId: deal.PortfolioId,
		Price:       deal.Price,
		SecurityId:  deal.SecurityId,
		Type:        deal.Type,
		Ticker:      deal.Ticker,
	}
}
