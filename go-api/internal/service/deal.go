package service

import (
	"context"
	"errors"

	"github.com/pttrulez/investor-go/internal/model"
)

func (s *DealService) CreateDeal(ctx context.Context, deal *model.Deal, userId int) error {
	if deal.Exchange == model.EXCH_Moex {
		return s.services.Moex.Deal.CreateDeal(ctx, deal, userId)
	}
	return errors.New("Неправильно введена биржа")
}

func (s *DealService) DeleteDeal(ctx context.Context, deal *model.Deal, userId int) error {
	if deal.Exchange == model.EXCH_Moex {
		return s.services.Moex.Deal.DeleteDeal(ctx, deal, userId)
	}
	return errors.New("Неправильно введена биржа")
}

type DealService struct {
	services *Container
}
