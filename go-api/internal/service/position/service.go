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

func (s *Service) AddInfo(ctx context.Context, i entity.PositionUpdateInfo) error {
	const op = "PositionService.UpdatePosition"

	err := s.repo.AddInfo(ctx, i)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type Repository interface {
	GetListByUserID(ctx context.Context, userID int) ([]entity.Position, error)
	AddInfo(ctx context.Context, i entity.PositionUpdateInfo) error
}

type Service struct {
	repo Repository
}

func NewPositionService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
