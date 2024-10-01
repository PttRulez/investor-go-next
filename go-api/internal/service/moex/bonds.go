package moex

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) GetBondByTicker(ctx context.Context, ticker string) (domain.Bond, error) {
	const op = "MoexBondService.GetByTicker"

	// Пробуем достать из нашей бд
	bond, err := s.bondRepo.GetByTicker(ctx, ticker)

	// Если её там нет то делаем запрос на МОЕХ и записываем в бд
	if errors.Is(err, database.ErrNotFound) {
		var err error
		bond, err = s.createNewBondFromMoexISS(ctx, ticker)
		if err != nil {
			return domain.Bond{}, fmt.Errorf("%s.createNewBondFromMoex, (ticker %s): %w", op, ticker, err)
		}
	} else if err != nil {
		return domain.Bond{}, fmt.Errorf("%s: %w", op, err)
	}

	// если уже была в базе, то просто возвращаем
	return bond, nil
}

func (s *Service) createNewBondFromMoexISS(ctx context.Context, ticker string) (domain.Bond, error) {
	secInfo, err := s.issClient.GetSecurityInfoByTicker(ctx, ticker)
	bond := domain.Bond{
		SecurityCommonInfo: domain.SecurityCommonInfo{
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
		return domain.Bond{}, err
	}

	if bond.Market != domain.MoexMarketBonds {
		return domain.Bond{}, service.NewArgumentsError(fmt.Sprintf(
			"ticker %s не принадлежит рынку облигаций, рынок тикера  - %s", ticker, bond.Market))
	}

	// Добираем доп инфу более подробный запросом на мск биржу. Пока что это только размер лота
	fullInfo, err := s.issClient.GetSecurityFullInfo(ctx, secInfo.Engine, secInfo.Market,
		secInfo.Board, ticker)
	if err != nil {
		return domain.Bond{}, err
	}

	// Это всё что нам нужно было из фулинфо
	bond.LotSize = fullInfo.LotSize

	// сохраняем в бд
	b, err := s.bondRepo.Insert(ctx, bond)
	if err != nil {
		return domain.Bond{}, err
	}

	return b, nil
}
