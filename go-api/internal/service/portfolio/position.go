package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (s *Service) AddInfoToPosition(ctx context.Context, i domain.PositionUpdateInfo) error {
	const op = "PositionService.UpdatePosition"

	err := s.repo.AddPositionInfo(ctx, i)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetPositionList(ctx context.Context, userID int) ([]domain.Position, error) {
	const op = "PositionService.GetListByUserID"

	positions, err := s.repo.GetUserPositionList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return positions, nil
}

func (s *Service) newPosition(ctx context.Context, exchange domain.Exchange, portfolioID int,
	ticker string, securityType domain.SecurityType, userID int) (domain.Position, error) {
	const op = "DealService.newPosition"

	position := domain.Position{
		PortfolioID:  portfolioID,
		Ticker:       ticker,
		SecurityType: securityType,
		UserID:       userID,
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

func (s *Service) updatePositionInDB(ctx context.Context, portfolioID int, exchange domain.Exchange,
	securityType domain.SecurityType, ticker string, userID int, decimalCount int) error {
	const op = "DealService.updatePositionInDB"

	allDeals, err := s.repo.GetDealListForSecurity(ctx, exchange, portfolioID, securityType, ticker)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var position domain.Position
	oldPosition, err := s.repo.GetPosition(
		ctx,
		exchange,
		portfolioID,
		securityType,
		ticker,
	)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}
	if errors.Is(err, storage.ErrNotFound) {
		position, err = s.newPosition(ctx, exchange, portfolioID, ticker, securityType, userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		position = oldPosition
	}

	position = position.UpdateByDeals(allDeals, decimalCount)

	if position.ID == 0 {
		err = s.repo.InsertPosition(ctx, position)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		err = s.repo.UpdatePosition(ctx, position)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
