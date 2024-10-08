package converter

import (
	"github.com/pttrulez/investor-go/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go/internal/domain"
)

func FromCreatePortfolioRequestToPortfolio(dto contracts.CreatePortfolioRequest) domain.Portfolio {
	return domain.Portfolio{
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromUpdatePortfolioRequestToPortfolio(dto contracts.UpdatePortfolioRequest) domain.Portfolio {
	portfolio := domain.Portfolio{
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

func FromPortfolioToPortfolioResponse(portfolio domain.Portfolio) contracts.PortfolioResponse {
	return contracts.PortfolioResponse{
		Id:       portfolio.ID,
		Compound: portfolio.Compound,
		Name:     portfolio.Name,
	}
}

func FromPortfolioToFullPortfolioResponse(portfolio domain.Portfolio) contracts.FullPortfolioResponse {
	bondPositions := make([]contracts.PositionResponse, 0, len(portfolio.BondPositions))
	for _, p := range portfolio.BondPositions {
		bondPositions = append(bondPositions, FromPositionToResponse(p))
	}

	sharePositions := make([]contracts.PositionResponse, 0, len(portfolio.SharePositions))
	for _, p := range portfolio.SharePositions {
		sharePositions = append(sharePositions, FromPositionToResponse(p))
	}

	deals := make([]contracts.DealResponse, 0, len(portfolio.Deals))
	for _, d := range portfolio.Deals {
		deals = append(deals, FromDealToResponse(d))
	}

	return contracts.FullPortfolioResponse{
		BondPositions:  bondPositions,
		Cash:           portfolio.Cash,
		CashoutsSum:    portfolio.CashoutsSum,
		Compound:       portfolio.Compound,
		CouponsSum:     portfolio.CouponsSum,
		Deals:          deals,
		DepositsSum:    portfolio.DepositsSum,
		DividendsSum:   portfolio.DividendsSum,
		ExpensesSum:    portfolio.ExpensesSum,
		Id:             portfolio.ID,
		Name:           portfolio.Name,
		Profitability:  portfolio.Profitability,
		SharePositions: sharePositions,
		TotalCost:      portfolio.TotalCost,
	}
}
