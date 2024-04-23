package types

type Position struct {
	Id           int      `json:"id" db:"id"`
	Amount       int      `json:"amount" db:"amount"`
	AveragePrice float64  `json:"averagePrice" db:"average_price"`
	Comment      string   `json:"comment" db:"comment,omitempty"`
	CurrentPrice float64  `json:"currentPrice" db:"-"`
	CurrentCost  float64  `json:"currentCost" db:"-"`
	Exchange     Exchange `json:"exchange" db:"exchange"`
	PortfolioId  int      `json:"portfolioId" db:"portfolio_id"`
	SecurityId   int      `json:"securityId" db:"security_id"`
	TargetPrice  float64  `json:"targetPrice" db:"target_price,omitempty"`
	Ticker       string   `json:"ticker" db:"ticker"`
}

type BondPosition struct {
	Position
}

type SharePosition struct {
	Position
}
