package entity

type Position struct {
	Amount       int
	AveragePrice float64
	Board        ISSMoexBoard
	Comment      string
	CurrentPrice float64
	CurrentCost  int
	Exchange     Exchange
	Id           int
	PortfolioId  int
	SecurityType SecurityType
	ShortName    string
	TargetPrice  float64
	Ticker       string
	UserId       int
}
