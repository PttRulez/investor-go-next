package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go-next/go-api/internal/service"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetFullPortfolioByID(ctx context.Context, portfolioID int,
	userID int) (domain.Portfolio, error) {
	const op = "PortfolioService.GetFullPortfolioByID"
	var (
		result         domain.Portfolio
		deals          []domain.Deal
		coupons        []domain.Coupon
		dividends      []domain.Dividend
		expenses       []domain.Expense
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
		deals, err = s.repo.GetDealList(ctx, portfolioID, userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	// позиции
	eg.Go(func() error {
		// получили массив позиций по айди портфолио
		positions, err := s.repo.GetPortfolioPositionList(ctx, portfolioID, userID)
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

	// депозиты и кэшауты
	eg.Go(func() error {
		var err error
		transactions, err = s.repo.GetTransactionList(ctx, portfolioID, utils.GetCurrentUserID(ctx))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	// дивиденды
	eg.Go(func() error {
		var err error
		dividends, err = s.repo.GetDividendList(ctx, portfolioID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	// купоны
	eg.Go(func() error {
		var err error
		coupons, err = s.repo.GetCouponList(ctx, portfolioID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	// другие затраты, возвраты
	eg.Go(func() error {
		var err error
		expenses, err = s.repo.GetExpenseList(ctx, portfolioID)
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

	var spentToComissions float64
	var spentToBuys float64
	var spentToOtherExpenses float64
	var receivedFromCoupons float64
	var receivedFromDividends float64
	var receivedFromSells float64

	// сделки
	for _, d := range deals {
		if d.Type == domain.DTBuy {
			spentToBuys += float64(d.Amount) * d.Price
		} else {
			receivedFromSells += float64(d.Amount) * d.Price
		}
		spentToComissions += d.Commission
	}
	result.Deals = deals

	// купоны
	for _, c := range coupons {
		receivedFromCoupons += c.CouponAmount * float64(c.BondsCount)
	}

	// дивиденды
	for _, d := range dividends {
		receivedFromDividends += d.PaymentPerShare * float64(d.SharesCount)
	}

	// прочие траты/возвраты
	for _, e := range expenses {
		spentToOtherExpenses += e.Amount
	}

	// позиции
	result.BondPositions = make([]domain.Position, len(bondPositions))
	result.SharePositions = make([]domain.Position, len(sharePositions))

	// сумма стоимости облигаций
	for i := range len(bondPositions) {
		result.BondsCost += bondPositions[i].CurrentCost
		result.BondPositions[i] = bondPositions[i]
	}

	// сумма стоимости акций
	for i := range len(sharePositions) {
		result.SharesCost += sharePositions[i].CurrentCost
		result.SharePositions[i] = sharePositions[i]
	}

	result.CouponsSum = int(receivedFromCoupons)
	result.DividendsSum = int(receivedFromDividends)
	result.ExpensesSum = int(spentToOtherExpenses)

	const percentageMultiplier = 100
	result.Cash = result.DepositsSum - result.CashoutsSum - int(spentToBuys) + int(receivedFromSells)

	result.TotalCost += result.BondsCost
	result.TotalCost += result.SharesCost
	result.TotalCost += result.Cash
	result.TotalCost -= int(spentToComissions)
	result.TotalCost += result.CouponsSum
	result.TotalCost += result.DividendsSum
	result.TotalCost -= result.ExpensesSum
	result.Profitability = int((float64(result.TotalCost+result.CashoutsSum)/
		float64(result.DepositsSum) - 1) * percentageMultiplier)

	return result, nil
}

func (s *Service) CreatePortfolio(ctx context.Context, p domain.Portfolio) (domain.Portfolio, error) {
	const op = "PortfolioPostgres.Insert"

	portfolio, err := s.repo.InsertPortfolio(ctx, p)
	if err != nil {
		return domain.Portfolio{}, fmt.Errorf("%s:, %w", op, err)
	}

	return portfolio, nil
}

func (s *Service) GetPortfolioList(ctx context.Context, userID int) ([]domain.Portfolio,
	error) {
	return s.repo.GetPortfolioList(ctx, userID)
}

func (s *Service) GetPortfolioListByChatID(ctx context.Context, chatId string) ([]domain.Portfolio, error) {
	return s.repo.GetPortfolioListByChatID(ctx, chatId)
}

func (s *Service) GetPortfolioByID(ctx context.Context, portfolioID int,
	userID int) (domain.Portfolio, error) {
	const op = "PortfolioService.GetPortfolioByID"

	portfolio, err := s.repo.GetPortfolio(ctx, portfolioID, userID)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.Portfolio{}, service.ErrdomainNotFound
		}
		return domain.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return portfolio, nil
}

func (s *Service) DeletePortfolio(ctx context.Context, portfolioID int, userID int) error {
	const op = "PortfolioService.DeletePortfolio"

	err := s.repo.DeletePortfolio(ctx, portfolioID, userID)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

func (s *Service) UpdatePortfolio(ctx context.Context, portfolio domain.Portfolio,
	userID int) (domain.Portfolio, error) {
	const op = "PortfolioService.UpdatePortfolio"

	updatedPortoflio, err := s.repo.UpdatePortfolio(ctx, portfolio, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.Portfolio{}, service.ErrdomainNotFound
		}
		return domain.Portfolio{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedPortoflio, err
}

func (s *Service) TgFullPortfolioMessage(ctx context.Context, portfolioID int,
	tgChatID int) (string, error) {
	const op = "PortfolioService.TgFullPortfolioMessage"

	p, err := s.GetFullPortfolioByID(ctx, portfolioID, utils.GetCurrentUserID(ctx))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	message := fmt.Sprintf(`Портфель: %s - %d
		Акции: %d руб.
		Облигации: %d руб.
		Рубли: %d руб.
	`, p.Name, p.TotalCost, p.BondsCost, p.SharesCost, p.Cash)

	return message, nil
}
