package entity

import (
	"time"

	"github.com/pttrulez/investor-go/internal/types"
)

type Deal struct {
	Amount       int
	Date         time.Time
	Exchange     types.Exchange
	Id           int
	PortfolioId  int
	Price        float64
	SecurityId   int
	SecurityType types.SecurityType
	Type         Type
	Ticker       string
	UserId       int
}

//type DeleteDealInfo struct {
//	Exchange     types.Exchange
//	SecurityType types.SecurityType
//	Id           int
//}

//type Position struct {
//	Amount       int
//	AveragePrice float64
//	Comment      string
//	CurrentPrice float64
//	CurrentCost  int
//	Exchange     types.Exchange
//	Id           int
//	PortfolioId  int
//	Secid        string
//	SecurityId   int
//	SecurityType types.SecurityType
//	ShortName    string
//	TargetPrice  float64
//}

type Type string

const (
	DtBuy  Type = "BUY"
	DtSell Type = "SELL"
)

func (e Type) Validate() bool {
	switch e {
	case DtBuy:
	case DtSell:
	default:
		return false
	}
	return true
}
