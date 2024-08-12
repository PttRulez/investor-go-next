package expert

import (
	"context"

	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/service"
)

func (s *Service) CreateNewExpert(ctx context.Context, expert *entity.Expert) error {
	return s.repo.Insert(ctx, expert)
}
func (s *Service) DeleteExpert(ctx context.Context, id int, userID int) error {
	if _, err := s.GetExpertByID(ctx, id, userID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
func (s *Service) GetListByUserID(ctx context.Context, userID int) ([]*entity.Expert, error) {
	return s.repo.GetListByUserID(ctx, userID)
}
func (s *Service) GetExpertByID(ctx context.Context, id int, userID int) (*entity.Expert, error) {
	expert, err := s.repo.GetOneByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if expert.UserID != userID {
		return nil, service.ErrEntityNotFound
	}
	return expert, nil
}
func (s *Service) UpdateExpert(ctx context.Context, expert *entity.Expert, userID int) error {
	if _, err := s.GetExpertByID(ctx, expert.ID, userID); err != nil {
		return err
	}
	return s.repo.Update(ctx, expert)
}

type Repository interface {
	Insert(ctx context.Context, expert *entity.Expert) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, expert *entity.Expert) error
	GetOneByID(ctx context.Context, id int) (*entity.Expert, error)
	GetListByUserID(ctx context.Context, userID int) ([]*entity.Expert, error)
}
type Service struct {
	repo Repository
}

func NewExpertService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
