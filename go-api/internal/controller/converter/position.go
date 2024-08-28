package converter

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromPositionToResponse(p entity.Position) api.PositionResponse {
	return api.PositionResponse{
		Amount:       p.Amount,
		AveragePrice: p.AveragePrice,
		Comment:      p.Comment,
		CurrentPrice: p.CurrentPrice,
		CurrentCost:  p.CurrentCost,
		SecurityType: api.SecurityType(p.SecurityType),
		ShortName:    p.ShortName,
		TargetPrice:  p.TargetPrice,
		Ticker:       p.Secid,
	}
}
