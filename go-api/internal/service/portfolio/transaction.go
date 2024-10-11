package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go-next/go-api/internal/service"
)

func (s *Service) CreateTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error) {
	const op = "TransactionService.CreateTransaction"

	newTr, err := s.repo.InsertTransaction(ctx, t)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	return newTr, nil
}
func (s *Service) DeleteTransaction(ctx context.Context, transactionID int, userID int) error {
	const op = "TransactionService.DeleteTransaction"

	err := s.repo.DeleteTransaction(ctx, transactionID, userID)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}
