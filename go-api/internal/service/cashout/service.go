package cashout

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
	ierrors "github.com/pttrulez/investor-go/internal/errors"
	"github.com/pttrulez/investor-go/internal/utils"
)

func (s *Service) CreateCashout(ctx context.Context, cashout *entity.Cashout) error {
	userId := utils.GetCurrentUserId(ctx)
	portfolio, err := s.portfolioRepo.GetById(ctx, cashout.PortfolioId, utils.GetCurrentUserId(ctx))
	if err != nil {
		return err
	}
	if userId != portfolio.UserId {
		return ierrors.ErrNotYours
	}

	return s.cashoutRepo.Insert(ctx, cashout)
}
func (s *Service) DeleteCashout(ctx context.Context, cashoutId int) error {
	userId := utils.GetCurrentUserId(ctx)
	cashout, err := s.cashoutRepo.GetById(ctx, cashoutId, utils.GetCurrentUserId(ctx))
	if err != nil {
		return err
	}
	if userId != cashout.UserId {
		return ierrors.ErrNotYours
	}

	err = s.cashoutRepo.Delete(ctx, cashoutId, utils.GetCurrentUserId(ctx))
	if err != nil {
		return err
	}
	return nil
}

type PortfolioRepository interface {
	GetById(ctx context.Context, portfolioId int, userId int) (*entity.Portfolio, error)
}

type Repository interface {
	Delete(ctx context.Context, id int, userId int) error
	GetById(ctx context.Context, id int, userId int) (*entity.Cashout, error)
	Insert(ctx context.Context, c *entity.Cashout) error
}

type Service struct {
	portfolioRepo PortfolioRepository
	cashoutRepo   Repository
}

func NewCashoutService(cashoutRepo Repository, portfolioRepo PortfolioRepository) *Service {
	return &Service{
		cashoutRepo:   cashoutRepo,
		portfolioRepo: portfolioRepo,
	}
}
