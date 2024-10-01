package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error) {
	const op = "TransactionService.CreateTransaction"

	newTr, err := s.transactionRepo.Insert(ctx, t)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	return newTr, nil
}
func (s *Service) DeleteTransaction(ctx context.Context, transactionID int, userID int) error {
	const op = "TransactionService.DeleteTransaction"

	err := s.transactionRepo.Delete(ctx, transactionID, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}
