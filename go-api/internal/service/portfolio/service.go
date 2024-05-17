package portfolio

import (
	"context"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"net/http"

	"github.com/pttrulez/investor-go/internal/types"
)

func (s *PortfolioService) GetFullPortfolioById(ctx context.Context, portfolioId int,
	userId int) (*entity.Portfolio, error) {
	test := entity.Bond{}
	var (
		done  = make(chan struct{})
		errCh = make(chan error)
		p     entity.Portfolio
	)

	// Проверяем владельца портфеля
	go func() {
		portfolioFronmDb, err := s.GetPortfolioById(ctx, portfolioId, userId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		p.Id = portfolioFronmDb.Id
		p.Compound = portfolioFronmDb.Compound
		p.Name = portfolioFronmDb.Name

		done <- struct{}{}
	}()

	// сделки по акциям
	go func() {
		shareDeals, err := s.repo.Deal.MoexShare.GetDealListByPortoflioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		p.ShareDeals = shareDeals

		done <- struct{}{}
	}()

	// сделки по облигациям
	go func() {
		bondDeals, err := s.repo.Deal.MoexBond.GetDealListByPortoflioId(ctx, portfolioId)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}
		p.BondDeals = bondDeals

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
			p.BondPositions = []*entity.Position{}
			done <- struct{}{}
			return
		}

		// запрашиваем цены с московской биржи
		bondIsins := make([]string, len(bondPositions))
		for i, b := range bondPositions {
			bondIsins[i] = b.Secid
		}

		bondPrices, err := s.services.Moex.Api.GetStocksCurrentPrices(ctx, iss_client.Market_Bonds, bondIsins)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		bondPositionsMap := make(map[string]*entity.Position, len(bondPositions))
		for _, b := range bondPositions {
			bondPositionsMap[b.Secid] = b
		}
		for _, b := range bondPrices.Securities.Data {
			isin := b[0].(string)
			price := b[2].(float64) * 10
			amount := float64(bondPositionsMap[isin].Amount)

			bondPositionsMap[isin].CurrentPrice = price
			bondPositionsMap[isin].CurrentCost = int(price * amount)
		}
		for i, s := range bondPositions {
			bondPositions[i] = bondPositionsMap[s.Secid]
		}

		// записываем в FullPortfolio
		p.BondPositions = bondPositions

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
			p.SharePositions = []*entity.Position{}
			done <- struct{}{}
			return
		}

		// запрашиваем цены с московской биржи
		shareTickers := make([]string, len(sharePositions))
		for i, s := range sharePositions {
			shareTickers[i] = s.Secid
		}

		sharePrices, err := s.services.Moex.Api.GetStocksCurrentPrices(ctx, iss_client.Market_Shares, shareTickers)
		if err != nil {
			errCh <- fmt.Errorf("<-[PortfolioService.GetPortfolio]: \n%w", err)
			return
		}

		// каждой позиции присваиваем текущую цену и общую текущую стоимость
		sharePositionsMap := make(map[string]*entity.Position, len(sharePositions))
		for _, s := range sharePositions {
			sharePositionsMap[s.Secid] = s
		}
		for _, s := range sharePrices.Securities.Data {
			ticker := s[0].(string)
			price := s[2].(float64)
			amount := float64(sharePositionsMap[s[0].(string)].Amount)

			sharePositionsMap[ticker].CurrentPrice = price
			sharePositionsMap[ticker].CurrentCost = int(price * amount)
		}
		for i, s := range sharePositions {
			sharePositions[i] = sharePositionsMap[s.Secid]
		}

		// записываем в p
		p.SharePositions = sharePositions

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
			cashouts = make([]*entity.Cashout, 0)
		}
		p.Cashouts = cashouts

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
			deposits = make([]*entity.Deposit, 0)
		}
		p.Deposits = deposits

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

	for _, c := range p.Cashouts {
		p.CashoutsSum += c.Amount
	}
	for _, d := range p.Deposits {
		p.DepositsSum += d.Amount
	}

	spentToBuys := 0
	receivedFromSells := 0

	for _, d := range p.ShareDeals {
		if d.Type == types.DT_Buy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
	}
	for _, d := range p.BondDeals {
		if d.Type == types.DT_Buy {
			spentToBuys += d.Amount * int(d.Price)
		} else {
			receivedFromSells += d.Amount * int(d.Price)
		}
	}

	for _, pos := range p.BondPositions {
		p.TotalCost += pos.CurrentCost
	}
	for _, pos := range p.SharePositions {
		p.TotalCost += pos.CurrentCost
	}

	p.Cash = p.DepositsSum - p.CashoutsSum - spentToBuys + receivedFromSells
	p.TotalCost += p.Cash
	p.Profitability = int((float64((p.TotalCost+p.CashoutsSum+p.Cash))/float64(p.DepositsSum) - 1) * 100)

	return &p, nil
}

func (s *PortfolioService) CreatePortfolio(ctx context.Context, p *entity.Portfolio) error {
	return s.repo.Portfolio.Insert(ctx, p)
}

func (s *PortfolioService) GetListByUserId(ctx context.Context, userId int) ([]*entity.Portfolio, error) {
	portfolios, err := s.repo.Portfolio.GetListByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (s *PortfolioService) GetPortfolioById(ctx context.Context, portfolioId int,
	userId int) (*entity.Portfolio, error) {
	portfolio, err := s.repo.Portfolio.GetById(ctx, portfolioId)
	if err != nil {
		return nil, err
	}
	if portfolio.UserId != userId {
		return nil, httpresponse.NewErrSendToClient("Такого портфолио не существует", http.StatusNotFound)
	}
	return portfolio, nil
}

func (s *PortfolioService) DeletePortfolio(ctx context.Context, portfolioId int, userId int) error {
	if _, err := s.GetPortfolioById(ctx, portfolioId, userId); err != nil {
		return err
	}
	return s.repo.Portfolio.Delete(ctx, portfolioId)
}

func (s *PortfolioService) UpdatePortfolio(ctx context.Context, portfolio *entity.Portfolio, userId int) error {
	if _, err := s.GetPortfolioById(ctx, portfolio.Id, userId); err != nil {
		return err
	}
	return s.cashoutRepo.Update(ctx, portfolio)
}

type PortfolioService struct {
	cashoutRepo CashoutRepository
}

type CashoutRepository interface {
	GetListByPortfolioId(ctx context.Context, id int) ([]*entity.Cashout, error)
}

func NewPortfolioService(cashoutRepo CashoutRepository) *PortfolioService {
	return &PortfolioService{
		cashoutRepo: cashoutRepo,
	}
}
