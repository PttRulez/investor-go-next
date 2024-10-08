package converter

import (
	"github.com/pttrulez/investor-go/internal/infrastructure/http-server/contracts"
	"github.com/pttrulez/investor-go/internal/domain"
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
