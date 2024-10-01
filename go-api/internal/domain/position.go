package domain

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
