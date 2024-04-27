package service

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
)

func (s *MoexBondDealService) CreateBondDeal(ctx context.Context, dealData *model.Deal, userId int) error {
	if dealData.SecurityId == 0 {
		_, _ = s.services.Moex.Bond.GetByISIN(ctx, dealData.Ticker)
	}
	err := s.repo.Deal.MoexBond.Insert(ctx, dealData)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealService.CreateDeal]: %w", err)
	}
	err = s.services.Moex.Bond.UpdatePositionInDB(ctx, dealData.PortfolioId,
		dealData.SecurityId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealService.CreateDeal]: %w", err)
	}

	return nil
}

func (s *MoexBondDealService) DeleteBondDeal(ctx context.Context, dealId int, userId int) error {
	err := s.repo.Deal.MoexBond.Delete(ctx, dealId)
	if err != nil {
		return fmt.Errorf("\n<-[MoexBondDealService.DeleteDeal]: %w", err)
	}
	return nil
}

type MoexBondDealService struct {
	repo     *repository.Repository
	services *Container
}

func NewMoexBondDealService(repo *repository.Repository, services *Container) *MoexBondDealService {
	return &MoexBondDealService{
		repo:     repo,
		services: services,
	}
}
