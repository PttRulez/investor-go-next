package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateTransaction(ctx context.Context, t entity.Transaction) (entity.Transaction, error) {
	const op = "TransactionService.CreateTransaction"

	newTr, err := s.repo.Insert(ctx, t)
	if err != nil {
		return entity.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	return newTr, nil
}
func (s *Service) DeleteTransaction(ctx context.Context, transactionID int, userID int) error {
	const op = "TransactionService.DeleteTransaction"

	err := s.repo.Delete(ctx, transactionID, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrEntityNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

type PortfolioRepository interface {
	GetByID(ctx context.Context, portfolioID int, userID int) (entity.Portfolio, error)
}
type Repository interface {
	Delete(ctx context.Context, id int, userID int) error
	GetByID(ctx context.Context, id int, userID int) (entity.Transaction, error)
	Insert(ctx context.Context, t entity.Transaction) (entity.Transaction, error)
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
