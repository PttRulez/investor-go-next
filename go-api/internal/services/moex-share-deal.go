package services

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/types"
)

func (s *MoexShareDealService) CreateDeal(ctx context.Context, dealData *types.Deal, userId int) error {
	// Проверяем не чужой ли этой портфолио
	portfolio, err := s.repo.Portfolio.GetById(ctx, dealData.PortfolioId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}
	if portfolio.UserId != userId {
		return types.ErrNotYours
	}

	err = s.repo.Deal.MoexShares.Insert(ctx, dealData)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}
	err = s.services.Moex.Shares.UpdatePositionInDB(ctx, dealData.PortfolioId,
		dealData.SecurityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexShareDealService.CreateDeal]: %w", err)
	}

	return nil
}

type MoexShareDealService struct {
	repo     *types.Repository
	services *ServiceContainer
}

func NewMoexShareDealService(repo *types.Repository, services *ServiceContainer) *MoexShareDealService {
	return &MoexShareDealService{
		repo:     repo,
		services: services,
	}
}
