package service

import (
	"context"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

type DepositService struct {
	repo     *repository.Repository
	services *Container
}

func NewDepositService(repo *repository.Repository, services *Container) *DepositService {
	return &DepositService{
		repo:     repo,
		services: services,
	}
}

func (s *DepositService) CreateDeposit(ctx context.Context, depositData *model.Deposit, userId int) error {
	portfolio, err := s.repo.Portfolio.GetById(ctx, depositData.PortfolioId)
	if err != nil {
		return err
	}

	if portfolio.UserId != userId {
		return httpresponse.ErrNotYours
	}
	return s.repo.Deposit.Insert(ctx, depositData)
}

func (s *DepositService) DeleteDeposit(ctx context.Context, depositId int, userId int) error {
	deposit, err := s.repo.Deposit.GetById(ctx, depositId)
	if err != nil {
		return err
	}

	portoflio, err := s.repo.Portfolio.GetById(ctx, deposit.PortfolioId)
	if err != nil {
		return err
	}
	if portoflio.UserId != userId {
		return httpresponse.ErrNotYours
	}

	err = s.repo.Deposit.Delete(ctx, depositId)
	if err != nil {
		return err
	}
	return nil
}
