package converter

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreatePortfolioRequestToPortfolio(dto api.CreatePortfolioRequest) entity.Portfolio {
	return entity.Portfolio{
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromUpdatePortfolioRequestToPortfolio(dto api.UpdatePortfolioRequest) entity.Portfolio {
	portfolio := entity.Portfolio{
		ID: dto.Id,
	}
	if dto.Compound != nil {
		portfolio.Compound = *dto.Compound
	}
	if dto.Name != nil {
		portfolio.Name = *dto.Name
	}
	return portfolio
}

func FromPortfolioToPortfolioResponse(portfolio entity.Portfolio) api.PortfolioResponse {
	return api.PortfolioResponse{
		Id:       portfolio.ID,
		Compound: portfolio.Compound,
		Name:     portfolio.Name,
	}
}

func FromPortfolioToFullPortfolioResponse(portfolio entity.Portfolio) api.FullPortfolioResponse {
	bondPositions := make([]api.PositionResponse, 0, len(portfolio.BondPositions))
	for _, p := range portfolio.BondPositions {
		bondPositions = append(bondPositions, FromPositionToResponse(p))
	}

	sharePositions := make([]api.PositionResponse, 0, len(portfolio.SharePositions))
	for _, p := range portfolio.SharePositions {
		sharePositions = append(sharePositions, FromPositionToResponse(p))
	}

	deals := make([]api.DealResponse, 0, len(portfolio.Deals))
	for _, d := range portfolio.Deals {
		deals = append(deals, FromDealToResponse(d))
	}

	return api.FullPortfolioResponse{
		BondPositions:  bondPositions,
		Cash:           portfolio.Cash,
		CashoutSum:     portfolio.CashoutsSum,
		Compound:       portfolio.Compound,
		Deals:          deals,
		DepositsSum:    portfolio.DepositsSum,
		Id:             portfolio.ID,
		Name:           portfolio.Name,
		Profitability:  portfolio.Profitability,
		SharePositions: sharePositions,
	}
}
