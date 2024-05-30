package response

type Position struct {
	Amount       int     `json:"amount"`
	AveragePrice float64 `json:"averagePrice"`
	Comment      string  `json:"comment"`
	CurrentPrice float64 `json:"currentPrice"`
	CurrentCost  int     `json:"currentCost"`
	ShortName    string  `json:"shortName"`
	TargetPrice  float64 `json:"targetPrice"`
	Ticker       string  `json:"ticker"`
}

type BondPosition struct {
	Position
	Isin string `json:"isin"`
}

type SharePosition struct {
	Position
	Ticker string `json:"ticker"`
}
