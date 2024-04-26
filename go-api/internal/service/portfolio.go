package service

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/lib/response"
	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

func (s *PortfolioService) GetPortfolio(ctx context.Context, portfolioId int,
	userId int) (*types.FullPortfolioData, error) {
	var fp types.FullPortfolioData

	done := make(chan struct{})
	errCh := make(chan error)

	// Проверяем владельца портфеля
	go func() {
		err := s.repo.Portfolio.GetByIdAndScan(ctx, portfolioId, &fp)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		if userId != fp.UserId {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", response.ErrNotYours)
			return
		}

		done <- struct{}{}
	}()

	// сделки по акциям
	go func() {
		shareDeals, err := s.repo.Deal.MoexShare.GetDealListByPortoflioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		fp.ShareDeals = shareDeals

		done <- struct{}{}
	}()

	// сделки по облигациям
	go func() {
		bondDeals, err := s.repo.Deal.MoexBond.GetDealListByPortoflioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		fp.BondDeals = bondDeals

		done <- struct{}{}
	}()

	// позиции облигаций
	go func() {
		// получили массив позиций по айди портфолио
		bondPositions, err := s.repo.Position.MoexBond.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		if len(bondPositions) == 0 {
			fp.BondPositions = []*types.BondPosition{}
			done <- struct{}{}
			return
		}

		// запрашиваем цены с московской биржи
		bondIsins := make([]string, len(bondPositions))
		for i, b := range bondPositions {
			bondIsins[i] = b.Isin
		}

		bondPrices, err := s.services.IssApi.GetStocksCurrentPrices(ctx, tmoex.Market_Bonds, bondIsins)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		bondPositionsMap := make(map[string]*types.BondPosition, len(bondPositions))
		for _, b := range bondPositions {
			bondPositionsMap[b.Isin] = b
		}
		for _, b := range bondPrices.Securities.Data {
			isin := b[0].(string)
			price := b[2].(float64) * 10
			amount := float64(bondPositionsMap[isin].Amount)

			bondPositionsMap[isin].CurrentPrice = price
			bondPositionsMap[isin].CurrentCost = int(price * amount)
		}
		for i, s := range bondPositions {
			bondPositions[i] = bondPositionsMap[s.Isin]
		}

		// записываем в FullPortfolio
		fp.BondPositions = bondPositions

		done <- struct{}{}
	}()

	// позиции акций
	go func() {
		// получили массив позиций по айди портфолио
		sharePositions, err := s.repo.Position.MoexShare.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		if len(sharePositions) == 0 {
			fp.SharePositions = []*types.SharePosition{}
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
			amount := float64(sharePositionsMap[s[0].(string)].Amount)

			sharePositionsMap[ticker].CurrentPrice = price
			sharePositionsMap[ticker].CurrentCost = int(price * amount)
		}
		for i, s := range sharePositions {
			sharePositions[i] = sharePositionsMap[s.Ticker]
		}

		// записываем в fp
		fp.SharePositions = sharePositions

		done <- struct{}{}
	}()

	// кэшауты
	go func() {
		cashouts, err := s.repo.Cashout.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		if cashouts == nil {
			cashouts = make([]*types.Cashout, 0)
		}
		fp.Cashouts = cashouts

		done <- struct{}{}
	}()

	// депозиты
	go func() {
		deposits, err := s.repo.Deposit.GetListByPortfolioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		if deposits == nil {
			deposits = make([]*types.Deposit, 0)
		}
		fp.Deposits = deposits

		done <- struct{}{}
	}()

	routinesWorking := 7

	for routinesWorking > 0 {
		select {
		case <-done:
			routinesWorking--
		case err := <-errCh:
			return nil, err
		}
	}

	for _, c := range fp.Cashouts {
		fp.CashoutsSum += c.Amount
	}
	for _, d := range fp.Deposits {
		fp.DepositsSum += d.Amount
	}

	spentToBuys := 0
	receivedFromSells := 0

	for _, d := range fp.ShareDeals {
		if d.Type == types.Buy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
	}
	for _, d := range fp.BondDeals {
		if d.Type == types.Buy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
	}

	for _, p := range fp.BondPositions {
		fp.TotalCost += p.CurrentCost
	}
	for _, p := range fp.SharePositions {
		fp.TotalCost += p.CurrentCost
	}

	fp.Cash = fp.DepositsSum - fp.CashoutsSum - spentToBuys + receivedFromSells
	fp.TotalCost += fp.Cash
	fp.Profitability = int((float64((fp.TotalCost+fp.CashoutsSum+fp.Cash))/float64(fp.DepositsSum) - 1) * 100)

	return &fp, nil
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
