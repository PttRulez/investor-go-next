package converter

import (
	"github.com/pttrulez/investor-go/internal/api/model/dto"
	"github.com/pttrulez/investor-go/internal/api/model/response"
	"github.com/pttrulez/investor-go/internal/model"
)

func FromCreatePortfolioDtoToPortfolio(dto *dto.CreatePortfolio) *model.Portfolio {
	return &model.Portfolio{
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromUpdatePortfolioDtoToPortfolio(dto *dto.UpdatePortfolio) *model.Portfolio {
	return &model.Portfolio{
		Id:       dto.Id,
		Compound: dto.Compound,
		Name:     dto.Name,
	}
}

func FromPortfolioToShortPortfolio(portfolio *model.Portfolio) *response.ShortPortfolio {
	return &response.ShortPortfolio{
		Id:       portfolio.Id,
		Compound: portfolio.Compound,
		Name:     portfolio.Name,
	}
}
