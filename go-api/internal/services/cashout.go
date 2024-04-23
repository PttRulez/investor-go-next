package services

import (
	"context"

	"github.com/pttrulez/investor-go/internal/types"
)

func (s *CashoutService) CreateCashout(ctx context.Context, cashoutData *types.Cashout, userId int) error {
	portfolio, err := s.repo.Portfolio.GetById(ctx, cashoutData.PortfolioId)
	if err != nil {
		return err
	}

	if portfolio.UserId != userId {
		return types.ErrNotYours
	}
	return s.repo.Cashout.Insert(ctx, cashoutData)
}
func (s *CashoutService) DeleteCashout(ctx context.Context, cashoutId int, userId int) error {
	cashout, err := s.repo.Cashout.GetById(ctx, cashoutId)
	if err != nil {
		return err
	}

	portoflio, err := s.repo.Portfolio.GetById(ctx, cashout.PortfolioId)
	if err != nil {
		return err
	}
	if portoflio.UserId != userId {
		return types.ErrNotYours
	}

	err = s.repo.Cashout.Delete(ctx, cashoutId)
	if err != nil {
		return err
	}
	return nil
}

type CashoutService struct {
	repo     *types.Repository
	services *ServiceContainer
}

func NewCashoutService(repo *types.Repository, services *ServiceContainer) *CashoutService {
	return &CashoutService{
		repo:     repo,
		services: services,
	}
}
