package service

import (
	"context"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (s *CashoutService) CreateCashout(ctx context.Context, cashoutData *model.Cashout, userId int) error {
	portfolio, err := s.repo.Portfolio.GetById(ctx, cashoutData.PortfolioId)
	if err != nil {
		return err
	}

	if portfolio.UserId != userId {
		return httpresponse.ErrNotYours
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
		return httpresponse.ErrNotYours
	}

	err = s.repo.Cashout.Delete(ctx, cashoutId)
	if err != nil {
		return err
	}
	return nil
}

type CashoutService struct {
	repo     *repository.Repository
	services *Container
}

func NewCashoutService(repo *repository.Repository, services *Container) *CashoutService {
	return &CashoutService{
		repo:     repo,
		services: services,
	}
}
