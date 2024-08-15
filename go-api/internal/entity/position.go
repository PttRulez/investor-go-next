package entity

type Position struct {
	Amount       int
	AveragePrice float64
	Board        ISSMoexBoard
	Comment      string
	CurrentPrice float64
	CurrentCost  int
	Exchange     Exchange
	ID           int
	PortfolioID  int
	SecurityType SecurityType
	ShortName    string
	TargetPrice  float64
	Secid        string
	UserID       int
}
