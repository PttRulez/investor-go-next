package model

type Position struct {
	Amount       int
	AveragePrice float64
	Comment      string
	CurrentPrice float64
	CurrentCost  int
	Exchange     Exchange
	Id           int
	PortfolioId  int
	Secid        string
	SecurityId   int
	SecurityType SecurityType
	ShortName    string
	TargetPrice  float64
}
