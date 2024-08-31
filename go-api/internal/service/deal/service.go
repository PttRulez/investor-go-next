package deal

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/infrastracture/issclient"
	"github.com/pttrulez/investor-go/internal/service"
	"github.com/pttrulez/investor-go/internal/service/moexbond"
	"github.com/pttrulez/investor-go/internal/service/moexshare"
)

func (s *Service) CreateDeal(ctx context.Context, d entity.Deal) (entity.Deal, error) {
	const op = "DealService.Create"

	// Этот шаг создает запись в бд с бумагой, если её ещё нет
	var err error
	var decimalCount int
	fmt.Printf("CreateDeal d %+v\n", d)
	if d.Exchange == entity.EXCHMoex && d.SecurityType == entity.STShare {
		var share entity.Share
		share, err = s.moexShareService.GetByTicker(ctx, d.Ticker)
		decimalCount = share.PriceDecimals
	} else if d.Exchange == entity.EXCHMoex && d.SecurityType == entity.STBond {
		_, err = s.moexBondService.GetByTicker(ctx, d.Ticker)
		decimalCount = 2
	}
	if err != nil {
		return entity.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := s.repo.Insert(ctx, d)

	if err != nil {
		return entity.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.updatePositionInDB(ctx, d.PortfolioID, d.Exchange, d.SecurityType, d.Ticker,
		d.UserID, decimalCount)
	if err != nil {
		return entity.Deal{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *Service) DeleteDealByID(ctx context.Context, id int, userID int) error {
	const op = "DealService.DeleteDealByID"

	d, err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrEntityNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	var decimalCount = 2
	if d.SecurityType == entity.STShare {
		share, err := s.moexShareService.GetByTicker(ctx, d.Ticker)
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

func (s *Service) createNewPosition(ctx context.Context, exchange entity.Exchange, portfolioID int,
	ticker string, securityType entity.SecurityType) (entity.Position, error) {
	const op = "DealService.createNewPosition"

	position := entity.Position{
		PortfolioID:  portfolioID,
		Ticker:       ticker,
		SecurityType: securityType,
	}

	if exchange == entity.EXCHMoex && securityType == entity.STShare {
		share, err := s.moexShareService.GetByTicker(ctx, ticker)
		if err != nil {
			return entity.Position{}, fmt.Errorf("%s: %w", op, err)
		}
		position.Board = share.Board
		position.ShortName = share.ShortName
	} else if exchange == entity.EXCHMoex && securityType == entity.STBond {
		bond, err := s.moexBondService.GetByTicker(ctx, ticker)
		if err != nil {
			return entity.Position{}, fmt.Errorf("%s: %w", op, err)
		}
		position.Board = bond.Board
		position.ShortName = bond.ShortName
	}
	fmt.Println("position", position)
	return position, nil
}

func (s *Service) updatePositionInDB(ctx context.Context, portfolioID int, exchange entity.Exchange,
	securityType entity.SecurityType, ticker string, userID int, decimalCount int) error {
	const op = "DealService.updatePositionInDB"

	allDeals, err := s.repo.GetDealListForSecurity(ctx, exchange, portfolioID, securityType, ticker)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var position entity.Position
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
		if deal.Type == entity.DTSell {
			amount = -amount
		}
		totalAmount += amount
	}
	position.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	for _, deal := range allDeals {
		if deal.Type == entity.DTSell {
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

type Repository interface {
	Insert(ctx context.Context, d entity.Deal) (entity.Deal, error)
	GetDealListForSecurity(ctx context.Context, exchange entity.Exchange, portfolioID int,
		securityType entity.SecurityType, ticker string) ([]entity.Deal, error)
	Delete(ctx context.Context, id int, userID int) (entity.Deal, error)
}
type MoexShareRepo interface {
	GetByTicker(ctx context.Context, ticker string) (entity.Share, error)
}
type PositionRepo interface {
	GetPositionForSecurity(ctx context.Context, exchange entity.Exchange, portfolioID int,
		securityType entity.SecurityType, ticker string) (entity.Position, error)
	Insert(ctx context.Context, p entity.Position) error
	Update(ctx context.Context, p entity.Position) error
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
	positionRepo PositionRepo,
	repo Repository,
) *Service {
	return &Service{
		issClient:        issClient,
		moexBondService:  moexBondService,
		moexShareService: moexShareService,
		positionRepo:     positionRepo,
		repo:             repo,
	}
}
