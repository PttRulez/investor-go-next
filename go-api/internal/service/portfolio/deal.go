package portfolio

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateDeal(ctx context.Context, d domain.Deal) (domain.Deal, error) {
	const op = "DealService.Create"

	// Этот шаг создает запись в бд с бумагой, если её ещё нет
	var err error
	var decimalCount int
	fmt.Printf("CreateDeal d %+v\n", d)
	if d.Exchange == domain.EXCHMoex && d.SecurityType == domain.STShare {
		var share domain.Share
		share, err = s.moexService.GetShareByTicker(ctx, d.Ticker)
		decimalCount = share.PriceDecimals
	} else if d.Exchange == domain.EXCHMoex && d.SecurityType == domain.STBond {
		_, err = s.moexService.GetBondByTicker(ctx, d.Ticker)
		decimalCount = 2
	}
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := s.dealRepo.Insert(ctx, d)

	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.updatePositionInDB(ctx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker,
		d.UserID, decimalCount)
	if err != nil {
		return domain.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *Service) DeleteDealByID(ctx context.Context, id int, userID int) error {
	const op = "DealService.DeleteDealByID"

	d, err := s.dealRepo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
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
		d.UserID, decimalCount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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
	fmt.Println("position", position)
	return position, nil
}

func (s *Service) updatePositionInDB(ctx context.Context, portfolioID int, exchange domain.Exchange,
	securityType domain.SecurityType, ticker string, userID int, decimalCount int) error {
	const op = "DealService.updatePositionInDB"

	allDeals, err := s.dealRepo.GetDealListForSecurity(ctx, exchange, portfolioID, securityType, ticker)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var position domain.Position
	oldPosition, err := s.positionRepo.GetPositionForSecurity(
		ctx,
		exchange,
		portfolioID,
		securityType,
		ticker,
	)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}
	if errors.Is(err, database.ErrNotFound) {
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
	for _, deal := range allDeals {
		amount = deal.Amount
		if deal.Type == domain.DTSell {
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
		err = s.positionRepo.Insert(ctx, position)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		err = s.positionRepo.Update(ctx, position)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
