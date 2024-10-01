package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"github.com/pttrulez/investor-go/internal/service"
	"github.com/pttrulez/investor-go/internal/utils"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetFullPortfolioByID(ctx context.Context, portfolioID int,
	userID int) (domain.Portfolio, error) {
	const op = "PortfolioService.GetFullPortfolioByID"
	var (
		result         domain.Portfolio
		deals          []domain.Deal
		bondPositions  = make([]domain.Position, 0)
		sharePositions = make([]domain.Position, 0)
		transactions   = make([]domain.Transaction, 0)
		eg             = errgroup.Group{}
	)

	// Базовая инфа о портфолио
	eg.Go(func() error {
		p, err := s.GetPortfolioByID(ctx, portfolioID, userID)
		if errors.Is(err, sql.ErrNoRows) {
			return service.ErrdomainNotFound
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
		deals, err = s.dealRepo.GetDealListByPortoflioID(ctx, portfolioID, userID)
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
			if p.SecurityType == domain.STBond {
				bondPositionsCount++
			} else if p.SecurityType == domain.STShare {
				sharePositionsCount++
			}
		}

		sharePositions = make([]domain.Position, 0, bondPositionsCount)
		bondPositions = make([]domain.Position, 0, sharePositionsCount)
		bondBoards := make(map[string]domain.ISSMoexBoard)
		shareBoards := make(map[string]domain.ISSMoexBoard)

		for _, p := range positions {
			if p.SecurityType == domain.STShare {
				sharePositions = append(sharePositions, p)
				shareBoards[p.Ticker] = p.Board
			} else if p.SecurityType == domain.STBond {
				bondPositions = append(bondPositions, p)
				bondBoards[p.Ticker] = p.Board
			}
		}

		var bondPrices, sharePrices map[string]float64

		if bondPositionsCount > 0 {
			bondPrices, err = s.issClient.GetStocksCurrentPrices(ctx, domain.MoexMarketBonds,
				bondBoards)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		if sharePositionsCount > 0 {
			sharePrices, err = s.issClient.GetStocksCurrentPrices(ctx, domain.MoexMarketShares,
				shareBoards)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		const faceValue = 1000
		const hundredPercents = 100
		for i := range len(bondPositions) {
			bondPositions[i].CurrentPrice = bondPrices[bondPositions[i].Ticker]
			bondPositions[i].CurrentCost = int(
				(bondPositions[i].CurrentPrice / hundredPercents) * faceValue *
					float64(bondPositions[i].Amount))
		}

		for i := range len(sharePositions) {
			sharePositions[i].CurrentPrice = sharePrices[sharePositions[i].Ticker]
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
		return domain.Portfolio{}, err
	}

	for _, t := range transactions {
		if t.Type == domain.TTCashout {
			result.CashoutsSum += t.Amount
		} else if t.Type == domain.TTDeposit {
			result.DepositsSum += t.Amount
		}
		result.Transactions = append(result.Transactions, t)
	}

	var spentToBuys int
	var receivedFromSells int
	var spentToComissions int

	// сделки
	for _, d := range deals {
		if d.Type == domain.DTBuy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
		spentToComissions += int(d.Commission)
	}

	// позиции
	result.BondPositions = make([]domain.Position, len(bondPositions))
	result.SharePositions = make([]domain.Position, len(sharePositions))

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

func (s *Service) CreatePortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error) {
	const op = "PortfolioPostgres.Insert"

	portfolio, err := s.portfolioRepo.Insert(ctx, p)
	if err != nil {
		return domain.Portfolio{}, fmt.Errorf("%s:, %w", op, err)
	}

	return portfolio, nil
}

func (s *Service) GetListByUserID(ctx context.Context, userID int) ([]domain.Portfolio,
	error) {
	return s.portfolioRepo.GetListByUserID(ctx, userID)
}

func (s *Service) GetPortfolioByID(ctx context.Context, portfolioID int,
	userID int) (domain.Portfolio, error) {
	const op = "PortfolioService.GetPortfolioByID"

	portfolio, err := s.portfolioRepo.GetByID(ctx, portfolioID, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return domain.Portfolio{}, service.ErrdomainNotFound
		}
		return domain.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return portfolio, nil
}

func (s *Service) DeletePortfolio(ctx context.Context, portfolioID int, userID int) error {
	const op = "PortfolioService.DeletePortfolio"

	err := s.portfolioRepo.Delete(ctx, portfolioID, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

func (s *Service) UpdatePortfolio(ctx context.Context, portfolio domain.Portfolio,
	userID int) (domain.Portfolio, error) {
	const op = "PortfolioService.UpdatePortfolio"

	updatedPortoflio, err := s.portfolioRepo.Update(ctx, portfolio, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return domain.Portfolio{}, service.ErrdomainNotFound
		}
		return domain.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedPortoflio, err
}
