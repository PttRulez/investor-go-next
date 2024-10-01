package portfolio

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
)

func (s *Service) GetPositionList(ctx context.Context, userID int) ([]domain.Position, error) {
	const op = "PositionService.GetListByUserID"

	positions, err := s.positionRepo.GetListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return positions, nil
}

func (s *Service) AddInfoToPosition(ctx context.Context, i domain.PositionUpdateInfo) error {
	const op = "PositionService.UpdatePosition"

	err := s.positionRepo.AddInfo(ctx, i)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
