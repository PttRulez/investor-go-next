package services

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

func (s *PortfolioService) GetPortfolio(ctx context.Context, portfolioId int, userId int) (*types.FullPortfolioData, error) {
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
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", types.ErrNotYours)
			return
		}

		done <- struct{}{}
	}()

	// Получаем сделки по акциям МОЕХ по портфелю из БД и записываем в fp
	go func() {
		shareDeals, err := s.repo.Deal.MoexShare.GetDealsListByPortoflioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		fp.ShareDeals = shareDeals

		done <- struct{}{}
	}()

	// Получаем позиции облигаций по портфелю из БД, запрашиваем текущие цены для низ и записываем в FullPortfolio
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
		fmt.Println("bondIsins:", bondIsins)
		bondPrices, err := s.services.IssApi.GetStocksCurrentPrices(ctx, tmoex.Market_Shares, bondIsins)
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
			ticker := b[0].(string)
			price := b[2].(float64)
			amount := float64(bondPositionsMap[ticker].Amount)

			bondPositionsMap[ticker].CurrentPrice = price
			bondPositionsMap[ticker].CurrentCost = price * amount
		}
		for i, s := range bondPositions {
			bondPositions[i] = bondPositionsMap[s.Isin]
		}

		// записываем в FullPortfolio
		fp.BondPositions = bondPositions

		done <- struct{}{}
	}()

	// Получаем позиции акций по портфелю из БД, запрашиваем текущие цены для низ и записываем в fp
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
			fmt.Printf("%v %v %v\n", ticker, s[1], price)
			amount := float64(sharePositionsMap[s[0].(string)].Amount)

			sharePositionsMap[ticker].CurrentPrice = price
			sharePositionsMap[ticker].CurrentCost = price * amount
		}
		for i, s := range sharePositions {
			sharePositions[i] = sharePositionsMap[s.Ticker]
		}

		// записываем в fp
		fp.SharePositions = sharePositions

		done <- struct{}{}
	}()

	// Получаем кэшауты по портфелю из БД и записываем в fp
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

	// Получаем депозиты по портфелю из БД и записываем в fp
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

	routinesWorking := 6

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
			spentToBuys += d.Amount
		} else {
			receivedFromSells += d.Amount
		}
	}
	for _, d := range fp.BondDeals {
		if d.Type == types.Buy {
			spentToBuys += d.Amount
		} else {
			receivedFromSells += d.Amount
		}
	}

	fp.Cash = fp.DepositsSum - fp.CashoutsSum - spentToBuys + receivedFromSells

	fp.Positions = make([]*types.Position, len(fp.SharePositions)+len(fp.BondPositions))

	for i, p := range fp.SharePositions {
		fp.Positions[i] = &p.Position
		fp.Positions[i].SecurityType = types.ST_Share
		fp.TotalCost += int(p.CurrentCost)
	}
	for i, p := range fp.BondPositions {
		fp.Positions[i] = &p.Position
		fp.Positions[i].SecurityType = types.ST_Bond
		fp.TotalCost += int(p.CurrentCost)
	}

	fp.Profitability = float64((fp.TotalCost + fp.CashoutsSum + fp.Cash) / fp.DepositsSum)

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
