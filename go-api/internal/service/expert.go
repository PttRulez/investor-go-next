package service

import (
	"context"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (s *ExpertService) CreateNewExpert(ctx context.Context, expert *model.Expert) error {
	return s.repo.Expert.Insert(ctx, expert)
}
func (s *ExpertService) DeleteExpert(ctx context.Context, id int, userId int) error {
	if _, err := s.GetExpertById(ctx, id, userId); err != nil {
		return err
	}
	return s.repo.Expert.Delete(ctx, id)
}
func (s *ExpertService) GetExpertsListByUserId(ctx context.Context, userId int) ([]*model.Expert, error) {
	return s.repo.Expert.GetListByUserId(ctx, userId)
}
func (s *ExpertService) GetExpertById(ctx context.Context, id int, userId int) (*model.Expert, error) {
	expert, err := s.repo.Expert.GetOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	if expert.UserId != userId {
		return nil, httpresponse.ErrNotYours
	}
	return expert, nil
}
func (s *ExpertService) UpdateExpert(ctx context.Context, expert *model.Expert, userId int) error {
	if _, err := s.GetExpertById(ctx, expert.Id, userId); err != nil {
		return err
	}
	return s.repo.Expert.Update(ctx, expert)
}

type ExpertService struct {
	repo     *repository.Repository
	services *Container
}

func NewExpertService(repo *repository.Repository, services *Container) *ExpertService {
	return &ExpertService{
		repo:     repo,
		services: services,
	}
}
