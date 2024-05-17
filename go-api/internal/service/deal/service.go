package deal

import (
	"context"
	"errors"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	errors2 "github.com/pttrulez/investor-go/internal/errors"
	"math"

	"github.com/pttrulez/investor-go/internal/types"
)

func (s *Service) CreateDeal(ctx context.Context, dealInfo *entity.CreateDealInfo) error {
	if dealInfo.Exchange == types.EXCH_Moex {
		return s.CreateMoexDeal(ctx, dealInfo)
	}
	return errors.New("неправильно введена биржа")
}
func (s *Service) CreateMoexDeal(ctx context.Context, dealInfo *entity.CreateDealInfo) error {
	// Проверяем не чужой ли этой портфолио
	belongs, err := s.portfolio.BelongsToUser(ctx, dealInfo.PortfolioId, dealInfo.UserId)
	if err != nil {
		return fmt.Errorf("\n<-[DealService.CreateDeal]: %w", err)
	}
	if !belongs {
		return errors2.ErrNotYours
	}

	if dealInfo.SecurityType == types.ST_Bond {
		return s.CreateMoexBondDeal(
			ctx,
			entity.CreateDealInfo{
				Deal: dealInfo.Deal,
				Isin: dealInfo.Isin,
			},
		)
	} else if dealInfo.SecurityType == types.ST_Share {
		return s.CreateMoexShareDeal(ctx, dealInfo)
	}

	return fmt.Errorf("\n<-[DealService.CreateDeal]: Логика описана только для BOND и SHARE")
}
func (s *Service) CreateMoexBondDeal(ctx context.Context, dealInfo entity.CreateDealInfo) error {
	if dealInfo.SecurityId == 0 {
		_, _ = s.moexBond.GetByISIN(ctx, dealInfo.Isin)
	}
	err := s.moexBondDealRepo.Insert(ctx, &dealInfo.Deal)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealService.CreateDeal]: %w", err)
	}
	err = s.UpdateMoexBondPositionInDB(ctx, dealInfo.PortfolioId,
		dealInfo.SecurityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealService.CreateDeal]: %w", err)
	}

	return nil
}
func (s *Service) UpdateMoexBondPositionInDB(ctx context.Context, portfolioId int,
	securityId int) error {

	allDeals, err := s.moexBondDealRepo.GetDealListByBondId(ctx, portfolioId, securityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondService.UpdatePositionInDB]: %w", err)
	}

	var position *entity.Position
	oldPosition, err := s.moexSharePositionRepo.Get(ctx, portfolioId, securityId)
	if err != nil {
		return err
	}
	if oldPosition != nil {
		position = oldPosition
	} else {
		position = &entity.Position{}
		position.Exchange = types.EXCH_Moex
		position.PortfolioId = portfolioId
		position.SecurityId = securityId
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
		err = s.moexBondPositionRepo.Insert(ctx, position)
		if err != nil {
			return err
		}
	} else {
		err = s.moexBondPositionRepo.Update(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *Service) CreateMoexShareDeal(ctx context.Context, dealInfo *entity.CreateDealInfo) error {
	if dealInfo.SecurityId == 0 {
		_, _ = s.moexShare.GetByTicker(ctx, dealInfo.Ticker)
	}
	err := s.moexShareDealRepo.Insert(ctx, &dealInfo.Deal)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}
	err = s.UpdateMoexSharePositionInDB(ctx, dealInfo.PortfolioId, dealInfo.SecurityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}

	return nil
}
func (s *Service) UpdatePositionInDB(ctx context.Context, portfolioId int,
	securityId int) error {
	allDeals, err := s.moexShareDealRepo.GetDealListByShareId(ctx, portfolioId, securityId)
	if err != nil {
		return fmt.Errorf("\n<-[ShareService.UpdatePositionInDB]: %w", err)
	}

	var position *entity.Position
	oldPosition, err := s.moexSharePositionRepo.Get(ctx, portfolioId, securityId)
	if err != nil {
		return err
	}
	if oldPosition != nil {
		position = oldPosition
	} else {
		position = &entity.Position{}
		position.Exchange = types.EXCH_Moex
		position.PortfolioId = portfolioId
		position.SecurityId = securityId
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
		err = s.moexSharePositionRepo.Insert(ctx, position)
		if err != nil {
			return err
		}
	} else {
		err = s.moexSharePositionRepo.Update(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteDeal(ctx context.Context, dealInfo *entity.DeleteDealInfo, userId int) error {
	if dealInfo.Exchange == types.EXCH_Moex {
		return s.DeleteMoexDeal(ctx, dealInfo, userId)
	}
	return errors.New("неправильно введена биржа")
}
func (s *Service) DeleteMoexDeal(ctx context.Context, dealInfo *entity.DeleteDealInfo, userId int) error {
	if dealInfo.SecurityType == types.ST_Bond {
		return s.moexShareDealRepo.Delete(ctx, dealInfo.Id)
	} else if dealInfo.SecurityType == types.ST_Share {
		return s.DeleteMoexBondDeal(ctx, dealInfo.Id, userId)
	}

	return nil
}
func (s *Service) DeleteMoexBondDeal(ctx context.Context, dealId int, userId int) error {
	err := s.moexBondDealRepo.Delete(ctx, dealId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealService.DeleteDeal]: %w", err)
	}
	return nil
}
func (s *Service) DeleteMoexShareDeal(ctx context.Context, dealId int, userId int) error {
	err := s.moexShareDealRepo.Delete(ctx, dealId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.DeleteDeal]: %w", err)
	}
	return nil
}

type Service struct {
	moexBond              MoexBondService
	moexShare             MoexShareService
	moexBondDealRepo      MoexBondDealRepository
	moexShareDealRepo     MoexShareDealRepository
	moexBondPositionRepo  MoexBondPositionRepository
	moexSharePositionRepo MoexSharePositionRepository
	portfolio             PortfolioService
}

type PortfolioService interface {
	BelongsToUser(ctx context.Context, portfolioId int, userId int) (bool, error)
}

type MoexBondService interface {
	GetById(ctx context.Context, id int) (*entity.Bond, error)
	GetByISIN(ctx context.Context, isin string) (*entity.Bond, error)
}
type MoexShareService interface {
	GetByTicker(ctx context.Context, ticker string) (*entity.Share, error)
}

func NewDealService(portfolioService PortfolioService) *Service {
	return &Service{
		portfolio: portfolioService,
	}
}
