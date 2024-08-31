package position

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
)

func (s *Service) GetListByUserID(ctx context.Context, userID int) ([]entity.Position, error) {
	const op = "PositionService.GetListByUserID"

	positions, err := s.repo.GetListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return positions, nil
}

type Repository interface {
	GetListByUserID(ctx context.Context, userID int) ([]entity.Position, error)
}

type Service struct {
	repo Repository
}

func NewPositionService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
