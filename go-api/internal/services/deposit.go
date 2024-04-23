package services

import (
	"context"

	"github.com/pttrulez/investor-go/internal/types"
)

type DepositService struct {
	repo     *types.Repository
	services *ServiceContainer
}

func NewDepositService(repo *types.Repository, services *ServiceContainer) *DepositService {
	return &DepositService{
		repo:     repo,
		services: services,
	}
}

func (s *DepositService) CreateDeposit(ctx context.Context, depositData *types.Deposit, userId int) error {
	portfolio, err := s.repo.Portfolio.GetById(ctx, depositData.PortfolioId)
	if err != nil {
		return err
	}

	if portfolio.UserId != userId {
		return types.ErrNotYours
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
		return types.ErrNotYours
	}

	err = s.repo.Deposit.Delete(ctx, depositId)
	if err != nil {
		return err
	}
	return nil
}
