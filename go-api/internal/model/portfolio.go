package model

type Portfolio struct {
	BondDeals      []*Deal
	BondPositions  []*Position
	Cash           int
	Cashouts       []*Cashout
	CashoutsSum    int
	Compound       bool
	Deals          []*Deal
	Deposits       []*Deposit
	DepositsSum    int
	Id             int
	Name           string
	Positions      []*Position
	Profitability  int
	ShareDeals     []*Deal
	SharePositions []*Position
	TotalCost      int
	UserId         int
}
