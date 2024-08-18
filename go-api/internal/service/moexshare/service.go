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

func (s *Service) GetBySecid(ctx context.Context, secID string) (*entity.Share, error) {
	const op = "MoexShareService.GetBySecid"

	// Пробуем достать из нашей бд
	share, err := s.repo.GetBySecid(ctx, secID)

	// Если её там нет то делаем запрос на МОЕХ и записываем в бд
	if errors.Is(err, database.ErrNotFound) {
		var err error
		share, err = s.createNewShareFromMoex(ctx, secID)
		if err != nil {
			return nil, fmt.Errorf("%s.createNewShareFromMoex, (secid %s): %w", op, secID, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// если уже была в базе, то просто возвращаем
	return share, nil
}

func (s *Service) createNewShareFromMoex(ctx context.Context, secID string) (*entity.Share, error) {
	// если бумаги нет в БД то делаем запрос
	// на информацию по бумаге из апишки московской биржи
	secInfo, err := s.issClient.GetSecurityInfoBySecid(ctx, secID)
	share := &entity.Share{
		SecurityCommonInfo: entity.SecurityCommonInfo{
			Name:      secInfo.Name,
			ShortName: secInfo.ShortName,
			Market:    secInfo.Market,
			Board:     secInfo.Board,
			Engine:    secInfo.Engine,
			Secid:     secID,
		},
	}
	if err != nil {
		return nil, err
	}

	if share.Market != entity.MoexMarketShares {
		return nil, service.NewArgumentsError(fmt.Sprintf(
			"secid %s не принадлежит рынку акций, рынок тикера  - %s", secID, share.Market))
	}

	// сохраняем в бд
	err = s.repo.Insert(ctx, share)
	if err != nil {
		return nil, err
	}

	return share, nil
}

type Repository interface {
	GetBySecid(ctx context.Context, ticker string) (*entity.Share, error)
	Insert(ctx context.Context, share *entity.Share) error
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
