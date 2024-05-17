package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/dto"
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromCreatePortfolioDtoToPortfolio(dto *dto.CreatePortfolio) *entity.Portfolio {
	return &entity.Portfolio{
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromUpdatePortfolioDtoToPortfolio(dto *dto.UpdatePortfolio) *entity.Portfolio {
	return &entity.Portfolio{
		Id:       dto.Id,
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromPortfolioToShortPortfolio(portfolio *entity.Portfolio) *response.ShortPortfolio {
	return &response.ShortPortfolio{
		Id:       portfolio.Id,
		Compound: portfolio.Compound,
		Name:     portfolio.Name,
	}
}
