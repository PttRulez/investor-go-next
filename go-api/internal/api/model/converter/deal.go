package converter

import (
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/model"
)

func FromCreateDealDtoToDeal(dto *dto.CreateDeal) *model.Deal {
	return &model.Deal{
		Amount:       dto.Amount,
		Date:         dto.Date,
		Exchange:     dto.Exchange,
		PortfolioId:  dto.PortfolioId,
		Price:        dto.Price,
		Ticker:       dto.Ticker,
		SecurityType: dto.SecurityType,
		Type:         dto.Type,
	}
}

func FromDeleteDealDtoToDeal(dto *dto.DeleteDeal) *model.Deal {
	return &model.Deal{
		Id:           dto.Id,
		Exchange:     dto.Exchange,
		SecurityType: dto.SecurityType,
	}
}
