package transaction

import (
	"context"

	"github.com/pttrulez/investor-go/internal/entity"
)

func (s *Service) CreateTransaction(ctx context.Context, t *entity.Transaction) error {
	return s.repo.Insert(ctx, t)
}
func (s *Service) DeleteTransaction(ctx context.Context, transactionID int, userID int) error {
	err := s.repo.Delete(ctx, transactionID, userID)
	if err != nil {
		return err
	}
	return nil
}

type PortfolioRepository interface {
	GetByID(ctx context.Context, portfolioID int, userID int) (*entity.Portfolio, error)
}
type Repository interface {
	Delete(ctx context.Context, id int, userID int) error
	GetByID(ctx context.Context, id int, userID int) (*entity.Transaction, error)
	Insert(ctx context.Context, t *entity.Transaction) error
}
type Service struct {
	portfolioRepo PortfolioRepository
	repo          Repository
}

func NewTransactionService(transactionRepo Repository, portfolioRepo PortfolioRepository) *Service {
	return &Service{
		repo:          transactionRepo,
		portfolioRepo: portfolioRepo,
	}
}
