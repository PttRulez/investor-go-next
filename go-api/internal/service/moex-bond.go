package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (s *MoexBondService) GetByISIN(ctx context.Context, isin string) (*model.MoexBond, error) {
	bond, err := s.repo.Moex.Bond.GetByISIN(ctx, isin)
	if errors.Is(err, sql.ErrNoRows) {
		// если бумаги нет в БД то делаем запрос
		// на информацию по бумаге из апишки московской биржи
		ISSecurityInfo, err := s.services.Moex.Api.GetSecurityInfoBySecid(isin)
		bond := &model.MoexBond{
			Name:      ISSecurityInfo.Name,
			ShortName: ISSecurityInfo.ShortName,
			Isin:      isin,
			Market:    ISSecurityInfo.Market,
			Board:     ISSecurityInfo.Board,
			Engine:    ISSecurityInfo.Engine,
		}
		if err != nil {
			return nil, err
		}

		if bond.Market != model.Market_Bonds {
			e := httpresponse.NewErrSendToClient("Вы ввели некорректный ISIN облигации", http.StatusBadRequest)
			return nil, e
		}

		// сохраняем в бд
		err = s.repo.Moex.Bond.Insert(ctx, bond)
		if err != nil {
			return nil, err
		}

		// ищем её же в бд
		bond, err = s.repo.Moex.Bond.GetByISIN(ctx, isin)
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
	allDeals, err := s.repo.Deal.MoexBond.GetDealListByBondId(ctx, portfolioId, securityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondService.UpdatePositionInDB]: %w", err)
	}

	var position *model.Position
	oldPosition, err := s.repo.Position.MoexBond.Get(ctx, portfolioId, securityId)
	if err != nil {
		return err
	}
	if oldPosition != nil {
		position = oldPosition
	} else {
		position = &model.Position{}
		position.Exchange = model.EXCH_Moex
		position.PortfolioId = portfolioId
		position.SecurityId = securityId
	}

	// Calculate and add to position Amount
	var amount int
	var totalAmount int
	for _, deal := range allDeals {
		amount = deal.Amount
		if deal.Type == model.DT_Sell {
			amount = -amount
		}
		totalAmount += amount
	}
	position.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	left := position.Amount
	for _, deal := range allDeals {
		if deal.Type == model.DT_Sell {
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
		err = s.repo.Position.MoexBond.Insert(ctx, position)
		if err != nil {
			return err
		}
	} else {
		err = s.repo.Position.MoexBond.Update(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}

type MoexBondService struct {
	repo     *repository.Repository
	services *Container
}

func NewMoexBondService(repo *repository.Repository, services *Container) *MoexBondService {
	return &MoexBondService{
		repo: repo,
	}
}
