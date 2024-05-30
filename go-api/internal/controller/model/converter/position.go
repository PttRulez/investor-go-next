package converter

import (
	"github.com/pttrulez/investor-go/internal/controller/model/response"
	"github.com/pttrulez/investor-go/internal/entity"
)

func FromPositionToResponse(d *entity.Position) response.Position {
	return response.Position{
		Amount:       d.Amount,
		AveragePrice: d.AveragePrice,
		Comment:      d.Comment,
		CurrentPrice: d.CurrentPrice,
		CurrentCost:  d.CurrentCost,
		ShortName:    d.ShortName,
		TargetPrice:  d.TargetPrice,
		Ticker:       d.Ticker,
	}
}
