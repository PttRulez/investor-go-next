package portfolio

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
)

func (s *Service) CreateDividend(ctx context.Context, d domain.Dividend,
	userID int) error {
	const op = "PortfolioService.AddDividend"

	err := s.repo.InsertDividend(ctx, d, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteDividend(ctx context.Context, dividendID int, userID int) error {
	const op = "PortfolioService.DeleteDividend"

	err := s.repo.DeleteDividend(ctx, dividendID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
