package converter

import (
	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/http-server/contracts"
)

func FromCreateCouponRequestToCoupon(
	req contracts.CreateCouponRequest) (domain.Coupon, error) {
	exch, err := exchange(req.Exchange)
	if err != nil {
		return domain.Coupon{}, err
	}

	return domain.Coupon{
		BondsCount:    req.BondsCount,
		CouponAmount:  req.CouponAmount,
		Date:          req.Date.Time,
		Exchange:      exch,
		PaymentPeriod: req.PaymentPeriod,
		PortfolioID:   req.PortfolioId,
		Ticker:        req.Ticker,
	}, nil
}
