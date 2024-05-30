package portfolio

import (
	"context"
	"fmt"
	"github.com/pttrulez/investor-go/internal/controller/model/converter"
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"github.com/pttrulez/investor-go/internal/types"
	"github.com/pttrulez/investor-go/internal/utils"
)

func (s *Service) GetFullPortfolioById(ctx context.Context, portfolioId int,
	userId int) (*response.FullPortfolio, error) {
	var (
		done           = make(chan struct{})
		errCh          = make(chan error)
		p              entity.Portfolio
		result         *response.FullPortfolio
		cashouts       []entity.Cashout
		deposits       []entity.Deposit
		bondDeals      []*entity.Deal
		shareDeals     []*entity.Deal
		bondPositions  []*entity.Position
		sharePositions []*entity.Position
	)

	// Базовая инфа о портфолио
	go func() {
		pdb, err := s.GetPortfolioById(ctx, portfolioId, userId)
		if err != nil {
			errCh <- fmt.Errorf("[PortfolioService.GetPortfolio]: %w", err)
			return
		}

		p.Id = pdb.Id
		p.Compound = pdb.Compound
		p.Name = pdb.Name

		done <- struct{}{}
	}()

	// сделки
	go func() {
		deals, err := s.dealRepository.GetDealListByPortoflioId(ctx, portfolioId, userId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		var bondDealsCount, shareDealsCount int
		for _, d := range deals {
			if d.SecurityType == types.ST_Bond {
				bondDealsCount++
			} else if d.SecurityType == types.ST_Share {
				shareDealsCount++
			}
		}
		for _, d := range deals {
			if d.SecurityType == types.ST_Share {
				shareDeals = append(shareDeals, d)
			} else if d.SecurityType == types.ST_Bond {
				bondDeals = append(bondDeals, d)
			}
		}

		done <- struct{}{}
	}()

	// позиции
	go func() {
		// получили массив позиций по айди портфолио
		positions, err := s.positionRepo.GetListByPortfolioId(ctx, portfolioId, userId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		var bondPositionsCount, sharePositionsCount int
		for _, p := range positions {
			if p.SecurityType == types.ST_Bond {
				bondPositionsCount++
			} else if p.SecurityType == types.ST_Share {
				sharePositionsCount++
			}
		}
		sharePositions = make([]*entity.Position, bondPositionsCount)
		bondPositions = make([]*entity.Position, sharePositionsCount)
		bondTickers := make([]string, bondPositionsCount)
		shareTickers := make([]string, sharePositionsCount)
		for _, p := range positions {
			if p.SecurityType == types.ST_Share {
				sharePositions = append(sharePositions, p)
				shareTickers = append(shareTickers, p.Ticker)
			} else if p.SecurityType == types.ST_Bond {
				bondPositions = append(bondPositions, p)
				bondTickers = append(bondTickers, p.Ticker)
			}
		}
		var bondPrices, sharePrices iss_client.Prices

		if bondPositionsCount > 0 {
			bondPrices, err = s.issClient.GetStocksCurrentPrices(ctx, entity.MoexMarketBonds, bondTickers)
			if err != nil {
				errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
				return
			}
		}
		if sharePositionsCount > 0 {
			sharePrices, err = s.issClient.GetStocksCurrentPrices(ctx, entity.MoexMarketBonds, shareTickers)
			if err != nil {
				errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
				return
			}
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		for i := 0; i < len(bondPositions); i++ {
			bondPositions[i].CurrentPrice = bondPrices[bondPositions[i].Ticker][bondPositions[i].Board]
			bondPositions[i].CurrentCost = int(bondPositions[i].CurrentPrice * float64(bondPositions[i].Amount))
		}
		for i := 0; i < len(sharePositions); i++ {
			sharePositions[i].CurrentPrice = sharePrices[sharePositions[i].Ticker][sharePositions[i].Board]
			sharePositions[i].CurrentCost = int(sharePositions[i].CurrentPrice * float64(sharePositions[i].Amount))
		}

		done <- struct{}{}
	}()

	// кэшауты
	go func() {
		var err error
		cashouts, err = s.cashoutRepo.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("[PortfolioService.GetPortfolio]: %w", err)
			return
		}

		done <- struct{}{}
	}()

	// депозиты
	go func() {
		var err error
		deposits, err = s.depositRepo.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("[PortfolioService.GetPortfolio]: %w", err)
			return
		}

		done <- struct{}{}
	}()

	routinesWorking := 5

	for routinesWorking > 0 {
		select {
		case <-done:
			routinesWorking--
		case err := <-errCh:
			return nil, err
		}
	}

	for _, c := range cashouts {
		result.CashoutsSum += c.Amount
		result.Cashouts = append(result.Cashouts, converter.FromCashoutToResponse(c))
	}
	for _, d := range deposits {
		result.DepositsSum += d.Amount
		result.Deposits = append(result.Deposits, converter.FromDepositToResponse(d))
	}

	spentToBuys := 0
	receivedFromSells := 0

	// сделки
	result.BondDeals = make([]response.Deal, len(bondDeals))
	result.ShareDeals = make([]response.Deal, len(shareDeals))
	for _, d := range shareDeals {
		if d.Type == entity.DtBuy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
	}
	for _, d := range bondDeals {
		if d.Type == entity.DtBuy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
	}

	// позиции
	result.BondPositions = make([]response.Position, len(bondPositions))
	result.SharePositions = make([]response.Position, len(bondPositions))

	for i, pos := range bondPositions {
		result.TotalCost += pos.CurrentCost
		result.BondPositions[i] = converter.FromPositionToResponse(pos)
	}
	for i, pos := range sharePositions {
		result.TotalCost += pos.CurrentCost
		result.SharePositions[i] = converter.FromPositionToResponse(pos)
	}

	result.Cash = result.DepositsSum - result.CashoutsSum - spentToBuys + receivedFromSells
	result.TotalCost += result.Cash
	result.Profitability = int((float64(result.TotalCost+result.CashoutsSum+result.Cash)/float64(result.DepositsSum) - 1) * 100)

	return result, nil
}

func (s *Service) CreatePortfolio(ctx context.Context, p *entity.Portfolio) error {
	return s.repo.Insert(ctx, p)
}

func (s *Service) GetListByUserId(ctx context.Context, userId int) ([]*entity.Portfolio, error) {
	return s.repo.GetListByUserId(ctx, userId)
}

func (s *Service) GetPortfolioById(ctx context.Context, portfolioId int,
	userId int) (*entity.Portfolio, error) {
	return s.repo.GetById(ctx, portfolioId, userId)
}

func (s *Service) DeletePortfolio(ctx context.Context, portfolioId int, userId int) error {
	if _, err := s.GetPortfolioById(ctx, portfolioId, userId); err != nil {
		return err
	}
	return s.repo.Delete(ctx, portfolioId, userId)
}

func (s *Service) UpdatePortfolio(ctx context.Context, portfolio *entity.Portfolio, userId int) error {
	return s.repo.Update(ctx, portfolio, utils.GetCurrentUserId(ctx))
}

type Service struct {
	cashoutRepo    CashoutRepository
	dealRepository DealRepository
	depositRepo    DepositRepository
	issClient      iss_client.IssClient
	positionRepo   PositionRepository
	repo           Repository
}

type Repository interface {
	Delete(ctx context.Context, id int, userId int) error
	GetById(ctx context.Context, id int, userId int) (*entity.Portfolio, error)
	GetListByUserId(ctx context.Context, id int) ([]*entity.Portfolio, error)
	Insert(ctx context.Context, p *entity.Portfolio) error
	Update(ctx context.Context, p *entity.Portfolio, userId int) error
}
type DepositRepository interface {
	GetListByPortfolioId(ctx context.Context, portfolioId int) ([]entity.Deposit, error)
}
type CashoutRepository interface {
	GetListByPortfolioId(ctx context.Context, portfolioId int) ([]entity.Cashout, error)
}
type DealRepository interface {
	GetDealListByPortoflioId(ctx context.Context, portfolioId int, userId int) ([]*entity.Deal, error)
}
type PositionRepository interface {
	GetListByPortfolioId(ctx context.Context, portfolioId int, userId int) ([]*entity.Position, error)
}

func NewPortfolioService(
	cashoutRepo CashoutRepository,
	dealRepository DealRepository,
	depositRepo DepositRepository,
	issClient iss_client.IssClient,
	positionRepo PositionRepository,
	repository Repository,
) *Service {
	return &Service{
		cashoutRepo:    cashoutRepo,
		dealRepository: dealRepository,
		depositRepo:    depositRepo,
		issClient:      issClient,
		positionRepo:   positionRepo,
		repo:           repository,
	}
}
