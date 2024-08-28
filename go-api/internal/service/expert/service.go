package expert

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/infrastracture/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateNewExpert(ctx context.Context, expert entity.Expert) (entity.Expert, error) {
	const op = "ExpertService.CreateNewExpert"

	e, err := s.repo.Insert(ctx, expert)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}
func (s *Service) DeleteExpert(ctx context.Context, id int, userID int) error {
	const op = "ExpertService.DeleteExpert"

	err := s.repo.Delete(ctx, id, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrEntityNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}
func (s *Service) GetListByUserID(ctx context.Context, userID int) ([]entity.Expert, error) {
	const op = "ExpertService.GetListByUserID"

	e, err := s.repo.GetListByUserID(ctx, userID)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}
func (s *Service) GetExpertByID(ctx context.Context, id int, userID int) (entity.Expert, error) {
	const op = "ExpertService.GetExpertByID"

	expert, err := s.repo.GetOneByID(ctx, id, userID)
	if errors.Is(err, database.ErrNotFound) {
		return expert, service.ErrEntityNotFound
	}
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return expert, nil
}
func (s *Service) UpdateExpert(ctx context.Context, expert entity.Expert, userID int) (entity.Expert, error) {
	const op = "ExpertService.UpdateExpert"

	expert, err := s.repo.Update(ctx, expert, userID)
	if errors.Is(err, database.ErrNotFound) {
		return expert, service.ErrEntityNotFound
	}
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return s.repo.Update(ctx, expert, userID)
}

type Repository interface {
	Insert(ctx context.Context, expert entity.Expert) (entity.Expert, error)
	Delete(ctx context.Context, id int, userID int) error
	Update(ctx context.Context, expert entity.Expert, userID int) (entity.Expert, error)
	GetOneByID(ctx context.Context, id int, userID int) (entity.Expert, error)
	GetListByUserID(ctx context.Context, userID int) ([]entity.Expert, error)
}
type Service struct {
	repo Repository
}

func NewExpertService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
