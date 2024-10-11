package converter

import (
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
)

func FromPositionToResponse(p domain.Position) contracts.PositionResponse {
	opinions := make([]contracts.OpinionResponse, 0, len(p.Opinions))
	for _, o := range p.Opinions {
		opinions = append(opinions, FromOpinionToOpinionResponse(o))
	}

	return contracts.PositionResponse{
		Amount:        p.Amount,
		AveragePrice:  p.AveragePrice,
		Comment:       p.Comment,
		CurrentPrice:  p.CurrentPrice,
		CurrentCost:   p.CurrentCost,
		Id:            p.ID,
		Opinions:      opinions,
		OpinionIds:    p.OpinionIDs,
		PortfolioName: p.PortfolioName,
		SecurityType:  contracts.SecurityType(p.SecurityType),
		ShortName:     p.ShortName,
		TargetPrice:   p.TargetPrice,
		Ticker:        p.Ticker,
	}
}
