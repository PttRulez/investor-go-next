package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromCreateDealDtoToDeal(dto *dto.CreateDeal) *entity.Deal {
	return &entity.Deal{
		Amount:       dto.Amount,
		Date:         dto.Date,
		Exchange:     dto.Exchange,
		PortfolioId:  dto.PortfolioId,
		Price:        dto.Price,
		SecurityType: dto.SecurityType,
		Ticker:       dto.Ticker,
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

func FromDealToResponse(d *entity.Deal) response.Deal {
	return response.Deal{
		Amount: d.Amount,
		Date:   d.Date,
		Id:     d.Id,
		Price:  d.Price,
		Ticker: d.Ticker,
	}
}
