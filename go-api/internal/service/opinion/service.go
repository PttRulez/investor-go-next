package opinion

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) AttachToPosition(ctx context.Context, opinionID, positionID int) error {
	const op = "OpinionService.AttachToPosition"

	err := s.repo.AttachToPosition(ctx, opinionID, positionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) CreateOpinion(ctx context.Context, o entity.Opinion) (entity.Opinion, error) {
	const op = "OpinionService.CreateOpinion"

	e, err := s.repo.Insert(ctx, o)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (s *Service) DeleteOpinionByID(ctx context.Context, id int, userID int) error {
	const op = "OpinionService.DeleteOpinion"

	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrEntityNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetOpinions(ctx context.Context, f entity.OpinionFilters, userID int) (
	[]entity.Opinion, error) {
	const op = "OpinionService.GetOpinions"

	o, err := s.repo.GetOpinionsList(ctx, f, userID)
	if err != nil {
		return o, fmt.Errorf("%s: %w", op, err)
	}

	return o, nil
}

type Repository interface {
	AttachToPosition(ctx context.Context, opinionID, positionID int) error
	Delete(ctx context.Context, id int, userID int) error
	GetOpinionsList(ctx context.Context, f entity.OpinionFilters, userID int) ([]entity.Opinion,
		error)
	Insert(ctx context.Context, o entity.Opinion) (entity.Opinion, error)
}

type Service struct {
	repo Repository
}

func NewOpinionService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
