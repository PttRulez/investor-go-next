package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateDeal(ctx context.Context, d domain.Deal, userID int) (domain.Deal, error) {
	const op = "DealService.Create"
	var res domain.Deal

	err := s.repo.ExecAsTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		decimalCount, err := s.checkSecurity(ctx, d.Exchange, d.SecurityType, d.Ticker)
		if err != nil {
			return err
		}

		res, err = s.repo.InsertDeal(ctx, tx, d)
		if err != nil {
			return err
		}

		err = s.updatePositionInDB(ctx, tx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker,
			userID, decimalCount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *Service) DeleteDealByID(ctx context.Context, id int, userID int) error {
	const op = "DealService.DeleteDealByID"

	err := s.repo.ExecAsTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
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

		err = s.updatePositionInDB(ctx, tx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker,
			userID, decimalCount)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}

func (s *Service) createNewPosition(ctx context.Context, exchange domain.Exchange, portfolioID int,
	ticker string, securityType domain.SecurityType) (domain.Position, error) {
	const op = "DealService.createNewPosition"

	position := domain.Position{
		PortfolioID:  portfolioID,
		Ticker:       ticker,
		SecurityType: securityType,
	}

	if exchange == domain.EXCHMoex && securityType == domain.STShare {
		share, err := s.moexService.GetShareByTicker(ctx, ticker)
		if err != nil {
			return domain.Position{}, fmt.Errorf("%s: %w", op, err)
		}
		position.Board = share.Board
		position.ShortName = share.ShortName
	} else if exchange == domain.EXCHMoex && securityType == domain.STBond {
		bond, err := s.moexService.GetBondByTicker(ctx, ticker)
		if err != nil {
			return domain.Position{}, fmt.Errorf("%s: %w", op, err)
		}
		position.Board = bond.Board
		position.ShortName = bond.ShortName
	}

	return position, nil
}
