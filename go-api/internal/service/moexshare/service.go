package moexshare

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/infrastracture/issclient"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) GetByTicker(ctx context.Context, ticker string) (entity.Share, error) {
	const op = "MoexShareService.GetByTicker"

	// Пробуем достать из нашей бд
	share, err := s.repo.GetByTicker(ctx, ticker)

	// Если её там нет то делаем запрос на МОЕХ и записываем в бд
	if errors.Is(err, database.ErrNotFound) {
		var err error
		share, err = s.createNewShareFromMoex(ctx, ticker)
		if err != nil {
			return entity.Share{}, fmt.Errorf("%s.createNewShareFromMoex, (ticker %s): %w", op, ticker, err)
		}
	} else if err != nil {
		return entity.Share{}, fmt.Errorf("%s: %w", op, err)
	}

	// если уже была в базе, то просто возвращаем
	return share, nil
}

func (s *Service) createNewShareFromMoex(ctx context.Context, ticker string) (entity.Share, error) {
	// если бумаги нет в БД то делаем запрос
	// на информацию по бумаге из апишки московской биржи
	secInfo, err := s.issClient.GetSecurityInfoByTicker(ctx, ticker)
	if err != nil {
		return entity.Share{}, err
	}

	share := entity.Share{
		SecurityCommonInfo: entity.SecurityCommonInfo{
			Name:      secInfo.Name,
			ShortName: secInfo.ShortName,
			Market:    secInfo.Market,
			Board:     secInfo.Board,
			Engine:    secInfo.Engine,
			Ticker:    ticker,
		},
	}

	if share.Market != entity.MoexMarketShares {
		return entity.Share{}, service.NewArgumentsError(fmt.Sprintf(
			"ticker %s не принадлежит рынку акций, рынок тикера  - %s", ticker, share.Market))
	}

	// Добираем доп инфу более подробный запросом на мск биржу. Пока что это только размер лота
	fullInfo, err := s.issClient.GetSecurityFullInfo(ctx, secInfo.Engine, secInfo.Market,
		secInfo.Board, ticker)
	if err != nil {
		return entity.Share{}, err
	}

	// Это всё что нам нужно было из фулинфо
	share.LotSize = fullInfo.LotSize
	share.PriceDecimals = fullInfo.PriceDecimals

	// сохраняем в бд
	share, err = s.repo.Insert(ctx, share)
	if err != nil {
		return entity.Share{}, err
	}

	return share, nil
}

type Repository interface {
	GetByTicker(ctx context.Context, ticker string) (entity.Share, error)
	Insert(ctx context.Context, share entity.Share) (entity.Share, error)
}
type Service struct {
	issClient *issclient.IssClient
	repo      Repository
}

func NewMoexShareService(repo Repository, issClient *issclient.IssClient) *Service {
	return &Service{
		issClient: issClient,
		repo:      repo,
	}
}
