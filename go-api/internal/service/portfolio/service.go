package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/infrastracture/issclient"
	"github.com/pttrulez/investor-go/internal/service"
	"github.com/pttrulez/investor-go/internal/utils"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetFullPortfolioByID(ctx context.Context, portfolioID int,
	userID int) (entity.Portfolio, error) {
	const op = "PortfolioService.GetFullPortfolioByID"
	var (
		result         entity.Portfolio
		deals          []entity.Deal
		bondPositions  = make([]entity.Position, 0)
		sharePositions = make([]entity.Position, 0)
		transactions   = make([]entity.Transaction, 0)
		eg             = errgroup.Group{}
	)

	// Базовая инфа о портфолио
	eg.Go(func() error {
		p, err := s.GetPortfolioByID(ctx, portfolioID, userID)
		if errors.Is(err, sql.ErrNoRows) {
			return service.ErrEntityNotFound
		}
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		result.ID = p.ID
		result.Compound = p.Compound
		result.Name = p.Name

		return nil
	})

	// сделки
	eg.Go(func() error {
		var err error
		deals, err = s.dealRepository.GetDealListByPortoflioID(ctx, portfolioID, userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	// позиции
	eg.Go(func() error {
		// получили массив позиций по айди портфолио
		positions, err := s.positionRepo.GetListByPortfolioID(ctx, portfolioID, userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		var bondPositionsCount, sharePositionsCount int
		for _, p := range positions {
			if p.SecurityType == entity.STBond {
				bondPositionsCount++
			} else if p.SecurityType == entity.STShare {
				sharePositionsCount++
			}
		}

		sharePositions = make([]entity.Position, 0, bondPositionsCount)
		bondPositions = make([]entity.Position, 0, sharePositionsCount)
		bondBoards := make(map[string]entity.ISSMoexBoard)
		shareBoards := make(map[string]entity.ISSMoexBoard)

		for _, p := range positions {
			if p.SecurityType == entity.STShare {
				sharePositions = append(sharePositions, p)
				shareBoards[p.Secid] = p.Board
			} else if p.SecurityType == entity.STBond {
				bondPositions = append(bondPositions, p)
				bondBoards[p.Secid] = p.Board
			}
		}

		var bondPrices, sharePrices map[string]float64

		if bondPositionsCount > 0 {
			bondPrices, err = s.issClient.GetStocksCurrentPrices(ctx, entity.MoexMarketBonds,
				bondBoards)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		if sharePositionsCount > 0 {
			sharePrices, err = s.issClient.GetStocksCurrentPrices(ctx, entity.MoexMarketShares,
				shareBoards)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		const faceValue = 1000
		const hundredPercents = 100
		for i := range len(bondPositions) {
			bondPositions[i].CurrentPrice = bondPrices[bondPositions[i].Secid]
			bondPositions[i].CurrentCost = int(
				(bondPositions[i].CurrentPrice / hundredPercents) * faceValue *
					float64(bondPositions[i].Amount))
		}

		fmt.Printf("sharePrices %#v", sharePrices)

		for i := range len(sharePositions) {
			sharePositions[i].CurrentPrice = sharePrices[sharePositions[i].Secid]
			sharePositions[i].CurrentCost = int(
				sharePositions[i].CurrentPrice * float64(sharePositions[i].Amount))
		}

		return nil
	})

	// транзакции
	eg.Go(func() error {
		var err error
		transactions, err = s.transactionRepo.GetListByPortfolioID(ctx, portfolioID, utils.GetCurrentUserID(ctx))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		return entity.Portfolio{}, err
	}

	for _, t := range transactions {
		if t.Type == entity.TTCashout {
			result.CashoutsSum += t.Amount
		} else if t.Type == entity.TTDeposit {
			result.DepositsSum += t.Amount
		}
		result.Transactions = append(result.Transactions, t)
	}

	var spentToBuys int
	var receivedFromSells int
	var spentToComissions int

	// сделки
	for _, d := range deals {
		if d.Type == entity.DTBuy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
		spentToComissions += int(d.Commission)
	}

	// позиции
	result.BondPositions = make([]entity.Position, len(bondPositions))
	result.SharePositions = make([]entity.Position, len(sharePositions))

	for i := range len(bondPositions) {
		result.TotalCost += bondPositions[i].CurrentCost
		result.BondPositions[i] = bondPositions[i]
	}

	for i := range len(sharePositions) {
		result.TotalCost += sharePositions[i].CurrentCost
		result.SharePositions[i] = sharePositions[i]
	}

	const percentageMultiplier = 100
	result.Cash = result.DepositsSum - result.CashoutsSum - spentToBuys + receivedFromSells
	result.TotalCost += result.Cash
	result.TotalCost -= spentToComissions
	result.Profitability = int((float64(result.TotalCost+result.CashoutsSum)/
		float64(result.DepositsSum) - 1) * percentageMultiplier)

	return result, nil
}

func (s *Service) CreatePortfolio(ctx context.Context, p entity.Portfolio) (entity.Portfolio, error) {
	const op = "PortfolioPostgres.Insert"

	portfolio, err := s.repo.Insert(ctx, p)
	if err != nil {
		return entity.Portfolio{}, fmt.Errorf("%s:, %w", op, err)
	}

	return portfolio, nil
}

func (s *Service) GetListByUserID(ctx context.Context, userID int) ([]entity.Portfolio,
	error) {
	return s.repo.GetListByUserID(ctx, userID)
}

func (s *Service) GetPortfolioByID(ctx context.Context, portfolioID int,
	userID int) (entity.Portfolio, error) {
	const op = "PortfolioService.GetPortfolioByID"

	portfolio, err := s.repo.GetByID(ctx, portfolioID, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return entity.Portfolio{}, service.ErrEntityNotFound
		}
		return entity.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return portfolio, nil
}

func (s *Service) DeletePortfolio(ctx context.Context, portfolioID int, userID int) error {
	const op = "PortfolioService.DeletePortfolio"

	err := s.repo.Delete(ctx, portfolioID, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrEntityNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

func (s *Service) UpdatePortfolio(ctx context.Context, portfolio entity.Portfolio,
	userID int) (entity.Portfolio, error) {
	const op = "PortfolioService.UpdatePortfolio"

	updatedPortoflio, err := s.repo.Update(ctx, portfolio, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return entity.Portfolio{}, service.ErrEntityNotFound
		}
		return entity.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedPortoflio, err
}

type Service struct {
	dealRepository  DealRepository
	issClient       *issclient.IssClient
	positionRepo    PositionRepository
	repo            Repository
	transactionRepo TransactionRepository
}

type Repository interface {
	Delete(ctx context.Context, id int, userID int) error
	GetByID(ctx context.Context, id int, userID int) (entity.Portfolio, error)
	GetListByUserID(ctx context.Context, id int) ([]entity.Portfolio, error)
	Insert(ctx context.Context, p entity.Portfolio) (entity.Portfolio, error)
	Update(ctx context.Context, p entity.Portfolio, userID int) (entity.Portfolio, error)
}
type DealRepository interface {
	GetDealListByPortoflioID(ctx context.Context, portfolioID int, userID int) ([]entity.Deal, error)
}
type PositionRepository interface {
	GetListByPortfolioID(ctx context.Context, portfolioID int, userID int) ([]entity.Position, error)
}
type TransactionRepository interface {
	GetListByPortfolioID(ctx context.Context, portfolioID int, userID int) ([]entity.Transaction, error)
}

func NewPortfolioService(
	dealRepository DealRepository,
	issClient *issclient.IssClient,
	positionRepo PositionRepository,
	repository Repository,
	transactionRepo TransactionRepository,
) *Service {
	return &Service{
		dealRepository:  dealRepository,
		transactionRepo: transactionRepo,
		issClient:       issClient,
		positionRepo:    positionRepo,
		repo:            repository,
	}
}
