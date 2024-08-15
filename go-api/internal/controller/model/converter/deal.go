package converter

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreateDealRequestToDeal(req api.CreateDealRequest) (entity.Deal, error) {
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

	return entity.Deal{
		Amount:       req.Amount,
		Date:         req.Date.Time,
		Exchange:     exch,
		Price:        req.Price,
		SecurityType: secType,
		Ticker:       req.Ticker,
		Type:         dealType,
	}, nil
}

func FromDealToResponse(deal entity.Deal) api.DealResponse {
	return api.DealResponse{
		Amount:       deal.Amount,
		Comission:    deal.Commission,
		Date:         openapi_types.Date{Time: deal.Date},
		Exchange:     api.Exchange(deal.Exchange),
		Id:           &deal.ID,
		PortfolioId:  deal.PortfolioID,
		Price:        deal.Price,
		SecurityId:   deal.SecurityID,
		SecurityType: api.SecurityType(deal.SecurityType),
		Ticker:       deal.Ticker,
		Type:         api.DealType(deal.Type),
	}
}
