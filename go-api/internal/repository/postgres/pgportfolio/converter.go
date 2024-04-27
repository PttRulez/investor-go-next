package pgportfolio

import "github.com/pttrulez/investor-go/internal/model"

func FromDBToPortfolio(p *Portfolio) *model.Portfolio {
	return &model.Portfolio{
		Id:       p.Id,
		Name:     p.Name,
		Compound: p.Compound,
		UserId:   p.UserId,
	}
}

func FromPortfolioToDB(p *model.Portfolio) *Portfolio {
	return &Portfolio{
		Id:       p.Id,
		Name:     p.Name,
		Compound: p.Compound,
		UserId:   p.UserId,
	}
}
