package converter

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
)

func FromCreateDividendRequestToDividend(
	req contracts.CreateDividendRequest) (domain.Dividend, error) {
	exch, err := exchange(req.Exchange)
	if err != nil {
		return domain.Dividend{}, err
	}

	return domain.Dividend{
		Date:          req.Date.Time,
		Exchange:      exch,
		PaymentPeriod: req.PaymentPeriod,
		PortfolioID:   req.PortfolioId,
		TaxPaid:       *req.TaxPaid,
		Ticker:        req.Ticker,
		TotalPayment:  req.TotalPayment,
		SharesCount:   req.SharesCount,
	}, nil
}

func FromDividendToDividendResponse(
	d domain.Dividend) contracts.DividendResponse {

	return contracts.DividendResponse{
		Date:          openapi_types.Date{Time: d.Date},
		Id:            d.ID,
		PaymentPeriod: d.PaymentPeriod,
		SharesCount:   d.SharesCount,
		ShortName:     d.ShortName,
		TaxPaid:       d.TaxPaid,
		TotalPayment:  d.TotalPayment,
	}
}
