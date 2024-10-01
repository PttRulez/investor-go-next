package opinion

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) AttachToPosition(ctx context.Context, opinionID, positionID int) error {
	const op = "OpinionService.AttachToPosition"

	err := s.opinionRepo.AttachToPosition(ctx, opinionID, positionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) CreateOpinion(ctx context.Context, o domain.Opinion) (domain.Opinion, error) {
	const op = "OpinionService.CreateOpinion"

	e, err := s.opinionRepo.Insert(ctx, o)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (s *Service) DeleteOpinionByID(ctx context.Context, id int, userID int) error {
	const op = "OpinionService.DeleteOpinion"

	err := s.opinionRepo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetOpinionsList(ctx context.Context, f domain.OpinionFilters, userID int) (
	[]domain.Opinion, error) {
	const op = "OpinionService.GetOpinions"

	o, err := s.opinionRepo.GetOpinionsList(ctx, f, userID)
	if err != nil {
		return o, fmt.Errorf("%s: %w", op, err)
	}

	return o, nil
}
