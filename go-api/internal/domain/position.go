package domain

import "math"

func (p Position) UpdateByDeals(allDeals []Deal, decimalCount int) Position {
	// Calculate position amount
	var amount int
	var totalAmount int
	for _, d := range allDeals {
		amount = d.Amount
		if d.Type == DTSell {
			amount = -amount
		}
		totalAmount += amount
	}
	p.Amount = totalAmount

	// Calculate and add AveragePrice to position
	m := make(map[float64]int)
	for _, deal := range allDeals {
		if deal.Type == DTSell {
			continue
		}
		m[deal.Price] += deal.Amount
	}

	var avPrice float64
	for price, amount := range m {
		avPrice += float64(amount) / float64(p.Amount) * price
	}

	// Приведение средней цены к определенному кол-ву знаков после запятой
	d := float64(math.Pow10(decimalCount))
	p.AveragePrice = math.Floor(avPrice*d) / d
	return p
}

type Position struct {
	Amount        int
	AveragePrice  float64
	Board         ISSMoexBoard
	Comment       *string
	CurrentPrice  float64
	CurrentCost   int
	Exchange      Exchange
	ID            int
	OpinionIDs    []int
	Opinions      []Opinion `json:"opinions"`
	PortfolioID   int
	PortfolioName string
	SecurityType  SecurityType
	ShortName     string
	TargetPrice   *float64
	Ticker        string
	UserID        int
}

type PositionUpdateInfo struct {
	ID          int
	Comment     *string
	TargetPrice *float64
	UserID      int
}
