package portfolio

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
)

func (s *Service) CreateExpense(ctx context.Context, d domain.Expense,
	userID int) error {
	const op = "PortfolioService.CreateExpense"

	err := s.repo.InsertExpense(ctx, d, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteExpense(ctx context.Context, expenseID int, userID int) error {
	const op = "PortfolioService.DeleteExpense"

	err := s.repo.DeleteExpense(ctx, expenseID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
