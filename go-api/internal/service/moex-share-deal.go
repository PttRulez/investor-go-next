package service

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
)

func (s *MoexShareDealService) CreateDeal(ctx context.Context, dealData *model.Deal, userId int) error {
	if dealData.SecurityId == 0 {
		_, _ = s.services.MoexShare.GetByTicker(ctx, dealData.Secid)
	}
	err := s.repo.Deal.MoexShare.Insert(ctx, dealData)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}
	err = s.services.MoexShare.UpdatePositionInDB(ctx, dealData.PortfolioId,
		dealData.SecurityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}

	return nil
}

func (s *MoexShareDealService) DeleteDeal(ctx context.Context, dealId int, userId int) error {
	err := s.repo.Deal.MoexShare.Delete(ctx, dealId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.DeleteDeal]: %w", err)
	}
	return nil
}



type MoexShareDealService struct {
	repo     *repository.Repository
	services *Container
}

func NewMoexShareDealService(repo *repository.Repository, services *Container) *MoexShareDealService {
	return &MoexShareDealService{
		repo:     repo,
		services: services,
	}
}
