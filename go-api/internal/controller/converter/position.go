package converter

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromPositionToResponse(p entity.Position) api.PositionResponse {
	opinions := make([]api.OpinionResponse, 0, len(p.Opinions))
	for _, o := range p.Opinions {
		opinions = append(opinions, FromOpinionToOpinionResponse(o))
	}

	return api.PositionResponse{
		Amount:        p.Amount,
		AveragePrice:  p.AveragePrice,
		Comment:       p.Comment,
		CurrentPrice:  p.CurrentPrice,
		CurrentCost:   p.CurrentCost,
		Id:            p.ID,
		Opinions:      opinions,
		OpinionIds:    p.OpinionIDs,
		PortfolioName: p.PortfolioName,
		SecurityType:  api.SecurityType(p.SecurityType),
		ShortName:     p.ShortName,
		TargetPrice:   p.TargetPrice,
		Ticker:        p.Ticker,
	}
}
