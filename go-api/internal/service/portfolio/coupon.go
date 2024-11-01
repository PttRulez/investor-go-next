package portfolio

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
)

func (s *Service) CreateCoupon(ctx context.Context, d domain.Coupon,
	userID int) error {
	const op = "PortfolioService.CreateCoupon"

	bond, err := s.moexService.GetBondByTicker(ctx, d.Ticker)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	d.ShortName = bond.ShortName

	err = s.repo.InsertCoupon(ctx, d, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteCoupon(ctx context.Context, couponID int, userID int) error {
	const op = "PortfolioService.DeleteCoupon"

	err := s.repo.DeleteCoupon(ctx, couponID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
