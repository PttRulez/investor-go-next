package converter

import (
	"context"

	"github.com/pttrulez/investor-go/internal/api/contracts"
	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/utils"
)

func FromCreateDealRequestToDeal(ctx context.Context, req contracts.CreateDealRequest) (domain.Deal, error) {
	exch, err := exchange(req.Exchange)
	if err != nil {
		return domain.Deal{}, err
	}

	secType, err := securityType(req.SecurityType)
	if err != nil {
		return domain.Deal{}, err
	}

	dealType, err := dealType(req.Type)
	if err != nil {
		return domain.Deal{}, err
	}

	userID := utils.GetCurrentUserID(ctx)

	return domain.Deal{
		Amount:       req.Amount,
		Commission:   req.Comission,
		Date:         req.Date,
		Exchange:     exch,
		PortfolioID:  req.PortfolioId,
		Price:        req.Price,
		SecurityType: secType,
		Ticker:       req.Ticker,
		Type:         dealType,
		UserID:       userID,
	}, nil
}

func FromDealToResponse(deal domain.Deal) contracts.DealResponse {
	return contracts.DealResponse{
		Amount:       deal.Amount,
		Comission:    deal.Commission,
		Date:         deal.Date,
		Exchange:     contracts.Exchange(deal.Exchange),
		Id:           &deal.ID,
		PortfolioId:  deal.PortfolioID,
		Price:        deal.Price,
		Ticker:       deal.Ticker,
		SecurityType: contracts.SecurityType(deal.SecurityType),
		Type:         contracts.DealType(deal.Type),
	}
}
