package entity

import (
	"github.com/pttrulez/investor-go/internal/types"
)

type Position struct {
	Amount       int
	AveragePrice float64
	Board        ISSMoexBoard
	Comment      string
	CurrentPrice float64
	CurrentCost  int
	Exchange     types.Exchange
	Id           int
	PortfolioId  int
	SecurityType types.SecurityType
	ShortName    string
	TargetPrice  float64
	Ticker       string
	UserId       int
}
