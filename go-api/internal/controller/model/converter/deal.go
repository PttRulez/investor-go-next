package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromCreateDealDtoToDeal(dto *dto.CreateDeal) *entity.Deal {
	return &entity.Deal{
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

func FromDeleteDealDtoToDeal(dto *dto.DeleteDeal) *entity.Deal {
	return &entity.Deal{
		Id:           dto.Id,
		Exchange:     dto.Exchange,
		SecurityType: dto.SecurityType,
	}
}
