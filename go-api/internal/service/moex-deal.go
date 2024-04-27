package service

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/model"
	"github.com/pttrulez/investor-go/internal/repository"
	httpresponse "github.com/pttrulez/investor-go/pkg/http-response"
)

func (s *MoexDealService) CreateDeal(ctx context.Context, deal *model.Deal, userId int) error {
	// Проверяем не чужой ли этой портфолио
	portfolio, err := s.repo.Portfolio.GetById(ctx, deal.PortfolioId)
	if err != nil {
		return fmt.Errorf("\n<-[DealService.CreateDeal]: %w", err)
	}
	if portfolio.UserId != userId {
		return httpresponse.ErrNotYours
	}

	if deal.SecurityType == model.ST_Bond {
		return s.services.Moex.BondDeal.CreateDeal(ctx, deal, userId)
	} else if deal.SecurityType == model.ST_Share {
		return s.services.Moex.ShareDeal.CreateDeal(ctx, deal, userId)
	}

	return fmt.Errorf("\n<-[DealService.CreateDeal]: Логика описана только для BOND и SHARE")
}

func (s *MoexDealService) DeleteDeal(ctx context.Context, deal *model.Deal, userId int) error {
	if deal.SecurityType == model.ST_Bond {
		return s.repo.Deal.MoexShare.Delete(ctx, deal.Id)
	} else if deal.SecurityType == model.ST_Share {
		return s.services.Moex.BondDeal.DeleteDeal(ctx, deal.Id, userId)
	}

	return nil
}

type MoexDealService struct {
	repo     *repository.Repository
	services *Container
}

func NewMoexDealService(repo *repository.Repository, services *Container) *MoexDealService {
	return &MoexDealService{
		repo:     repo,
		services: services,
	}
}
