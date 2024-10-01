package portfolio

import (
	"context"

	"github.com/pttrulez/investor-go/internal/domain"
)

func (s *Service) AddDividendsPaid(ctx context.Context, portfolio domain.Portfolio,
	userID int) (domain.Portfolio, error) {
	return domain.Portfolio{}, nil
}
