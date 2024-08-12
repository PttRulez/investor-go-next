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
		ID:       dto.ID,
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromPortfolioToShortPortfolio(portfolio *entity.Portfolio) *response.ShortPortfolio {
	return &response.ShortPortfolio{
		ID:       portfolio.ID,
		Compound: portfolio.Compound,
		Name:     portfolio.Name,
	}
}
