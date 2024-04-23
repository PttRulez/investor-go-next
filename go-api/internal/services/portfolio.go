package services

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

func (s *PortfolioService) GetPortfolio(ctx context.Context, portfolioId int, userId int) (*types.FullPortfolioData, error) {
	var fullPortfolio types.FullPortfolioData

	done := make(chan struct{})
	errCh := make(chan error)

	// Проверяем владельца портфеля
	go func() {
		err := s.repo.Portfolio.GetByIdAndScan(ctx, portfolioId, &fullPortfolio)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		if userId != fullPortfolio.UserId {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", types.ErrNotYours)
			return
		}

		done <- struct{}{}
	}()

	// Получаем сделки по акциям МОЕХ по портфелю из БД и записываем в fullPortfolio
	go func() {
		shareDeals, err := s.repo.Deal.MoexShares.GetDealsListByPortoflioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		fullPortfolio.ShareDeals = shareDeals

		done <- struct{}{}
	}()

	// Получаем позиции облигаций по портфелю из БД, запрашиваем текущие цены для низ и записываем в fullPortfolio
	go func() {
		// получили массив позиций по айди портфолио
		bondPositions, err := s.repo.MoexBondPosition.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		if len(bondPositions) == 0 {
			fullPortfolio.BondPositions = []*types.BondPosition{}
			done <- struct{}{}
			return
		}

		// запрашиваем цены с московской биржи
		bondTickers := make([]string, len(bondPositions))
		for i, b := range bondPositions {
			bondTickers[i] = b.Ticker
		}
		fmt.Println("bondTickers:", bondTickers)
		bondPrices, err := s.services.IssApi.GetStocksCurrentPrices(ctx, tmoex.Market_Shares, bondTickers)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		bondPositionsMap := make(map[string]*types.BondPosition, len(bondPositions))
		for _, b := range bondPositions {
			bondPositionsMap[b.Ticker] = b
		}
		for _, b := range bondPrices.Securities.Data {
			ticker := b[0].(string)
			price := b[2].(float64)
			amount := float64(bondPositionsMap[ticker].Amount)

			bondPositionsMap[ticker].CurrentPrice = price
			bondPositionsMap[ticker].CurrentCost = price * amount
		}
		for i, s := range bondPositions {
			bondPositions[i] = bondPositionsMap[s.Ticker]
		}

		// записываем в fullPortfolio
		fullPortfolio.BondPositions = bondPositions

		done <- struct{}{}
	}()

	// Получаем позиции акций по портфелю из БД, запрашиваем текущие цены для низ и записываем в fullPortfolio
	go func() {
		// получили массив позиций по айди портфолио
		sharePositions, err := s.repo.MoexSharePosition.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		if len(sharePositions) == 0 {
			fullPortfolio.SharePositions = []*types.SharePosition{}
			done <- struct{}{}
			return
		}

		// запрашиваем цены с московской биржи
		shareTickers := make([]string, len(sharePositions))
		for i, s := range sharePositions {
			shareTickers[i] = s.Ticker
		}

		sharePrices, err := s.services.IssApi.GetStocksCurrentPrices(ctx, tmoex.Market_Shares, shareTickers)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		sharePositionsMap := make(map[string]*types.SharePosition, len(sharePositions))
		for _, s := range sharePositions {
			sharePositionsMap[s.Ticker] = s
		}
		for _, s := range sharePrices.Securities.Data {
			ticker := s[0].(string)
			price := s[2].(float64)
			fmt.Printf("%v %v %v\n", ticker, s[1], price)
			amount := float64(sharePositionsMap[s[0].(string)].Amount)

			sharePositionsMap[ticker].CurrentPrice = price
			sharePositionsMap[ticker].CurrentCost = price * amount
		}
		for i, s := range sharePositions {
			sharePositions[i] = sharePositionsMap[s.Ticker]
		}

		// записываем в fullPortfolio
		fullPortfolio.SharePositions = sharePositions

		done <- struct{}{}
	}()

	// Получаем кэшауты по портфелю из БД и записываем в fullPortfolio
	go func() {
		cashouts, err := s.repo.Cashout.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		if cashouts == nil {
			cashouts = make([]*types.Cashout, 0)
		}
		fullPortfolio.Cashouts = cashouts

		done <- struct{}{}
	}()

	// Получаем депозиты по портфелю из БД и записываем в fullPortfolio
	go func() {
		deposits, err := s.repo.Deposit.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		if deposits == nil {
			deposits = make([]*types.Deposit, 0)
		}
		fullPortfolio.Deposits = deposits

		done <- struct{}{}
	}()

	routinesWorking := 6

	for routinesWorking > 0 {
		select {
		case <-done:
			routinesWorking--
		case err := <-errCh:
			return nil, err
		}
	}

	for _, c := range fullPortfolio.Cashouts {
		fullPortfolio.CashoutsSum += c.Amount
	}
	for _, d := range fullPortfolio.Deposits {
		fullPortfolio.DepositsSum += d.Amount
	}

	spentToBuys := 0
	receivedFromSells := 0

	for _, d := range fullPortfolio.ShareDeals {
		if d.Type == types.Buy {
			spentToBuys += d.Amount
		} else {
			receivedFromSells += d.Amount
		}
	}
	for _, d := range fullPortfolio.BondDeals {
		if d.Type == types.Buy {
			spentToBuys += d.Amount
		} else {
			receivedFromSells += d.Amount
		}
	}

	fullPortfolio.Cash = fullPortfolio.DepositsSum - fullPortfolio.CashoutsSum - spentToBuys + receivedFromSells

	for _, p := range fullPortfolio.BondPositions {
		fullPortfolio.TotalCost += p.Amount
	}
	for _, p := range fullPortfolio.SharePositions {
		fullPortfolio.TotalCost += p.Amount
	}

	return &fullPortfolio, nil
}

type PortfolioService struct {
	repo     *types.Repository
	services *ServiceContainer
}

func NewPortfolioService(repo *types.Repository, services *ServiceContainer) *PortfolioService {
	return &PortfolioService{
		repo:     repo,
		services: services,
	}
}
