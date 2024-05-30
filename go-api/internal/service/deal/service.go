package deal

import (
	"context"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"github.com/pttrulez/investor-go/internal/service/moex_bond"
	"github.com/pttrulez/investor-go/internal/service/moex_share"
	"math"

	"github.com/pttrulez/investor-go/internal/types"
)

func (s *Service) CreateDeal(ctx context.Context, d *entity.Deal) error {
	if d.Exchange != types.EXCH_Moex {
		if d.SecurityType == types.ST_Share {
			_, err := s.moexShareService.GetByTicker(ctx, d.Ticker)
			if err != nil {
				return err
			}
		} else {
			_, err := s.moexBondService.GetByISIN(ctx, d.Ticker)
			if err != nil {
				return err
			}
		}
	}

	err := s.repo.Insert(ctx, d)
	if err != nil {
		return fmt.Errorf("[DealService.Create]: %w", err)
	}

	err = s.UpdatePositionInDB(ctx, d.PortfolioId, d.Exchange, d.SecurityType, d.Ticker, d.UserId)
	if err != nil {
		return fmt.Errorf("[DealService.Create]: %w", err)
	}
	return nil
}
func (s *Service) UpdatePositionInDB(ctx context.Context, portfolioId int, exchange types.Exchange,
	securityType types.SecurityType, ticker string, userId int) error {
	allDeals, err := s.repo.GetDealListForSecurity(ctx, exchange, portfolioId, securityType, ticker)
	if err != nil {
		return fmt.Errorf("[DealService.UpdatePositionInDB]: %w", err)
	}

	var position *entity.Position
	oldPosition, err := s.positionRepo.GetForSecurity(
		ctx,
		exchange,
		portfolioId,
		securityType,
		ticker,
	)
	if err != nil {
		return err
	}
	if oldPosition != nil {
		position = oldPosition
	} else {
		position = new(entity.Position)
		if exchange == types.EXCH_Moex {
			if securityType == types.ST_Share {
				share, err := s.moexShareService.GetByTicker(ctx, ticker)
				if err != nil {
					return err
				}
				position.Board = share.Board
			} else if securityType == types.ST_Bond {
				bond, err := s.moexBondService.GetByISIN(ctx, ticker)
				if err != nil {
					return err
				}
				position.Board = bond.Board
			}
		}
		position.Exchange = exchange
		position.PortfolioId = portfolioId
		position.SecurityType = securityType
		position.Ticker = ticker
		position.UserId = userId
	}

	// Calculate and add to position Amount
	var amount int
	var totalAmount int
	for _, deal := range allDeals {
		amount = deal.Amount
		if deal.Type == entity.DtSell {
			amount = -amount
		}
		totalAmount += amount
	}
	position.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	left := position.Amount
	for _, deal := range allDeals {
		if deal.Type == entity.DtSell {
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

	var avPrice float64 = 0
	for price, amount := range m {
		avPrice += float64(amount) / float64(position.Amount) * price
	}
	position.AveragePrice = math.Floor(avPrice*100) / 100

	if position.Id == 0 {
		err = s.positionRepo.Insert(ctx, position)
		if err != nil {
			return err
		}
	} else {
		err = s.positionRepo.Update(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *Service) DeleteDealById(ctx context.Context, id int, userId int) error {
	d, err := s.repo.Delete(ctx, id, userId)
	if err != nil {
		return fmt.Errorf("[DealService.DeleteDealById.repo.Delete]: %w", err)
	}
	err = s.UpdatePositionInDB(ctx, d.PortfolioId, d.Exchange, d.SecurityType, d.Ticker, d.UserId)
	if err != nil {
		return fmt.Errorf("[DealService.DeleteDealById.UpdatePositionInDB]: %w", err)
	}

	return nil
}

type Repository interface {
	Insert(ctx context.Context, d *entity.Deal) error
	GetDealListForSecurity(ctx context.Context, exchange types.Exchange, portfolioId int,
		securityType types.SecurityType, ticker string) ([]*entity.Deal, error)
	Delete(ctx context.Context, id int, userId int) (*entity.Deal, error)
}
type PositionRepo interface {
	GetForSecurity(ctx context.Context, exchange types.Exchange, portfolioId int,
		securityType types.SecurityType, ticker string) (*entity.Position, error)
	Insert(ctx context.Context, p *entity.Position) error
	Update(ctx context.Context, p *entity.Position) error
}

type Service struct {
	issClient        *iss_client.IssClient
	moexBondService  *moex_bond.Service
	moexShareService *moex_share.Service
	positionRepo     PositionRepo
	repo             Repository
}

func NewDealService(
	issClient *iss_client.IssClient,
	moexBondService *moex_bond.Service,
	moexShareService *moex_share.Service,
	repo Repository,
) *Service {
	return &Service{
		issClient:        issClient,
		moexBondService:  moexBondService,
		moexShareService: moexShareService,
		repo:             repo,
	}
}
