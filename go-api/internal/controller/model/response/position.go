package response

type Position struct {
	Amount       int     `json:"amount"`
	AveragePrice float64 `json:"averagePrice"`
	Comment      string  `json:"comment"`
	CurrentPrice float64 `json:"currentPrice"`
	CurrentCost  int     `json:"currentCost"`
	ShortName    string  `json:"shortName"`
	TargetPrice  float64 `json:"targetPrice"`
	Secid        string  `json:"secid"`
}
