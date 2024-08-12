package deal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/issclient"
	"github.com/pttrulez/investor-go/internal/service"
	"github.com/pttrulez/investor-go/internal/service/moexbond"
	"github.com/pttrulez/investor-go/internal/service/moexshare"
)

func (s *Service) CreateDeal(ctx context.Context, d *entity.Deal) error {
	const op = "DealService.Create"

	// Этот шаг создает запись в бд с бумагой, если её ещё нет
	var err error
	if d.Exchange != entity.EXCHMoex && d.SecurityType == entity.STShare {
		_, err = s.moexShareService.GetBySecid(ctx, d.Ticker)
	} else if d.Exchange != entity.EXCHMoex && d.SecurityType == entity.STBond {
		_, err = s.moexBondService.GetBySecid(ctx, d.Ticker)
	}
	if err != nil {
		return err
	}

	err = s.repo.Insert(ctx, d)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.UpdatePositionInDB(ctx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker, d.UserID)
	if err != nil {
		return fmt.Errorf("[DealService.Create]: %w", err)
	}
	return nil
}

func (s *Service) UpdatePositionInDB(ctx context.Context, portfolioID int, exchange entity.Exchange,
	securityType entity.SecurityType, ticker string, userID int) error {
	const op = "DealService.UpdatePositionInDB"

	allDeals, err := s.repo.GetDealListForSecurity(ctx, exchange, portfolioID, securityType, ticker)
	if err != nil {
		return fmt.Errorf("DealService.UpdatePositionInDB -> %w", err)
	}

	var position *entity.Position
	oldPosition, err := s.positionRepo.GetForSecurity(
		ctx,
		exchange,
		portfolioID,
		securityType,
		ticker,
	)
	if err != nil {
		return err
	}
	if oldPosition != nil {
		position = oldPosition
	} else {
		position, err = s.getPosition(ctx, exchange, ticker, securityType)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		position.Exchange = exchange
		position.PortfolioID = portfolioID
		position.SecurityType = securityType
		position.Ticker = ticker
		position.UserID = userID
	}

	// Calculate and add to position Amount
	var amount int
	var totalAmount int
	for _, deal := range allDeals {
		amount = deal.Amount
		if deal.Type == entity.DTSell {
			amount = -amount
		}
		totalAmount += amount
	}
	position.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	left := position.Amount
	for _, deal := range allDeals {
		if deal.Type == entity.DTSell {
			continue
		}
		if left > deal.Amount {
			m[deal.Price] += deal.Amount
			left -= deal.Amount
		} else {
			m[deal.Price] += left
			break
		}
	}

	var avPrice float64
	const hundred = 100
	for price, amount := range m {
		avPrice += float64(amount) / float64(position.Amount) * price
	}
	position.AveragePrice = math.Floor(avPrice*hundred) / hundred

	if position.ID == 0 {
		err = s.positionRepo.Insert(ctx, position)
		if err != nil {
			return fmt.Errorf("DealService.UpdatePositionInDB -> %w", err)
		}
	} else {
		err = s.positionRepo.Update(ctx, position)
		if err != nil {
			return fmt.Errorf("DealService.UpdatePositionInDB -> %w", err)
		}
	}

	return nil
}

func (s *Service) getPosition(ctx context.Context, exchange entity.Exchange, ticker string,
	securityType entity.SecurityType) (*entity.Position, error) {
	const op = "getPosition"

	position := new(entity.Position)
	if exchange == entity.EXCHMoex && securityType == entity.STShare {
		share, err := s.moexShareService.GetBySecid(ctx, ticker)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		position.Board = share.Board
	} else if exchange == entity.EXCHMoex && securityType == entity.STBond {
		bond, err := s.moexBondService.GetBySecid(ctx, ticker)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		position.Board = bond.Board
	}

	return position, nil
}

func (s *Service) DeleteDealByID(ctx context.Context, id int, userID int) error {
	d, err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.ErrEntityNotFound
		}
		return fmt.Errorf("DealService.DeleteDealById -> %w", err)
	}
	err = s.UpdatePositionInDB(ctx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker, d.UserID)
	if err != nil {
		return fmt.Errorf("DealService.DeleteDealById -> %w", err)
	}

	return nil
}

type Repository interface {
	Insert(ctx context.Context, d *entity.Deal) error
	GetDealListForSecurity(ctx context.Context, exchange entity.Exchange, portfolioID int,
		securityType entity.SecurityType, ticker string) ([]*entity.Deal, error)
	Delete(ctx context.Context, id int, userID int) (*entity.Deal, error)
}
type PositionRepo interface {
	GetForSecurity(ctx context.Context, exchange entity.Exchange, portfolioID int,
		securityType entity.SecurityType, ticker string) (*entity.Position, error)
	Insert(ctx context.Context, p *entity.Position) error
	Update(ctx context.Context, p *entity.Position) error
}

type Service struct {
	issClient        *issclient.IssClient
	moexBondService  *moexbond.Service
	moexShareService *moexshare.Service
	positionRepo     PositionRepo
	repo             Repository
}

func NewDealService(
	issClient *issclient.IssClient,
	moexBondService *moexbond.Service,
	moexShareService *moexshare.Service,
	repo Repository,
) *Service {
	return &Service{
		issClient:        issClient,
		moexBondService:  moexBondService,
		moexShareService: moexShareService,
		repo:             repo,
	}
}
