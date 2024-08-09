package moex_bond

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pttrulez/investor-go/internal/entity"
	errors2 "github.com/pttrulez/investor-go/internal/errors"
	"github.com/pttrulez/investor-go/internal/infrastracture/iss_client"
	"net/http"
)

func (s *Service) GetBySecid(ctx context.Context, secid string) (*entity.Bond, error) {
	bond, err := s.repo.GetBySecid(ctx, secid)
	if errors.Is(err, sql.ErrNoRows) {
		// если бумаги нет в БД то делаем запрос
		// на информацию по бумаге из апишки московской биржи
		ISSecurityInfo, err := s.issClient.GetSecurityInfoBySecid(secid)
		bond = &entity.Bond{
			SecurityCommonInfo: entity.SecurityCommonInfo{
				Board:     ISSecurityInfo.Board,
				Engine:    ISSecurityInfo.Engine,
				Market:    ISSecurityInfo.Market,
				Name:      ISSecurityInfo.Name,
				ShortName: ISSecurityInfo.ShortName,
				Secid:     secid,
			},
			CouponPercent:   ISSecurityInfo.CouponPercent,
			CouponValue:     ISSecurityInfo.CouponValue,
			CouponFrequency: ISSecurityInfo.CouponFrequency,
			IssueDate:       ISSecurityInfo.IssueDate,
			FaceValue:       ISSecurityInfo.FaceValue,
			MatDate:         ISSecurityInfo.MatDate,
		}
		if err != nil {
			return nil, err
		}

		if bond.Market != entity.MoexMarketBonds {
			e := errors2.NewErrSendToClient(fmt.Sprintf("secid %s не принадлежит рынку облигаций, рынок тикера  - %s", secid, bond.Market),
				http.StatusBadRequest)
			return nil, e
		}

		// сохраняем в бд
		err = s.repo.Insert(ctx, bond)
		if err != nil {
			return nil, err
		}

		// ищем её же в бд
		bond, err = s.repo.GetBySecid(ctx, secid)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// если уже была в базе, то просто возвращаем
	return bond, nil
}

type Repository interface {
	GetBySecid(ctx context.Context, secid string) (*entity.Bond, error)
	Insert(ctx context.Context, share *entity.Bond) error
}
type Service struct {
	issClient *iss_client.IssClient
	repo      Repository
}

func NewMoexBondService(repo Repository, issClient *iss_client.IssClient) *Service {
	return &Service{
		issClient: issClient,
		repo:      repo,
	}
}
