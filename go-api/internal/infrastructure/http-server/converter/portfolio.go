package converter

import (
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
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

	coupons := make([]contracts.CouponResponse, 0, len(portfolio.Coupons))
	for _, d := range portfolio.Coupons {
		coupons = append(coupons, FromCouponToCouponResponse(d))
	}

	deals := make([]contracts.DealResponse, 0, len(portfolio.Deals))
	for _, d := range portfolio.Deals {
		deals = append(deals, FromDealToResponse(d))
	}

	dividends := make([]contracts.DividendResponse, 0, len(portfolio.Dividends))
	for _, d := range portfolio.Dividends {
		dividends = append(dividends, FromDividendToDividendResponse(d))
	}

	sharePositions := make([]contracts.PositionResponse, 0, len(portfolio.SharePositions))
	for _, p := range portfolio.SharePositions {
		sharePositions = append(sharePositions, FromPositionToResponse(p))
	}

	transactions := make([]contracts.TransactionResponse, 0, len(portfolio.Transactions))
	for _, t := range portfolio.Transactions {
		transactions = append(transactions, FromTransactionToResponse(t))
	}

	return contracts.FullPortfolioResponse{
		BondPositions:  bondPositions,
		Cash:           portfolio.Cash,
		CashoutsSum:    portfolio.CashoutsSum,
		Compound:       portfolio.Compound,
		Coupons:        coupons,
		CouponsSum:     portfolio.CouponsSum,
		Deals:          deals,
		DepositsSum:    portfolio.DepositsSum,
		Dividends:      dividends,
		DividendsSum:   portfolio.DividendsSum,
		ExpensesSum:    portfolio.ExpensesSum,
		Id:             portfolio.ID,
		Name:           portfolio.Name,
		Profitability:  portfolio.Profitability,
		SharePositions: sharePositions,
		TotalCost:      portfolio.TotalCost,
		Transactions:   transactions,
	}
}
