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

func (s *MoexShareService) GetByTicker(ctx context.Context, ticker string) (*model.MoexShare, error) {
	share, err := s.repo.Moex.Share.GetByTicker(ctx, ticker)

	if errors.Is(err, sql.ErrNoRows) {
		// если бумаги нет в БД то делаем запрос
		// на информацию по бумаге из апишки московской биржи

		ISSecurityInfo, err := s.services.MoexApi.GetSecurityInfoBySecid(ticker)
		share := &model.MoexShare{
			Name:      ISSecurityInfo.Name,
			ShortName: ISSecurityInfo.ShortName,
			Ticker:    ticker,
			Market:    ISSecurityInfo.Market,
			Board:     ISSecurityInfo.Board,
			Engine:    ISSecurityInfo.Engine,
		}
		if err != nil {
			return nil, err
		}
		
		if share.Market != model.Market_Shares {
			e := httpresponse.NewErrSendToClient("Вы ввели некорректный Ticker акции", http.StatusBadRequest)
			return nil, e
		}

		// сохраняем в бд
		err = s.repo.Moex.Share.Insert(ctx, share)
		if err != nil {
			return nil, err
		}

		// ищем её же в бд
		share, err = s.repo.Moex.Share.GetByTicker(ctx, ticker)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// если уже была в базе, то просто возвращаем
	return share, nil
}

func (s *MoexShareService) UpdatePositionInDB(ctx context.Context, portfolioId int,
	securityId int) error {
	allDeals, err := s.repo.Deal.MoexShare.GetDealListByShareId(ctx, portfolioId, securityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareService.UpdatePositionInDB]: %w", err)
	}

	var position *model.Position
	oldPosition, err := s.repo.Position.MoexShare.Get(ctx, portfolioId, securityId)
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
		err = s.repo.Position.MoexShare.Insert(ctx, position)
		if err != nil {
			return err
		}
	} else {
		err = s.repo.Position.MoexShare.Update(ctx, position)
		if err != nil {
			return err
		}
	}

	return nil
}

type MoexShareService struct {
	repo     *repository.Repository
	services *Container
}

func NewMoexShareService(repo *repository.Repository, services *Container) *MoexShareService {
	return &MoexShareService{
		repo:     repo,
		services: services,
	}
}
