package converter

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
)

func FromCreateCouponRequestToCoupon(
	req contracts.CreateCouponRequest) (domain.Coupon, error) {
	exch, err := exchange(req.Exchange)
	if err != nil {
		return domain.Coupon{}, err
	}

	return domain.Coupon{
		BondsCount:    req.BondsCount,
		Date:          req.Date.Time,
		Exchange:      exch,
		PaymentPeriod: req.PaymentPeriod,
		PortfolioID:   req.PortfolioId,
		TaxPaid:       *req.TaxPaid,
		Ticker:        req.Ticker,
		TotalPayment:  req.TotalPayment,
	}, nil
}

func FromCouponToCouponResponse(
	d domain.Coupon) contracts.CouponResponse {

	return contracts.CouponResponse{
		BondsCount:    d.BondsCount,
		Date:          openapi_types.Date{Time: d.Date},
		Id:            d.ID,
		PaymentPeriod: d.PaymentPeriod,
		ShortName:     d.ShortName,
		TaxPaid:       d.TaxPaid,
		TotalPayment:  d.TotalPayment,
	}
}
