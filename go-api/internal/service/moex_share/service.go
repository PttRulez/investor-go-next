package moex_share

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	ierrors "github.com/pttrulez/investor-go/internal/errors"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"net/http"
)

func (s *Service) GetBySecid(ctx context.Context, secid string) (*entity.Share, error) {
	share, err := s.repo.GetBySecid(ctx, secid)

	if errors.Is(err, sql.ErrNoRows) {
		// если бумаги нет в БД то делаем запрос
		// на информацию по бумаге из апишки московской биржи

		ISSecurityInfo, err := s.issClient.GetSecurityInfoBySecid(secid)
		share = &entity.Share{
			SecurityCommonInfo: entity.SecurityCommonInfo{
				Name:      ISSecurityInfo.Name,
				ShortName: ISSecurityInfo.ShortName,
				Market:    ISSecurityInfo.Market,
				Board:     ISSecurityInfo.Board,
				Engine:    ISSecurityInfo.Engine,
				Secid:     secid,
			},
		}
		if err != nil {
			return nil, err
		}

		if share.Market != entity.MoexMarketShares {
			e := ierrors.NewErrSendToClient("Вы ввели некорректный Ticker акции",
				http.StatusBadRequest)
			return nil, e
		}

		// сохраняем в бд
		err = s.repo.Insert(ctx, share)
		if err != nil {
			return nil, err
		}

		// ищем её же в бд
		share, err = s.repo.GetBySecid(ctx, secid)
		if err != nil {
			return nil, err
		}

		if share.Market != entity.MoexMarketShares {
			return nil, errors.New(fmt.Sprintf("secid %s не принадлежит рынку акций, рынок тикера  - %s", secid, share.Market))
		}

	} else if err != nil {
		return nil, err
	}

	// если уже была в базе, то просто возвращаем
	return share, nil
}

type Repository interface {
	GetBySecid(ctx context.Context, ticker string) (*entity.Share, error)
	Insert(ctx context.Context, share *entity.Share) error
}
type Service struct {
	issClient *iss_client.IssClient
	repo      Repository
}

func NewMoexShareService(repo Repository, issClient *iss_client.IssClient) *Service {
	return &Service{
		issClient: issClient,
		repo:      repo,
	}
}
