package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/pttrulez/investor-go/internal/types"
	tmoex "github.com/pttrulez/investor-go/internal/types/moex"
)

type MoexBondService struct {
	moexApi *IssApiService
	repo    *types.Repository
}

func (s *MoexBondService) GetByTicker(ctx context.Context, ticker string) (*tmoex.Bond, error) {
	bond, err := s.repo.Moex.Bonds.GetByTicker(ctx, ticker)
	if errors.Is(err, sql.ErrNoRows) {
		// если бумаги нет в БД то делаем запрос
		// на информацию по бумаге из апишки московской биржи
		security, err := s.moexApi.GetSecurityInfoByTicker(ticker)
		if err != nil {
			return nil, err
		}

		// сохраняем в бд
		err = s.repo.Moex.Bonds.Insert(ctx, &tmoex.Bond{Security: *security})
		if err != nil {
			return nil, err
		}

		// ищем её же в бд
		bond, err = s.repo.Moex.Bonds.GetByTicker(ctx, ticker)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// если уже была в базе, то просто возвращаем
	return bond, nil
}

func (s *MoexBondService) UpdatePositionInDB(ctx context.Context, portfolioId int,
	securityId int) error {
	allDeals, err := s.repo.Deal.MoexBonds.GetDealsListForSecurity(ctx, portfolioId, securityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondService.UpdatePositionInDB]: %w", err)
	}

	var position *types.BondPosition
	oldPosition, err := s.repo.MoexBondPosition.Get(ctx, portfolioId, securityId)
	if err != nil {
		return err
	}
	if oldPosition != nil {
		position = oldPosition
	} else {
		position = &types.BondPosition{}
		position.Exchange = types.EXCH_Moex
		position.PortfolioId = portfolioId
		position.SecurityId = securityId
	}

	// Calculate and add to position Amount
	var amount int
	var totalAmount int
	for _, deal := range allDeals {
		amount = deal.Amount
		if deal.Type == types.Sell {
			amount = -amount
		}
		totalAmount += amount
	}
	position.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	left := position.Amount
	for _, deal := range allDeals {
		if deal.Type == types.Sell {
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
		err = s.repo.MoexBondPosition.Insert(ctx, position)
		if err != nil {
			return err
		}
	} else {
		err = s.repo.MoexBondPosition.Update(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewBondService(repo *types.Repository) *MoexBondService {
	return &MoexBondService{
		moexApi: NewIssApiService(),
		repo:    repo,
	}
}
