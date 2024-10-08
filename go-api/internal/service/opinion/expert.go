package opinion

import (
	"context"
	"errors"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/storage"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateNewExpert(ctx context.Context, expert domain.Expert) (domain.Expert, error) {
	const op = "OpinionService.CreateNewExpert"

	e, err := s.repo.InsertExpert(ctx, expert)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (s *Service) DeleteExpert(ctx context.Context, id int, userID int) error {
	const op = "OpinionService.DeleteExpert"

	err := s.repo.DeleteExpert(ctx, id, userID)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return service.ErrdomainNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

func (s *Service) GetExpertsList(ctx context.Context, userID int) ([]domain.Expert, error) {
	const op = "OpinionService.GetExpertsList"

	e, err := s.repo.GetExpertList(ctx, userID)
	if err != nil {
		return e, fmt.Errorf("%s: %w", op, err)
	}

	return e, nil
}

func (s *Service) GetExpert(ctx context.Context, id int, userID int) (domain.Expert, error) {
	const op = "OpinionService.GetExpert"

	expert, err := s.repo.GetExpert(ctx, id, userID)
	if errors.Is(err, storage.ErrNotFound) {
		return expert, service.ErrdomainNotFound
	}
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return expert, nil
}

func (s *Service) UpdateExpert(ctx context.Context, expert domain.Expert, userID int) (
	domain.Expert, error) {
	const op = "OpinionService.UpdateExpert"

	expert, err := s.repo.UpdateExpert(ctx, expert, userID)
	if errors.Is(err, storage.ErrNotFound) {
		return expert, service.ErrdomainNotFound
	}
	if err != nil {
		return expert, fmt.Errorf("%s: %w", op, err)
	}

	return s.repo.UpdateExpert(ctx, expert, userID)
}
