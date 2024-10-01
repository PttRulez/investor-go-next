package opinion

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/database"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateNewExpert(ctx context.Context, expert domain.Expert) (domain.Expert, error) {
	const op = "ExpertService.CreateNewExpert"

	e, err := s.expertRepo.Insert(ctx, expert)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}
func (s *Service) DeleteExpert(ctx context.Context, id int, userID int) error {
	const op = "ExpertService.DeleteExpert"

	err := s.expertRepo.Delete(ctx, id, userID)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}
func (s *Service) GetExpertsList(ctx context.Context, userID int) ([]domain.Expert, error) {
	const op = "ExpertService.GetListByUserID"

	e, err := s.expertRepo.GetListByUserID(ctx, userID)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}
func (s *Service) GetExpert(ctx context.Context, id int, userID int) (domain.Expert, error) {
	const op = "ExpertService.GetExpertByID"

	expert, err := s.expertRepo.GetOneByID(ctx, id, userID)
	if errors.Is(err, database.ErrNotFound) {
		return expert, service.ErrdomainNotFound
	}
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return expert, nil
}
func (s *Service) UpdateExpert(ctx context.Context, expert domain.Expert, userID int) (domain.Expert, error) {
	const op = "ExpertService.UpdateExpert"

	expert, err := s.expertRepo.Update(ctx, expert, userID)
	if errors.Is(err, database.ErrNotFound) {
		return expert, service.ErrdomainNotFound
	}
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return s.expertRepo.Update(ctx, expert, userID)
}
