package moex

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) GetShareByTicker(ctx context.Context, ticker string) (domain.Share, error) {
	const op = "MoexShareService.GetByTicker"

	// Пробуем достать из нашей бд
	share, err := s.repo.GetMoexShare(ctx, ticker)

	// Если её там нет то делаем запрос на МОЕХ и записываем в бд
	if errors.Is(err, storage.ErrNotFound) {
		var err error
		share, err = s.createNewShareFromMoexISS(ctx, ticker)
		if err != nil {
			return domain.Share{}, fmt.Errorf("%s.createNewShareFromMoex, (ticker %s): %w", op, ticker, err)
		}
	} else if err != nil {
		return domain.Share{}, fmt.Errorf("%s: %w", op, err)
	}

	// если уже была в базе, то просто возвращаем
	return share, nil
}

func (s *Service) createNewShareFromMoexISS(ctx context.Context, ticker string) (domain.Share, error) {
	// если бумаги нет в БД то делаем запрос
	// на информацию по бумаге из апишки московской биржи
	secInfo, err := s.issClient.GetSecurityInfoByTicker(ctx, ticker)
	if err != nil {
		return domain.Share{}, err
	}

	share := domain.Share{
		SecurityCommonInfo: domain.SecurityCommonInfo{
			Name:      secInfo.Name,
			ShortName: secInfo.ShortName,
			Market:    secInfo.Market,
			Board:     secInfo.Board,
			Engine:    secInfo.Engine,
			Ticker:    ticker,
		},
	}

	if share.Market != domain.MoexMarketShares {
		return domain.Share{}, service.NewArgumentsError(fmt.Sprintf(
			"ticker %s не принадлежит рынку акций, рынок тикера  - %s", ticker, share.Market))
	}

	// Добираем доп инфу более подробный запросом на мск биржу. Пока что это только размер лота
	fullInfo, err := s.issClient.GetSecurityFullInfo(ctx, secInfo.Engine, secInfo.Market,
		secInfo.Board, ticker)
	if err != nil {
		return domain.Share{}, err
	}

	// Это всё что нам нужно было из фулинфо
	share.LotSize = fullInfo.LotSize
	share.PriceDecimals = fullInfo.PriceDecimals

	// сохраняем в бд
	share, err = s.repo.InsertMoexShare(ctx, share)
	if err != nil {
		return domain.Share{}, err
	}

	return share, nil
}
