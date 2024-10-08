package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/storage"
)

func (s *Service) GetPositionList(ctx context.Context, userID int) ([]domain.Position, error) {
	const op = "PositionService.GetListByUserID"

	positions, err := s.repo.GetUserPositionList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return positions, nil
}

func (s *Service) AddInfoToPosition(ctx context.Context, i domain.PositionUpdateInfo) error {
	const op = "PositionService.UpdatePosition"

	err := s.repo.AddPositionInfo(ctx, i)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) updatePositionInDB(ctx context.Context, tx *sql.Tx, portfolioID int, exchange domain.Exchange,
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
		position, err = s.createNewPosition(ctx, exchange, portfolioID, ticker, securityType)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		position.Exchange = exchange
		position.PortfolioID = portfolioID
		position.SecurityType = securityType
		position.Ticker = ticker
		position.UserID = userID
	} else {
		position = oldPosition
	}

	// Calculate position amount
	var amount int
	var totalAmount int
	for _, d := range allDeals {
		amount = d.Amount
		if d.Type == domain.DTSell {
			amount = -amount
		}
		totalAmount += amount
	}
	position.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	for _, deal := range allDeals {
		if deal.Type == domain.DTSell {
			continue
		}
		m[deal.Price] += deal.Amount
	}

	var avPrice float64
	for price, amount := range m {
		avPrice += float64(amount) / float64(position.Amount) * price
	}

	// Приведение средней цены к определенному кол-ву знаков после запятой
	d := float64(math.Pow10(decimalCount))
	position.AveragePrice = math.Floor(avPrice*d) / d

	if position.ID == 0 {
		err = s.repo.InsertPosition(ctx, position)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		err = s.repo.UpdatePosition(ctx, tx, position)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
