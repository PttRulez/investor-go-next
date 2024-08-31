package moexbond

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/infrastracture/issclient"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) GetByTicker(ctx context.Context, ticker string) (entity.Bond, error) {
	const op = "MoexBondService.GetByTicker"

	// Пробуем достать из нашей бд
	bond, err := s.repo.GetByTicker(ctx, ticker)

	// Если её там нет то делаем запрос на МОЕХ и записываем в бд
	if errors.Is(err, database.ErrNotFound) {
		var err error
		bond, err = s.createNewBondFromMoex(ctx, ticker)
		if err != nil {
			return entity.Bond{}, fmt.Errorf("%s.createNewBondFromMoex, (ticker %s): %w", op, ticker, err)
		}
	} else if err != nil {
		return entity.Bond{}, fmt.Errorf("%s: %w", op, err)
	}

	// если уже была в базе, то просто возвращаем
	return bond, nil
}

func (s *Service) createNewBondFromMoex(ctx context.Context, ticker string) (entity.Bond, error) {
	// если бумаги нет в БД то делаем запрос
	// на информацию по бумаге из апишки московской биржи
	secInfo, err := s.issClient.GetSecurityInfoByTicker(ctx, ticker)
	bond := entity.Bond{
		SecurityCommonInfo: entity.SecurityCommonInfo{
			Board:     secInfo.Board,
			Engine:    secInfo.Engine,
			Market:    secInfo.Market,
			Name:      secInfo.Name,
			ShortName: secInfo.ShortName,
			Ticker:    ticker,
		},
		CouponPercent:   secInfo.CouponPercent,
		CouponValue:     secInfo.CouponValue,
		CouponFrequency: secInfo.CouponFrequency,
		IssueDate:       secInfo.IssueDate,
		FaceValue:       secInfo.FaceValue,
		MatDate:         secInfo.MatDate,
	}
	if err != nil {
		return entity.Bond{}, err
	}

	if bond.Market != entity.MoexMarketBonds {
		return entity.Bond{}, service.NewArgumentsError(fmt.Sprintf(
			"ticker %s не принадлежит рынку облигаций, рынок тикера  - %s", ticker, bond.Market))
	}

	// Добираем доп инфу более подробный запросом на мск биржу. Пока что это только размер лота
	fullInfo, err := s.issClient.GetSecurityFullInfo(ctx, secInfo.Engine, secInfo.Market,
		secInfo.Board, ticker)
	if err != nil {
		return entity.Bond{}, err
	}

	// Это всё что нам нужно было из фулинфо
	bond.LotSize = fullInfo.LotSize

	// сохраняем в бд
	b, err := s.repo.Insert(ctx, bond)
	if err != nil {
		return entity.Bond{}, err
	}

	return b, nil
}

type Repository interface {
	GetByTicker(ctx context.Context, ticker string) (entity.Bond, error)
	Insert(ctx context.Context, bond entity.Bond) (entity.Bond, error)
}
type Service struct {
	issClient *issclient.IssClient
	repo      Repository
}

func NewMoexBondService(repo Repository, issClient *issclient.IssClient) *Service {
	return &Service{
		issClient: issClient,
		repo:      repo,
	}
}
