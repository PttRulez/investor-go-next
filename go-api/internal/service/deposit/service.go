package deposit

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
	ierrors "github.com/pttrulez/investor-go/internal/errors"
	"github.com/pttrulez/investor-go/internal/utils"
)

func (s *Service) CreateDeposit(ctx context.Context, deposit *entity.Deposit) error {
	userId := utils.GetCurrentUserId(ctx)
	portfolio, err := s.portfolioRepo.GetById(ctx, deposit.PortfolioId)
	if err != nil {
		return err
	}
	if userId != portfolio.UserId {
		return ierrors.ErrNotYours
	}

	return s.depositRepo.Insert(ctx, deposit)
}

func (s *Service) DeleteDeposit(ctx context.Context, depositId int) error {
	userId := utils.GetCurrentUserId(ctx)
	cashout, err := s.depositRepo.GetById(ctx, depositId)
	if err != nil {
		return err
	}
	if userId != cashout.UserId {
		return ierrors.ErrNotYours
	}

	err = s.depositRepo.Delete(ctx, depositId)
	if err != nil {
		return err
	}
	return nil
}

type PortfolioRepository interface {
	GetById(ctx context.Context, portfolioId int) (*entity.Portfolio, error)
}

type Repository interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*entity.Deposit, error)
	Insert(ctx context.Context, c *entity.Deposit) error
}

type Service struct {
	portfolioRepo PortfolioRepository
	depositRepo   Repository
}

func NewDepositService(depositRepo Repository, portfolioRepo PortfolioRepository) *Service {
	return &Service{
		depositRepo:   depositRepo,
		portfolioRepo: portfolioRepo,
	}
}
