package entity

import (
	"github.com/pttrulez/investor-go/internal/types"
)

type Position struct {
	Amount       int
	AveragePrice float64
	Comment      string
	CurrentPrice float64
	CurrentCost  int
	Exchange     types.Exchange
	Id           int
	PortfolioId  int
	Secid        string
	SecurityId   int
	SecurityType types.SecurityType
	ShortName    string
	TargetPrice  float64
}
