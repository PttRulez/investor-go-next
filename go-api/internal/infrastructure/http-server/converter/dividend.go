package converter

import (
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
		Date:            req.Date.Time,
		Exchange:        exch,
		PaymentPeriod:   req.PaymentPeriod,
		PaymentPerShare: req.PaymentPerShare,
		PortfolioID:     req.PortfolioId,
		Ticker:          req.Ticker,
		SharesCount:     req.SharesCount,
	}, nil
}
