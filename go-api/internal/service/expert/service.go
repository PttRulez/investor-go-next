package expert

import (
	"context"
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/internal/errors"
)

func (s *Service) CreateNewExpert(ctx context.Context, expert *entity.Expert) error {
	return s.repo.Insert(ctx, expert)
}
func (s *Service) DeleteExpert(ctx context.Context, id int, userId int) error {
	if _, err := s.GetExpertById(ctx, id, userId); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
func (s *Service) GetListByUserId(ctx context.Context, userId int) ([]*entity.Expert, error) {
	return s.repo.GetListByUserId(ctx, userId)
}
func (s *Service) GetExpertById(ctx context.Context, id int, userId int) (*entity.Expert, error) {
	expert, err := s.repo.GetOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	if expert.UserId != userId {
		return nil, errors.ErrNotYours
	}
	return expert, nil
}
func (s *Service) UpdateExpert(ctx context.Context, expert *entity.Expert, userId int) error {
	if _, err := s.GetExpertById(ctx, expert.Id, userId); err != nil {
		return err
	}
	return s.repo.Update(ctx, expert)
}

type Repository interface {
	Insert(ctx context.Context, expert *entity.Expert) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, expert *entity.Expert) error
	GetOneById(ctx context.Context, id int) (*entity.Expert, error)
	GetListByUserId(ctx context.Context, userId int) ([]*entity.Expert, error)
}
type Service struct {
	repo Repository
}

func NewExpertService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
