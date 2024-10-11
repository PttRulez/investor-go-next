package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go-next/go-api/internal/service"
)

func (s *Service) CreateDeal(ctx context.Context, d domain.Deal, userID int) (domain.Deal, error) {
	const op = "DealService.Create"

	var res domain.Deal
	err := s.repo.Transac(ctx, nil, func(ctx context.Context) error {
		decimalCount, err := s.checkSecurity(ctx, d.Exchange, d.SecurityType, d.Ticker)
		if err != nil {
			return err
		}

		res, err = s.repo.InsertDeal(ctx, d)
		if err != nil {
			return err
		}

		err = s.updatePositionInDB(ctx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker,
			userID, decimalCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	s.tg.SendMsg(ctx, fmt.Sprintf("Deal created:\n%s: %d шт.", res.Ticker, res.Amount))

	return res, nil
}

func (s *Service) DeleteDealByID(ctx context.Context, id int, userID int) error {
	const op = "DealService.DeleteDealByID"

	err := s.repo.Transac(ctx, nil, func(ctx context.Context) error {
		d, err := s.repo.DeleteDeal(ctx, id, userID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return service.ErrdomainNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		var decimalCount = 2
		if d.SecurityType == domain.STShare {
			share, err := s.moexService.GetShareByTicker(ctx, d.Ticker)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			decimalCount = share.PriceDecimals
		}

		err = s.updatePositionInDB(ctx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker,
			userID, decimalCount)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
