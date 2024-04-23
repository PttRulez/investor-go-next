package services

import (
	"context"

	"github.com/pttrulez/investor-go/internal/types"
)

func (s *MoexBondDealService) CreateDeal(ctx context.Context, dealData *types.Deal, userId int) error {
	// Проверяем не чужой ли этой портфолио
	portfolio, err := s.repo.Portfolio.GetById(ctx, dealData.PortfolioId)
	if err != nil {
		return err
	}
	if portfolio.UserId != userId {
		return types.ErrNotYours
	}

	err = s.repo.Deal.MoexBond.Insert(ctx, dealData)
	if err != nil {
		return err
	}

	err = s.services.MoexBond.UpdatePositionInDB(ctx, dealData.PortfolioId, dealData.SecurityId)
	if err != nil {
		return err
	}

	return nil
}

type MoexBondDealService struct {
	repo     *types.Repository
	services *ServiceContainer
}

func NewMoexBondDealService(repo *types.Repository, services *ServiceContainer) *MoexBondDealService {
	return &MoexBondDealService{
		repo:     repo,
		services: services,
	}
}
