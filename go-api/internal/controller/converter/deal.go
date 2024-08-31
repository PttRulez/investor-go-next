package converter

import (
	"context"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/utils"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreateDealRequestToDeal(ctx context.Context, req api.CreateDealRequest) (entity.Deal, error) {
	exch, err := exchange(req.Exchange)
	if err != nil {
		return entity.Deal{}, err
	}

	secType, err := securityType(req.SecurityType)
	if err != nil {
		return entity.Deal{}, err
	}

	dealType, err := dealType(req.Type)
	if err != nil {
		return entity.Deal{}, err
	}

	userID := utils.GetCurrentUserID(ctx)

	return entity.Deal{
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

func FromDealToResponse(deal entity.Deal) api.DealResponse {
	return api.DealResponse{
		Amount:       deal.Amount,
		Comission:    deal.Commission,
		Date:         deal.Date,
		Exchange:     api.Exchange(deal.Exchange),
		Id:           &deal.ID,
		PortfolioId:  deal.PortfolioID,
		Price:        deal.Price,
		Ticker:       deal.Ticker,
		SecurityType: api.SecurityType(deal.SecurityType),
		Type:         api.DealType(deal.Type),
	}
}
